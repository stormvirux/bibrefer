package doi

import (
	"fmt"
	"github.com/stormvirux/bibrefer/pkg/request"
	"os"
	"path/filepath"
	"strings"
)

type Doi struct {
	Arxiv   bool
	Clip    bool
	Verbose bool
}

func (a *Doi) Run(query []string) (string, error) {
	var (
		queryTxt       string
		extractedTitle string
		extractedDoi   string
		// arXivID        string
	)
	// a.Arxiv := flags[0]
	// clip := flags[1]
	// a.Verbose := flags[2]
	queryTxt = strings.Join(query, " ")

	if ext := filepath.Ext(queryTxt); ext == ".pdf" {
		var err error
		// _ is arXivID
		extractedTitle, extractedDoi, _, err = processPdf(queryTxt, a.Verbose)
		if err != nil {
			return "", err
		}

		if a.Verbose {
			fmt.Printf("\nDetected Title: %s\n", extractedTitle)
		}

		if extractedDoi != "" {
			verbosePrint(a.Verbose, fmt.Sprintf("Detected DOI: %s\n", extractedDoi), os.Stdout)
			return extractedDoi, nil
		}
		if extractedTitle != "" {
			queryTxt = extractedTitle
		}
	}

	if a.Arxiv {
		verbosePrint(a.Verbose, "[Info] Retrieving data from DataCite", os.Stdout)
		extractedDoi, err := request.DoiDataCite(queryTxt)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		verbosePrint(a.Verbose, fmt.Sprintf("\nDetected DOI: %s\n ", extractedDoi), os.Stdout)
		// fmt.Println(extractedDoi)
		return extractedDoi, nil
	}
	verbosePrint(a.Verbose, "[Info] Retrieving data from CrossRef", os.Stdout)
	extractedDoi, err := request.DoiCrossRef(queryTxt)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	verbosePrint(a.Verbose, fmt.Sprintf("\nDetected DOI: %s\n ", extractedDoi), os.Stderr)
	// fmt.Println(extractedDoi)
	return extractedDoi, nil
}
