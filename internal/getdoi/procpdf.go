package getdoi

import (
	"fmt"
	"github.com/stormvirux/bibrefer/pkg/fix"
	"github.com/stormvirux/pdf"
	"os"
	"regexp"
	"strings"
	"sync"
)

type Result struct {
	extractedDoi string
	arXivID      string
}

var wg = sync.WaitGroup{}

func readPdf(fileName string, verbose bool) (*os.File, *pdf.Reader, error) {
	f, r, err := pdf.Open(fileName)
	if err != nil {
		if verbose {
			fmt.Printf("[Warning !] Could not open the pdf %s\n", fileName)
		}
		return f, r, err
	}
	verbosePrint(verbose, fmt.Sprintf("[Info] Reading file: %s", fileName), os.Stdout)
	return f, r, err
}

func extractData(reader *pdf.Reader, tmpFile string) (string, string, string, error) {
	txtPdf, err := os.ReadFile(strings.ReplaceAll(tmpFile, ".pdf", ".txt"))
	if err != nil {
		return "", "", "", err
	}

	pageNumber := 1
	p := reader.Page(pageNumber)
	if p.V.IsNull() {
		return "", "", "", nil
	}

	textsMeta := p.Content().Text
	// textsMeta2 := p2.Content().Text
	// textsMeta = append(textsMeta, textsMeta2...)
	if len(textsMeta) < 1 {
		return "", "", "", fmt.Errorf("[Warning !] unable to open the pdf")
	}

	t := make([]pdf.Text, len(textsMeta))
	copy(t, textsMeta)

	wg.Add(2)
	resulTitle := make(chan string)
	go findTitle(textsMeta, txtPdf, resulTitle)
	extractedTitle := <-resulTitle

	resultDoi := make(chan Result)
	go findDoiID(t, txtPdf, resultDoi)
	result := <-resultDoi
	wg.Wait()
	return extractedTitle, result.extractedDoi, result.arXivID, nil
}

func findTitle(textMeta []pdf.Text, txtPdf []byte, resulTitle chan string) {
	defer close(resulTitle)
	// largestFont := 0.0
	largestIndex := 0
	nextLargestIndex := 0

	largestSize := 0
	nextLargestSize := 0

	var isElsevier bool
	var isMdpi bool

	/*titleFont := map[string]map[string]float64{
		"elsevier": {"max": 13.5, "min": 13.4},
		"mdpi":     {"max": 17.00, "min": 17.933},
	}*/

	isElsevier = strings.Contains(string(txtPdf), "elsevier")
	isMdpi = strings.Contains(string(txtPdf), "mdpi")

	largestIndex, largestSize, nextLargestIndex, nextLargestSize = index(textMeta)

	/*	fontStartIndex := 0
		prevSize := 0.0*/

	/*	if isElsevier || isMdpi {
		if isElsevier && textMeta[i].FontSize >= titleFont["elsevier"]["min"] && textMeta[i].
			FontSize <= titleFont["elsevier"]["max"] && prevSize != texts[i].FontSize {
			prevSize = texts[i].FontSize
			fontStartIndex = i
			titleSizeE = 0
		}
		if isElsevier && texts[i].FontSize >= titleFont["elsevier"]["min"] && texts[i].FontSize <= titleFont["elsevier"]["max"] {
			titleSizeE++
		}
	}*/

	var title strings.Builder
	title.Grow(300)
	var words []pdf.Text

	words = findWords(textMeta[largestIndex : largestSize+largestIndex])

	if isElsevier || isMdpi {
		words = findWords(textMeta[nextLargestIndex : nextLargestSize+nextLargestIndex])
	}

	for _, word := range words {
		title.WriteString(word.S)
	}
	extractedTitle := strings.ReplaceAll(title.String(), "\n", " ")

	resulTitle <- extractedTitle
	wg.Done()
}

func index(textMeta []pdf.Text) (largestIndex int, largestSize int, nextLargestIndex int, nextLargestSize int) {
	largest := 0.0
	largestIndex = 0
	largestSize = 0

	nextLargest := 0.0
	nextLargestIndex = 0
	nextLargestSize = 0

	for i := 0; i < 1000; i++ {
		if largest < textMeta[i].FontSize {
			nextLargest = largest
			nextLargestIndex = largestIndex
			nextLargestSize = largestSize
			largest = textMeta[i].FontSize
			largestIndex = i
			largestSize = 0
		} else if textMeta[i].FontSize > nextLargest && textMeta[i].FontSize < largest {
			nextLargest = textMeta[i].FontSize
			nextLargestIndex = i
			nextLargestSize = 0
		}
		if largest == textMeta[i].FontSize {
			largestSize++
		}
		if nextLargest == textMeta[i].FontSize {
			nextLargestSize++
		}
	}
	return largestIndex, largestSize, nextLargestIndex, nextLargestSize
}

func findDoiID(textsC []pdf.Text, txtPdf []byte, resultDoi chan Result) {
	defer close(resultDoi)
	regexDoi := regexp.MustCompile(`(?i)\b10\.(\d+\.*)+/(([^\s.])+\.*)+\b`)

	textForDoi := &strings.Builder{}
	for i := 0; i < len(textsC); i++ {
		textForDoi.WriteString(textsC[i].S)
	}

	var extractedDoi string
	var arXivID string
	/*for i := 0; i < len(textsC); i++ {
		textForDoi.WriteString(textsC[i].S)
	}*/

	isElsevier := strings.Contains(string(txtPdf), "elsevier")
	isMdpi := strings.Contains(string(txtPdf), "mdpi")

	if isElsevier || isMdpi {
		extractedDoi = regexDoi.FindString(string(txtPdf))
	}

	isArXiv := strings.Contains(textForDoi.String(), "arXiv")
	if isArXiv {
		var regexDoiOld *regexp.Regexp
		regexDoi, regexDoiOld = arxivRegexBuilder()
		extractedDoi = regexDoi.FindString(textForDoi.String())
		if extractedDoi == "" {
			extractedDoi = regexDoiOld.FindString(textForDoi.String())
		}
		arXivID = extractedDoi
		extractedDoi = strings.ReplaceAll(extractedDoi, "arXiv:", "")
		extractedDoi = "10.48550/arXiv." + stripVersion(extractedDoi)
	}
	if !isArXiv && !isElsevier && !isMdpi {
		/*if extractedDoi = regexDoi.FindString(textForDoi.String()); extractedDoi == "" {
			for i := 1001; i < len(textsC); i++ {
				textForDoi.WriteString(textsC[i].S)
			}
			extractedDoi = regexDoi.FindString(textForDoi.String())
		}*/
		allWord := findWords(textsC)
		var temp []string
		for _, t := range allWord {
			temp = append(temp, t.S)
		}

		extractedDoi = regexDoi.FindString(strings.Join(temp, " "))
	}
	res := new(Result)
	res.arXivID = arXivID
	res.extractedDoi = extractedDoi

	resultDoi <- *res
	wg.Done()
}

func processPdf(fileName string, verbose bool) (extractedTitle string, extractedDoi string, arXivID string, err error) {
	f, r, err := readPdf(fileName, verbose)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	if os.IsNotExist(err) {
		return "", "", "", err
	}

	pdF := fix.Pdf{Name: fileName}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			verbosePrint(verbose, fmt.Sprintf("%s", err), os.Stderr)
		}
	}(pdF.TmpDir)

	if err != nil {
		var msg = []string{"[Info] issue opening pdf -- %v\n" +
			"[Info] ghostScript installation detected.\n" +
			"[Info] Attempting to correct the pdf with ghostscript\\n",
			"pdf repair failed: %w", "pdf opening after repair failed: %w"}

		var isRecover bool
		isRecover, f, r, pdF, err = errorRecovery(verbose, pdF, err, msg)
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		if !isRecover {
			return "", "", "", err
		}
	}

	extractedTitle, extractedDoi, arXivID, err = extractData(r, pdF.TmpFile)
	if err != nil && extractedTitle == "" {
		var msg = []string{"%v\n[Info] ghostScript installation detected.\n[Info] Attempting to correct the pdf with ghostscript", "%w", "pdf opening failed: %w"}
		var isRecover bool
		isRecover, f, r, pdF, err = errorRecovery(verbose, pdF, err, msg)
		/*defer func(f *os.File) {
			_ = f.Close()
		}(f)*/

		if !isRecover {
			return "", "", "", err
		}
		extractedTitle, extractedDoi, arXivID, err = extractData(r, pdF.TmpFile)
		if err != nil {
			return "", "", "", fmt.Errorf("%w", err)
		}
	}

	return extractedTitle, extractedDoi, arXivID, nil
}

func errorRecovery(verbose bool, pdF fix.Pdf, err error, msg []string) (bool, *os.File, *pdf.Reader, fix.Pdf, error) {
	gsIsPresent, gsPath := detectGS()
	if !gsIsPresent {
		return false, &os.File{}, &pdf.Reader{}, pdF, fmt.Errorf("ghostScript not found. Check if bin and lib is in PATH")
	}

	verbosePrint(verbose, fmt.Sprintf(msg[0], err), os.Stderr)
	err = pdF.Fix(gsPath)
	if err != nil {
		return false, &os.File{}, &pdf.Reader{}, pdF, fmt.Errorf(msg[1], err)
	}

	f, r, err := readPdf(pdF.TmpFile, verbose)
	if err != nil {
		return false, &os.File{}, &pdf.Reader{}, pdF, fmt.Errorf(msg[2], err)
	}

	return true, f, r, pdF, nil
}
