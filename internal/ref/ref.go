package ref

import (
	"fmt"
	"github.com/hscells/doi"
	"github.com/stormvirux/bibrefer/pkg/request"
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
	/*r.BibKey := flags[0]
	r.FullJournal := flags[1]
	r.FullAuthor := flags[2]
	r.NoNewline := flags[3]
	r.Verbose := flags[4]*/

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
		return reference + "\n", nil
	}

	verbosePrint(r.Verbose, "Cleaning the obtained reference", os.Stdout)
	finalBib := bibCleanWithFlags(r.BibKey, r.FullJournal, r.FullAuthor, reference, r.Verbose)
	if r.NoNewline {
		return "\n" + finalBib, nil
	}
	return finalBib + "\n", nil
}

func verifyDOI(doiTxt string) (isValid bool, validDoi string) {
	r := regexp.MustCompile(`(?i)(https?://)?(?:[^/.\s]+\.)*([a-zA-Z]*\.[a-zA-Z]*/|dx\.doi\.org/|doi\.acm\.org/)`)
	strippedDOI := r.ReplaceAllString(doiTxt, "")
	d, err := doi.Parse(strippedDOI)
	if err != nil {
		println(err)
		println(d.ToString())
	}
	if !d.IsValid() {
		fmt.Println("Invalid DOI")
		return false, ""
	}
	return true, strippedDOI
}

func bibCleanWithFlags(bibKey bool, fullJournal bool, fullAuthor bool, bibEntry string, verbose bool) string {
	bibLineByLine := strings.Split(bibEntry, "\n")
	bibIndices := getIndexBib(bibLineByLine)
	var lName []string

	var indices []uint

	if bibIndices["authorIndex"] != 255 {
		verbosePrint(verbose, "Abbreviating the author names", os.Stdout)
		r1 := regexp.MustCompile(`[{](?P<name>[\p{L},?\s.'{}\\-]+)[}]`)
		authorField := r1.FindStringSubmatch(bibLineByLine[bibIndices["authorIndex"]])
		r1Index := r1.SubexpIndex("name")
		authors := strings.Split(authorField[r1Index], "and")
		if !fullAuthor {
			var newAuthor strings.Builder
			newAuthor.Grow(100)
			for i := 0; i < len(authors); i++ {
				author := strings.Split(strings.TrimSpace(authors[i]), " ")
				newAuthor.Reset()
				for j := 0; j < len(author)-1; j++ {
					newAuthor.WriteString(author[j][0:1])
					newAuthor.WriteString(". ")
				}
				authors[i] = newAuthor.String() + author[len(author)-1] // r2.ReplaceAllString(authors[i], "$1. $2")
			}
			// firstAuthLName = strings.Split(strings.TrimSpace(authors[0]), " ")
			bibLineByLine[bibIndices["authorIndex"]] = strings.ReplaceAll(bibLineByLine[bibIndices["authorIndex"]], authorField[r1Index], strings.Join(authors, " and "))
		}
		lName = strings.Split(strings.TrimSpace(authors[0]), " ")
	}

	isArXiv := strings.Contains(bibLineByLine[0], "arxiv")
	if !bibKey && !isArXiv {
		verbosePrint(verbose, "Updating the bib key", os.Stdout)
		bibLineByLine[0] = strings.Replace(bibLineByLine[0], "_", ":", 1)
		bibLineByLine[0] = strings.Replace(bibLineByLine[0], ",", "", 1)
		r := regexp.MustCompile(`\b10\.(\d+\.*)+/(?P<name>[a-zA-Z]+)\d*\.?(([^\s.])+\.*)+\b`)
		if strings.Contains(bibLineByLine[bibIndices["doiIndex"]], "10.1016") {
			r = regexp.MustCompile(`\b10.1016/j.(?P<name>[a-zA-Z]+)\d*\.(([^\s.])+\.*)+\b`)
		}
		venueName := r.FindStringSubmatch(bibLineByLine[bibIndices["doiIndex"]])
		namedIndex := r.SubexpIndex("name")
		if len(venueName) > namedIndex {
			bibLineByLine[0] += `:` + strings.ToUpper(venueName[namedIndex])
		}
		bibLineByLine[0] += `,`
	}

	if !bibKey && isArXiv {
		r := regexp.MustCompile(`\s*year\s?=\s?\{(\d{4})},`)
		yearString := r.ReplaceAllString(bibLineByLine[bibIndices["yearIndex"]], "$1")
		replaceString := `@misc{` + lName[len(lName)-1] + `:` + yearString + `:ArXiv,`
		bibLineByLine[0] = strings.ReplaceAll(bibLineByLine[0], bibLineByLine[0], replaceString)
	}

	if !fullJournal && bibIndices["journalIndex"] != 255 {
		verbosePrint(verbose, "Abbreviating the Journal names", os.Stdout)
		bibLineByLine[bibIndices["journalIndex"]] = replaceJournalStrings(bibLineByLine[bibIndices["journalIndex"]])
	}

	verbosePrint(verbose, "Removing url and month", os.Stdout)
	if bibIndices["urlIndex"] == 255 && bibIndices["monthIndex"] == 255 {
		return strings.Join(bibLineByLine, "\n")
	}

	if bibIndices["urlIndex"] != 255 && bibIndices["monthIndex"] != 255 {
		if bibIndices["urlIndex"] < bibIndices["monthIndex"] {
			indices = append(indices, bibIndices["urlIndex"], bibIndices["monthIndex"])
			return strings.TrimSpace(strings.Join(removeTwoIndexLinear(bibLineByLine, indices), "\n"))
		}
		indices = append(indices, bibIndices["monthIndex"], bibIndices["urlIndex"])
		return strings.TrimSpace(strings.Join(removeTwoIndexLinear(bibLineByLine, indices), "\n"))
	}

	if bibIndices["urlIndex"] == 255 {
		indices = append(indices, bibIndices["monthIndex"])
		return strings.TrimSpace(strings.Join(removeTwoIndexLinear(bibLineByLine, indices), "\n"))
	}
	indices = append(indices, bibIndices["urlIndex"])
	return strings.TrimSpace(strings.Join(removeTwoIndexLinear(bibLineByLine, indices), "\n"))
}

func removeTwoIndexLinear(s []string, indices []uint) []string {
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

func replaceJournalStrings(journalEntry string) string {
	var abbr = []string{"Journal", "J.", "Electrical", "Elect.", "Computer", "Comput.", "Engineering", "Eng.", "Communications", "Commun.", "Magazine", "Mag.", "Aerospace", "Aerosp.", "Electronics", "Electron.", "Systems", "Syst.", "Annals", "Ann.", "History", "Hist.", "Computing", "Comput.", "Propagation", "Propag.", "Letters", "Lett.", "Society", "Soc.", "Tutorials", "Tuts.", "Intelligence", "Intell.", "Computational", "Comput.", "Intelligence", "Intell.", "Science", "Sci.", "Applications", "Appl.", "Architecture", "Archit.", "Graphics", "Graph.", "Design", "Des.", "Distributed", "Distrib.", "Management", "Manag.", "Review", "Rev.", "Medicine", "Med.", "Sensing", "Sens.", "Professional", "Prof.", "Industry", "Ind.", "Industrial", "Ind.", "Instrumentation", "Instrum.", "Measurement", "Meas.", "Intelligent", "Intell.", "Transportation", "Transp.", "Networking", "Netw.", "Robotics", "Robot.", "Automation", "Autom.", "Selected", "Sel.", "Automation", "Autom.", "Applied", "Appl.", "Processing", "Process.", "Techniques", "Techn.", "Technology", "Technol.", "Sciences", "Sci.", "Software", "Softw.", "Transactions", "Trans.", "Advanced", "Adv.", "Information", "Inf.", "Knowledge", "Knowl.", "Learning", "Learn.", "Analysis", "Anal.", "Machine", "Mach.", "Reliability", "Rel.", "Optimization", "Optim.", "Research", "Res.", "Mechanics", "Mech.", "Proceedings", "Proc.", "Royal", "R.", "Society", "Soc", "Annals", "Ann.", "Resources", "Resour.", "Surface", "Surf.", "Processes", "Proc.", "National", "Nat.", "Computers", "Comput.", "Geotechnics", "Geotech.", "Academy", "Acad.", "Sciences", "Sci.", "Quaternary", "Quat.", "Physical", "Phys.", "Planetary", "Planet.", "Quarterly", "Q.", "Geological", "Geol.", "Statistical", "Stat.", "Applied", "Appl.", "Physics", "Phys.", "Geoscience", "Geosci.", "Landforms", "Land.", "Science", "Sci.", "Annual", "Ann.", "International", "Int.", "Numerical", "Numer.", "Methods", "Meth.", "Geomechanics", "Geomech.", "Analytical", "Anal.", "Advances", "Adv.", "Modeling", "Mod.", "for", "", "of", "", "and", "", "in", "", "on", "", "&", "", "the", ""}
	return strings.NewReplacer(abbr...).Replace(journalEntry)
}
