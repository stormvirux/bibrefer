package ref

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/nickng/bibtex"
	"github.com/stormvirux/bibrefer/pkg/bytereplacer"
	"github.com/stormvirux/bibrefer/pkg/request"
	"log"
	"math"
	"os"
	"strings"
	"sync"
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

type refMeta struct {
	doi        string
	keyOld     string
	keyNew     string
	title      string
	ref        string
	refIndex   int
	isValidDoi bool
	isArxiv    bool
	author     string
}

func (c *Clean) Run(query []string) (string, error) {
	var (
		references *bibtex.BibTex
		err        error
	)

	references, err = readBib(query, c)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	var (
		doiPresents = make([]refMeta, 0, len(references.Entries))
		doiMisses   = make([]refMeta, 0, len(references.Entries))
	)

	verbosePrint(c.Verbose, "Processing the bib file", os.Stdout)
	for i := 0; i < len(references.Entries); i++ {
		lowerDictKey(references.Entries[i].Fields)
		if doiS, ok := references.Entries[i].Fields["doi"]; ok {
			dp := refMeta{doi: doiS.String(),
				keyOld: references.Entries[i].CiteName, refIndex: i}
			if j, ok := references.Entries[i].Fields["journal"]; ok {
				dp.title = j.String()
			}
			if j, ok := references.Entries[i].Fields["author"]; ok {
				dp.author = j.String()
			}
			doiPresents = append(doiPresents, dp)
			continue
		}
		dm := refMeta{refIndex: i, title: references.Entries[i].Fields["title"].String(),
			isArxiv: false, keyOld: references.Entries[i].CiteName}
		if v, ok := references.Entries[i].Fields["journal"]; ok {
			dm.isArxiv = strings.Contains(strings.ToLower(v.String()), "arxiv")
		}
		if v, ok := references.Entries[i].Fields["author"]; ok {
			dm.author = v.String()
		}
		doiMisses = append(doiMisses, dm)
	}
	finalRef, oldRef := asyncReq(doiPresents, doiMisses, c)

	var key = make(map[string]string, len(references.Entries))
	keyTable(finalRef, key)
	bibBuilder := &strings.Builder{}
	if c.BibKey && c.FullAuthor && c.FullJournal {
		for i := 0; i < len(finalRef); i++ {
			bibBuilder.WriteString(finalRef[i].ref + "\n")
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
	for i := 0; i < len(finalRef); i++ {
		entry, _ := bibtex.Parse(strings.NewReader(finalRef[i].ref))
		if _, ok := entry.Entries[0].Fields["title"]; !ok && finalRef[i].title != "" {
			entry.Entries[0].Fields["title"] = bibtex.BibConst(finalRef[i].title)
		}
		if _, ok := entry.Entries[0].Fields["author"]; !ok && finalRef[i].author != "" {
			entry.Entries[0].Fields["author"] = bibtex.BibConst(finalRef[i].author)
		}
		s := bibCleanWithFlags(c.BibKey, c.FullJournal, c.FullAuthor,
			entry.Entries[0].String(), c.Verbose)
		b, _ := bibtex.Parse(strings.NewReader(s))
		key[references.Entries[i].CiteName] = b.Entries[0].CiteName
		// finalRef[i].keyNew = b.Entries[0].CiteName
		bibBuilder.WriteString(s + "\n")
	}
	for v := range oldRef {
		t := references.Entries[oldRef[v].refIndex]
		bibBuilder.WriteString(t.String() + "\n")
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

func readBib(query []string, c *Clean) (*bibtex.BibTex, error) {
	var references *bibtex.BibTex
	var err error
	switch {
	case len(query) > 0:
		queryTxt := strings.Join(query, " ")
		verbosePrint(c.Verbose, fmt.Sprintf("Provided reference:\n %s",
			queryTxt), os.Stdout)
		references, err = bibtex.Parse(strings.NewReader(queryTxt))
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	default:
		verbosePrint(c.Verbose, fmt.Sprintf("Reading bibtex file: %s", c.BibFile), os.Stdout)
		references, err = readBibtexFile(c.BibFile)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}
	return references, nil
}

func keyTable(doiPresents []refMeta, key map[string]string) {
	for i := range doiPresents {
		key[doiPresents[i].keyOld] = doiPresents[i].keyNew
	}
}

func streamInput(done <-chan struct{}, inputs []refMeta) <-chan refMeta {
	inputCh := make(chan refMeta)
	go func() {
		defer close(inputCh)
		for _, input := range inputs {
			select {
			case inputCh <- input:
			case <-done:
				break
			}
		}
	}()
	return inputCh
}

func loadBalance(len1, len2 int) (int, int, int) {
	const maxConcurrent = 20
	l := float64(len1) + float64(len2)
	var mx1 = int(math.Ceil((maxConcurrent / l) * float64(len1)))
	var mx2 = int(math.Floor((maxConcurrent / l) * float64(len2)))
	isLessThan := len1+len2 < 20
	if isLessThan {
		mx1 = len1
		mx2 = len2
	}
	var toAdd = mx1 + mx2
	return mx1, mx2, toAdd
}

func asyncReq(doiPresents []refMeta, doiMisses []refMeta, c *Clean) ([]refMeta, []refMeta) {
	type result struct {
		entry refMeta
		err   error
	}
	mx1, mx2, toAdd := loadBalance(len(doiPresents), len(doiMisses))
	done := make(chan struct{})
	defer close(done)

	inputChP := streamInput(done, doiPresents)
	inputChM := streamInput(done, doiMisses)
	resultChP := make(chan result, len(doiPresents))
	resultChM := make(chan result, len(doiMisses))

	var wg sync.WaitGroup
	wg.Add(toAdd)

	for i := 0; i < mx1; i++ {
		go func() {
			for input := range inputChP {
				res, err := fetchRef(input, c)
				resultChP <- result{entry: res, err: err}
			}
			wg.Done()
		}()
	}
	for i := 0; i < mx2; i++ {
		go func() {
			for input := range inputChM {
				res, err := fetchMissing(input, c)
				resultChM <- result{entry: res, err: err}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(resultChP)
		close(resultChM)
	}()

	var finalRef = make([]refMeta, 0, len(doiPresents)+len(doiMisses))
	var oldRef = make([]refMeta, 0, len(doiMisses))
	for result := range resultChP {
		if result.err != nil {
			log.Printf("Error for doi %s:%v", result.entry.doi, result.err)
		}
		finalRef = append(finalRef, result.entry)
	}
	for result := range resultChM {
		if result.err != nil {
			log.Printf("Error for title %s:%v", result.entry.title, result.err)
		}
		if result.entry.ref != "" {
			finalRef = append(finalRef, result.entry)
			continue
		}
		oldRef = append(oldRef, result.entry)
	}
	return finalRef, oldRef
}

func fetchRef(doiP refMeta, c *Clean) (refMeta, error) {
	isValidDoi, strippedDoi := verifyDOI(doiP.doi)
	doiP.isValidDoi = true
	if !isValidDoi {
		verbosePrint(c.Verbose, fmt.Sprintf("Invalid DOI %s. Searching with name: %s", doiP.doi,
			doiP.title), os.Stdout)
		doiP.isValidDoi = false
	}
	r, err := request.RefDoi(strippedDoi, "bibtex")

	if err != nil {
		verbosePrint(c.Verbose, fmt.Sprintf("Cannot find with DOI %s. Searching with name: %s", doiP.doi,
			doiP.title), os.Stdout)
		doiP.isValidDoi = false
		return doiP, nil
	}
	b, _ := bibtex.Parse(strings.NewReader(r))
	doiP.keyNew = b.Entries[0].CiteName
	doiP.ref = r
	return doiP, err
}

func fetchMissing(doiM refMeta, c *Clean) (refMeta, error) {
	var (
		d    string
		err1 error
		err2 error
	)

	if doiM.isArxiv {
		d, err1 = request.DoiDataCite(doiM.title)
	}

	if !doiM.isArxiv || d == "" {
		d, err2 = request.DoiCrossRef(doiM.title)
	}

	if err1 != nil && err2 != nil {
		return doiM, fmt.Errorf("%w\n %v", err1, err2)
	}

	verbosePrint(c.Verbose, fmt.Sprintf("Found doi %s for title %s", d, doiM.title), os.Stdout)
	r, err := request.RefDoi(d, "bibtex")
	if err != nil {
		return doiM, fmt.Errorf("%w", err)
	}

	isEq := checkEq(r, doiM.title)
	if isEq {
		b, _ := bibtex.Parse(strings.NewReader(r))
		doiM.keyNew = b.Entries[0].CiteName
		doiM.ref = r
	}
	if r == "" || !isEq {
		doiM.ref = ""
		// doiM.keepOld = true
	}
	return doiM, err
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
	if n == "" {
		return false
	}
	var nt string
	var replace = []string{"{", "", "}", "", "(", "", ")", "",
		"\\textemdash", " ", "--", " ", "-", " ", "[", "", "]", "",
		"Accepted From Open Call", "", "$", "", "\\mathsemicolon", "",
		";", "", ":", "", ",", "", ".", "", "?", "", "!", ""}

	p, _ := bibtex.Parse(strings.NewReader(n))
	if len(p.Entries) == 0 {
		return false
	}
	o = strings.NewReplacer(replace...).Replace(o)
	o = strings.ReplaceAll(o, " ", "")
	if t, ok := p.Entries[0].Fields["title"]; ok {
		nt = strings.NewReplacer(replace...).Replace(t.String())
		nt = strings.ReplaceAll(nt, " ", "")
	}
	return strings.EqualFold(nt, o)
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
	rKeys := make([]string, 0, 2*len(key))
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
