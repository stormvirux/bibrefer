package getdoi

import (
	"fmt"
	"github.com/stormvirux/bibrefer/pkg/fix"
	"github.com/stormvirux/pdf"
	"os"
	"regexp"
	"strings"
)

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

func extractData(reader *pdf.Reader) (extractedTitle string, extractedDoi string, arXivID string, err error) {
	largestFont := 0.0
	largestFontIndex := 0

	pageNumber := 1
	p := reader.Page(pageNumber)
	if p.V.IsNull() {
		return "", "", "", nil
	}

	titleSize := 0
	titleSizeE := 0
	texts := p.Content().Text
	if len(texts) < 1 {
		return "", "", "", fmt.Errorf("[Warning !] unable to open the pdf")
	}

	//fmt.Println(texts)
	textForDoi := &strings.Builder{}

	titleFont := map[string]float64{
		"elsevier": 13.4495,
	}
	fontStartIndex := 0
	prevSize := 0.0
	var isElsevier bool
	for i := 0; i < 1000; i++ {
		textForDoi.WriteString(texts[i].S)
		isElsevier = strings.Contains(textForDoi.String(), "elsevier")
		if isElsevier && texts[i].FontSize == titleFont["elsevier"] && prevSize != titleFont["elsevier"] {
			prevSize = titleFont["elsevier"]
			fontStartIndex = i
			titleSizeE = 0
		}
		if isElsevier && titleFont["elsevier"] == texts[i].FontSize {
			titleSizeE++
		}
		if largestFont < texts[i].FontSize {
			largestFont = texts[i].FontSize
			largestFontIndex = i
			titleSize = 0
		}
		if largestFont == texts[i].FontSize {
			titleSize++
		}
	}

	t := make([]pdf.Text, len(texts))

	for i := 0; i < len(texts); i++ {
		t[i] = texts[i]
	}

	extractedDoi, arXivID = findDoiID(textForDoi, t)

	var title strings.Builder
	title.Grow(300)

	if isElsevier {
		largestFontIndex = fontStartIndex
		titleSize = titleSizeE
	}

	//title := &strings.Builder{}

	words := findWords(texts[largestFontIndex : titleSize+largestFontIndex])

	for _, word := range words {
		title.WriteString(word.S)
	}
	extractedTitle = strings.ReplaceAll(title.String(), "\n", " ")

	return extractedTitle, extractedDoi, arXivID, nil
}

func findDoiID(textForDoi *strings.Builder, textsC []pdf.Text) (extractedDoi string, arXivID string) {

	regexDoi := regexp.MustCompile(`(?i)\b10\.(\d+\.*)+/(([^\s.])+\.*)+\b`)
	for i := 1001; i < len(textsC); i++ {
		textForDoi.WriteString(textsC[i].S)
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
	if !isArXiv {
		/*if extractedDoi = regexDoi.FindString(textForDoi.String()); extractedDoi == "" {
			for i := 1001; i < len(textsC); i++ {
				textForDoi.WriteString(textsC[i].S)
			}
			extractedDoi = regexDoi.FindString(textForDoi.String())
		}*/
		allWord := findWords(textsC)
		temp := []string{}
		for _, t := range allWord {
			temp = append(temp, t.S)
		}

		extractedDoi = regexDoi.FindString(strings.Join(temp, " "))
	}
	return extractedDoi, arXivID
}

func processPdf(fileName string, verbose bool) (extractedTitle string, extractedDoi string, arXivID string, err error) {
	f, r, err := readPdf(fileName, verbose)

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if os.IsNotExist(err) {
		return "", "", "", err
	}

	pdf_ := fix.Pdf{Name: fileName}

	if err != nil {
		var msg = []string{"[Info] issue opening pdf -- %v\n" +
			"[Info] ghostScript installation detected.\n" +
			"[Info] Attempting to correct the pdf with ghostscript\\n",
			"pdf repair failed: %w", "pdf opening after repair failed: %w"}

		var isRecover bool
		isRecover, f, r, err = errorRecovery(verbose, pdf_, err, msg)

		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				verbosePrint(verbose, fmt.Sprintf("%s", err), os.Stderr)
			}
		}(pdf_.TmpDir)

		if !isRecover {
			return "", "", "", err
		}
	}

	extractedTitle, extractedDoi, arXivID, err = extractData(r)
	if err != nil && extractedTitle == "" {
		var msg = []string{"%v\n[Info] ghostScript installation detected.\n[Info] Attempting to correct the pdf with ghostscript", "%w", "pdf opening failed: %w"}
		var isRecover bool
		isRecover, f, r, err = errorRecovery(verbose, pdf_, err, msg)
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Printf("%s", err)
			}
		}(pdf_.TmpDir)

		if !isRecover {
			return "", "", "", err
		}
		extractedTitle, extractedDoi, arXivID, err = extractData(r)
		if err != nil {
			return "", "", "", fmt.Errorf("%w", err)
		}
	}

	return extractedTitle, extractedDoi, arXivID, nil
}

func errorRecovery(verbose bool, pdf_ fix.Pdf, err error, msg []string) (bool, *os.File, *pdf.Reader, error) {
	gsIsPresent, gsPath := detectGS()
	if !gsIsPresent {
		return false, &os.File{}, &pdf.Reader{}, fmt.Errorf("ghostScript not found. Check if bin and lib is in PATH")
	}

	verbosePrint(verbose, fmt.Sprintf(msg[0], err), os.Stderr)
	err = pdf_.Fix(gsPath)
	if err != nil {
		return false, &os.File{}, &pdf.Reader{}, fmt.Errorf(msg[1], err)
	}

	f, r, err := readPdf(pdf_.TmpFile, verbose)
	if err != nil {
		return false, &os.File{}, &pdf.Reader{}, fmt.Errorf(msg[2], err)
	}

	return true, f, r, nil
}
