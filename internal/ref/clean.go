package ref

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/nickng/bibtex"
	"github.com/stormvirux/bibrefer/pkg/request"
	"go4.org/bytereplacer"
	"log"
	"os"
	"strings"
	"unicode"
)

type Clean struct {
	BibKey      bool
	FullJournal bool
	FullAuthor  bool
	NoNewline   bool
	Verbose     bool
	Output      string
	TexFile     []string
	BibFile     string
}

func (c *Clean) Run(query []string) (string, error) {
	var (
		queryTxt   string
		references *bibtex.BibTex
		err        error
	)

	switch {
	case len(query) > 0:
		queryTxt = strings.Join(query, " ")
		verbosePrint(c.Verbose, fmt.Sprintf("Provided reference:\n %s",
			queryTxt), os.Stdout)
		references, err = bibtex.Parse(strings.NewReader(queryTxt))
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
	default:
		verbosePrint(c.Verbose, fmt.Sprintf("Reading bibtex file: %s", c.BibFile), os.Stdout)
		references, err = readBibtexFile(c.BibFile)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
	}

	var (
		finalRef        = make([]string, 0, len(references.Entries))
		doiPresent      = make([]string, 0, len(references.Entries))
		doiMissingIndex = make([]int, 0, len(references.Entries))
		doiIssuesIndex  = make([]int, 0, len(references.Entries))
		doiNames        = make([]string, 0, len(references.Entries))
		doiArxiv        = make([]bool, 0, len(references.Entries))
		j               = 0
	)

	var key = make(map[string]string)

	verbosePrint(c.Verbose, "Processing the bib file", os.Stdout)
	for i := 0; i < len(references.Entries); i++ {
		lowerDictKey(references.Entries[i].Fields)
		if doiS, ok := references.Entries[i].Fields["doi"]; ok {
			key[references.Entries[i].CiteName] = references.Entries[i].CiteName
			doiPresent = append(doiPresent, doiS.String())
			continue
		}

		doiMissingIndex = append(doiMissingIndex, i)
		doiNames = append(doiNames, references.Entries[i].Fields["title"].String())
		doiArxiv = append(doiArxiv, false)
		if v, ok := references.Entries[i].Fields["journal"]; ok {
			doiArxiv[len(doiArxiv)-1] = strings.Contains(strings.ToLower(v.String()), "arxiv")
		}
		j++
	}

	j = 0
	for i, doiTxt := range doiPresent {
		isValidDoi, strippedDoi := verifyDOI(doiTxt)
		if !isValidDoi {
			verbosePrint(c.Verbose, fmt.Sprintf("Invalid DOI %s. Searching with name: %s", doiTxt,
				references.Entries[i].Fields["title"].String()), os.Stdout)
			doiIssuesIndex = append(doiIssuesIndex, i)
			j++
		}

		r, err := request.RefDoi(strippedDoi, "bibtex")
		b, _ := bibtex.Parse(strings.NewReader(r))
		key[references.Entries[i].CiteName] = b.Entries[0].CiteName
		finalRef = append(finalRef, r)
		if err != nil {
			verbosePrint(c.Verbose, fmt.Sprintf("Cannot find with DOI %s. Searching with name: %s", doiTxt,
				references.Entries[i].Fields["title"].String()), os.Stdout)
			doiIssuesIndex = append(doiIssuesIndex, i)
			j++
		}
	}

	j = 0
	var d string
	for i, doiTxt := range doiMissingIndex {
		switch doiArxiv[i] {
		case true:
			d, err = request.DoiDataCite(doiNames[i])
		case false:
			d, err = request.DoiCrossRef(doiNames[i])
		}
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		r, err := request.RefDoi(d, "bibtex")
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		isEq := checkEq(r, references.Entries[doiTxt].Fields["title"].String())
		if isEq {
			b, _ := bibtex.Parse(strings.NewReader(r))
			key[references.Entries[i].CiteName] = b.Entries[0].CiteName
			finalRef = append(finalRef, r)
		}
		if r == "" || !isEq {
			finalRef = append(finalRef, references.Entries[doiTxt].String())
			j++
		}
	}
	verbosePrint(c.Verbose, fmt.Sprintf("%d entries not found. Using old ones\n", j), os.Stdout)
	bibBuilder := &strings.Builder{}

	if c.BibKey && c.FullAuthor && c.FullJournal {
		for i := 0; i < len(finalRef); i++ {
			bibBuilder.WriteString(finalRef[i] + "\n")
		}
		if c.BibFile == "" {
			return bibBuilder.String(), nil
		}
		err := writeFile(bibBuilder.String(), c.BibFile)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		err = modifyTex(c.TexFile, key)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		return "", nil
	}
	// TODO: Add key for keeping original bibkey
	for i := 0; i < len(finalRef); i++ {
		s := bibCleanWithFlags(c.BibKey, c.FullJournal, c.FullAuthor,
			finalRef[i], c.Verbose)
		b, _ := bibtex.Parse(strings.NewReader(s))
		key[references.Entries[i].CiteName] = b.Entries[0].CiteName
		bibBuilder.WriteString(s + "\n")

	}
	if c.BibFile == "" {
		return bibBuilder.String(), nil
	}
	err = writeFile(bibBuilder.String(), c.BibFile)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	err = modifyTex(c.TexFile, key)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return "", nil

}

func writeFile(out string, file string) error {
	f, err := os.Create(fmt.Sprintf("%s-bibrefer.bib", strings.Replace(file, ".bib", "", 1)))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("%v", err)
		}
	}(f)
	w := bufio.NewWriter(f)
	_, _ = w.WriteString(out)
	err = w.Flush()
	if err != nil {
		return err
	}
	return err
}

func checkEq(n string, o string) bool {
	p, _ := bibtex.Parse(strings.NewReader(n))
	return strings.EqualFold(p.Entries[0].Fields["title"].String(), o)
	/*r := regexp.MustCompile(`{(?P<name>[\p{L}\d].*)}`)
	for i := 0; i < len(bib); i++ {
		if strings.Contains(bib[i], "title") {
			t := r.FindStringSubmatch(bib[i])
			namedIndex := r.SubexpIndex("name")
			if namedIndex == -1 {
				return false
			}
			if len(t) > namedIndex {

			}
		}
	}
	return false*/
}

func lowerDictKey(m map[string]bibtex.BibString) {
	for k, v := range m {
		if IsUpper(k) {
			m[strings.ToLower(k)] = v
			delete(m, k)
		}
	}
}

func modifyTex(texfiles []string, key map[string]string) error {
	var rKeys []string
	for k, v := range key {
		rKeys = append(rKeys, k, v)
	}
	//
	for _, texfile := range texfiles {
		b, err := os.ReadFile(texfile)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		o := bytereplacer.New(rKeys...).Replace(b)
		err = writeTex(o, texfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeTex(out []byte, file string) error {
	f, err := os.Create(fmt.Sprintf("%s-bibrefer.tex", strings.Replace(file, ".tex", "", 1)))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("%v", err)
		}
	}(f)
	w := bufio.NewWriter(f)
	_, _ = w.Write(out)
	err = w.Flush()
	if err != nil {
		return err
	}
	return err
}

func readBibtexFile(bibFile string) (*bibtex.BibTex, error) {
	b, err := os.ReadFile(bibFile)
	if err != nil {
		return bibtex.NewBibTex(), fmt.Errorf("cannot read %s: %w", bibFile, err)
	}
	s, err := bibtex.Parse(bytes.NewReader(b))
	if err != nil {
		return bibtex.NewBibTex(), fmt.Errorf("cannot parse valid bibtex file %s: %w", bibFile, err)
	}
	return s, nil
}

func IsUpper(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
