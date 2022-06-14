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
		// TODO: log here
		return f, r, err
	}
	verbosePrint(verbose, fmt.Sprintf("[Info] Reading file: %s", fileName))
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
	texts := p.Content().Text
	if len(texts) < 1 {
		return "", "", "", fmt.Errorf("[Warning !] unable to open the pdf")
	}

	textForDoi := &strings.Builder{}

	for i := 0; i < 1000; i++ {
		textForDoi.WriteString(texts[i].S)
		if largestFont < texts[i].FontSize {
			largestFont = texts[i].FontSize
			largestFontIndex = i
			titleSize = 0
		}
		if largestFont == texts[i].FontSize {
			titleSize++
		}
	}

	extractedDoi, arXivID = findDoiID(textForDoi, texts)

	var title strings.Builder
	title.Grow(300)
	//title := &strings.Builder{}

	words := findWords(texts[largestFontIndex : titleSize+largestFontIndex])
	for _, word := range words {
		title.WriteString(word.S)
	}
	extractedTitle = strings.ReplaceAll(title.String(), "\n", " ")

	return extractedTitle, extractedDoi, arXivID, nil
}

func findDoiID(textForDoi *strings.Builder, texts []pdf.Text) (extractedDoi string, arXivID string) {
	regexDoi := regexp.MustCompile(`(?i)\b10\.(\d+\.*)+/(([^\s.])+\.*)+\b`)
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
		if extractedDoi = regexDoi.FindString(textForDoi.String()); extractedDoi == "" {
			for i := 1001; i < len(texts); i++ {
				textForDoi.WriteString(texts[i].S)
			}
			extractedDoi = regexDoi.FindString(textForDoi.String())
		}
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
				fmt.Printf("%w", err)
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

	verbosePrint(verbose, fmt.Sprintf(msg[0], err))
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
