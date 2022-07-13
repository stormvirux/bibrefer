package ref

import (
	"fmt"
	"github.com/hscells/doi"
	"github.com/nickng/bibtex"
	"github.com/stormvirux/bibrefer/pkg/request"
	"log"
	"os"
	"regexp"
	"strings"
)

// TODO: Handling Arxiv parse using Atom feed and check for new Doi

type Ref struct {
	BibKey      bool
	FullJournal bool
	FullAuthor  bool
	NoNewline   bool
	Verbose     bool
	Output      string
}

func (r *Ref) Run(query []string) (string, error) {
	// doiUser := query
	doiTxt := query[0]

	isValidDoi, strippedDoi := verifyDOI(doiTxt)
	if !isValidDoi {
		fmt.Printf("Invalid DOI: %s\n", doiTxt)
		return "", nil
	}
	verbosePrint(r.Verbose, fmt.Sprintf("Provided valid DOI: %s", strippedDoi), os.Stdout)
	isValidOutput := r.Output == "json" || r.Output == "bibtex" || r.Output == "xml"
	if !isValidOutput {
		return "", fmt.Errorf("invalid output format: %s", r.Output)
	}
	reference, err := request.RefDoi(doiTxt, r.Output)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if strings.Contains(reference, "DOI cannot be found") {
		return "", fmt.Errorf("DOI cannot be found")
	}
	verbosePrint(r.Verbose, fmt.Sprintf("Refernce for DOI: %s obtained", strippedDoi), os.Stdout)

	if r.Output != "bibtex" {
		return reference, nil
	}
	if r.BibKey && r.FullAuthor && r.FullJournal {
		if r.NoNewline {
			return "\n" + reference, nil
		}
		return "\n" + reference + "\n", nil
	}

	verbosePrint(r.Verbose, "Cleaning the obtained reference", os.Stdout)
	finalBib := bibCleanWithFlags(r.BibKey, r.FullJournal, r.FullAuthor, reference, r.Verbose)
	if r.NoNewline {
		return "\n" + finalBib, nil
	}
	return "\n" + finalBib + "\n", nil
}

func verifyDOI(doiTxt string) (isValid bool, validDoi string) {
	r := regexp.MustCompile(`(?i)(https?://)?(?:[^/.\s]+\.)*([a-zA-Z]*\.[a-zA-Z]*/|dx\.doi\.org/|doi\.acm\.org/)`)
	strippedDOI := r.ReplaceAllString(doiTxt, "")
	// fmt.Printf("The %s: %v\n", doiTxt, strippedDOI)
	d, err := doi.Parse(strippedDOI)
	if err != nil {
		verbosePrint(true, fmt.Sprintf("%v", err), os.Stderr)
		toString, _ := d.ToString()
		verbosePrint(true, toString, os.Stderr)
	}
	if !d.IsValid() {
		fmt.Println("Invalid DOI")
		return false, ""
	}
	return true, strippedDOI
}

func reverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func correctFirstName(a []string) []string {
	fName := a[0]
	lName := a[len(a)-1]
	fnames := strings.Split(strings.TrimSpace(fName), " ")
	var lNames strings.Builder
	lNames.Grow(100)
	for _, f := range fnames[1:] {
		lNames.WriteString(strings.TrimSpace(f))
		lNames.WriteString(" ")
	}
	return []string{strings.TrimSpace(fnames[0]), strings.TrimSpace(lNames.String() + strings.TrimSpace(lName))}
}

// Todo: Use Bib parser to modify the reference.
func bibCleanWithFlags(bibKey bool, fullJournal bool, fullAuthor bool, bibEntry string, verbose bool) string {
	entry, err := bibtex.Parse(strings.NewReader(bibEntry))
	if err != nil {
		log.Println(err)
	}
	curEntry := entry.Entries[0]
	var lName []string

	if author, ok := curEntry.Fields["author"]; ok {
		authors := strings.Split(author.String(), "and")

		if !fullAuthor {
			var newAuthor strings.Builder
			newAuthor.Grow(100)
			for i := 0; i < len(authors); i++ {
				author := strings.Split(strings.TrimSpace(authors[i]), " ")
				if strings.Contains(authors[i], ",") {
					author = strings.Split(strings.TrimSpace(authors[i]), ",")
					author = reverseSlice(author)
					author = correctFirstName(author)
				}
				newAuthor.Reset()
				for j := 0; j < len(author)-1; j++ {
					newAuthor.WriteString(author[j][0:1])
					newAuthor.WriteString(". ")
				}
				authors[i] = newAuthor.String() + author[len(author)-1]
			}
			ns := strings.Join(authors, " and ")
			curEntry.AddField("author", bibtex.BibConst(ns))
		}
		lName = strings.Split(strings.TrimSpace(authors[0]), " ")
	}

	isArXiv := strings.Contains(strings.ToLower(curEntry.CiteName), "arxiv")
	if !bibKey && !isArXiv {
		verbosePrint(verbose, "Updating the bib key", os.Stdout)
		if t, ok := curEntry.Fields["title"]; ok {
			curEntry.CiteName = strings.Split(strings.TrimSpace(t.String()), " ")[0]
			curEntry.CiteName = strings.ReplaceAll(curEntry.CiteName, "{", "")
			curEntry.CiteName = strings.ReplaceAll(curEntry.CiteName, "}", "")
		}
		if len(lName) > 0 {
			curEntry.CiteName = strings.TrimSpace(lName[len(lName)-1])
		}
		if y, ok := curEntry.Fields["year"]; ok {
			curEntry.CiteName = curEntry.CiteName + `:` + y.String()
		}

		r := regexp.MustCompile(`\b10\.(\d+\.*)+/(?P<name>[a-zA-Z]+)\d*\.?(([^\s.])+\.*)+\b`)
		if _, ok := curEntry.Fields["doi"]; ok && strings.Contains(curEntry.Fields["doi"].String(), "10.1016") {
			r = regexp.MustCompile(`\b10.1016/j.(?P<name>[a-zA-Z]+)\d*\.(([^\s.])+\.*)+\b`)
		}
		var venueName []string
		if _, ok := curEntry.Fields["doi"]; ok {
			venueName = r.FindStringSubmatch(curEntry.Fields["doi"].String())
			namedIndex := r.SubexpIndex("name")
			if len(venueName) > namedIndex {
				curEntry.CiteName += `:` + strings.ToUpper(venueName[namedIndex])
			}
		}
	}

	if !bibKey && isArXiv {
		curEntry.CiteName = lName[len(lName)-1] + `:` + curEntry.Fields["year"].String() + `:ArXiv`
	}

	if j, ok := curEntry.Fields["journal"]; !fullJournal && ok {
		verbosePrint(verbose, "Abbreviating the Journal names", os.Stdout)
		var replace = []string{"{", "", "}", "", "(", "", ")", ""}
		t := strings.NewReplacer(replace...).Replace(strings.TrimSpace(j.String()))
		if val, ok := d[strings.ToLower(strings.TrimSpace(t))]; ok {
			curEntry.AddField("journal", bibtex.BibConst(val))
		} else {
			curEntry.AddField("journal", bibtex.BibConst(replaceJournalStrings(j.String())))
		}
	}

	verbosePrint(verbose, "Removing url and month", os.Stdout)

	delete(curEntry.Fields, "url")
	delete(curEntry.Fields, "month")

	return prettyPrint(curEntry)
}

/*func removeTwoIndexLinear(s []string, indices []uint) []string {
	ret := make([]string, len(s)-len(indices)+1)
	w := 0
loop:
	for i, x := range s {
		for _, index := range indices {
			if index == uint(i) {
				continue loop
			}
		}
		ret[w] = x
		w++
	}
	return ret[0:w]
}
*/

/*func replaceJournalStringsDB(key string) string {
	db, err := pogreb.Open("abbrv.db", nil)
	if err != nil {
		log.Fatal("Unable to open database", err)
	}
	defer func(db *pogreb.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Unable to close database", err)
		}
	}(db)

	/*it := db.Items()
	for {
		key, val, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s %s", key, val)
	}
	value, err := db.Get([]byte(strings.ToLower(key)))
	if err != nil {
		log.Panic("Unable to get value for key "+key, err)
	}
	return string(value)
}*/

func replaceJournalStrings(journalEntry string) string {
	var abbr = []string{"Journal", "J.", "Electrical",
		"Elect.", "Computer", "Comput.", "Engineering", "Eng.",
		"Communications", "Commun.", "Magazine", "Mag.", "Aerospace", "Aerosp.",
		"Electronics", "Electron.", "Systems", "Syst.", "Annals", "Ann.",
		"History", "Hist.", "Propagation", "Propag.", "Letters", "Lett.",
		"Society", "Soc.", "Tutorials", "Tuts.",
		"Computational", "Comput.", "Science", "Sci.",
		"Applications", "Appl.", "Architecture", "Archit.", "Graphics", "Graph.",
		"Design", "Des.", "Distributed", "Distrib.", "Management", "Manag.",
		"Review", "Rev.", "Medicine", "Med.", "Sensing", "Sens.",
		"Professional", "Prof.", "Industry", "Ind.", "Industrial", "Ind.",
		"Instrumentation", "Instrum.", "Intelligent", "Intell.",
		"Transportation", "Transp.", "Networking", "Netw.", "Robotics", "Robot.",
		"Automation", "Autom.", "Selected", "Sel.", "Automation", "Autom.",
		"Applied", "Appl.", "Processing", "Process.", "Techniques", "Techn.",
		"Technology", "Technol.", "Sciences", "Sci.",
		"Software", "Softw.", "Transactions", "Trans.", "Advanced", "Adv.",
		"Information", "Inf.", "Knowledge", "Knowl.", "Learning", "Learn.",
		"Analysis", "Anal.", "Machine", "Mach.", "Reliability", "Rel.",
		"Optimization", "Optim.", "Research", "Res.", "Mechanics", "Mech.",
		"Proceedings", "Proc.", "Royal", "R.", "Society", "Soc", "Annals", "Ann.",
		"Resources", "Resour.", "Surface", "Surf.", "Processes", "Proc.",
		"National", "Nat.", "Computers", "Comput.", "Geotechnics", "Geotech.",
		"Academy", "Acad.", "Sciences", "Sci.", "Quaternary", "Quat.",
		"Physical", "Phys.", "Planetary", "Planet.", "Quarterly", "Q.",
		"Geological", "Geol.", "Statistical", "Stat.", "Applied", "Appl.",
		"Physics", "Phys.", "Geoscience", "Geosci.", "Landforms", "Land.",
		"Science", "Sci.", "Annual", "Ann.", "International", "Int.",
		"Numerical", "Numer.", "Methods", "Meth.", "Geomechanics", "Geomech.",
		"Analytical", "Anal.", "Advances", "Adv.", "Modeling", "Mod.",
		"Measurement", "Meas",
		"for", "", "of", "", "and", "", "in", "", "&", ""}
	return strings.NewReplacer(abbr...).Replace(journalEntry)
}
