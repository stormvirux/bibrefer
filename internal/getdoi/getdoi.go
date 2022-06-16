package getdoi

import (
	"fmt"
	"github.com/stormvirux/bibrefer/pkg/request"
	"os"
	"path/filepath"
	"strings"
)

type App struct{}

// TODO: Pass flags to App and initialize with it

func (a *App) Run(query []string, flags []bool) (string, error) {
	var host = "https://api.crossref.org/works"

	var (
		queryTxt       string
		extractedTitle string
		extractedDoi   string
		// arXivID        string
	)
	arxiv := flags[0]
	// clip := flags[1]
	verbose := flags[2]
	queryTxt = strings.Join(query, " ")

	if ext := filepath.Ext(queryTxt); ext == ".pdf" {
		var err error
		// _ is arXivID
		extractedTitle, extractedDoi, _, err = processPdf(queryTxt, verbose)
		if err != nil {
			return "", err
		}

		if verbose {
			fmt.Printf("\nDetected Title: %s\n", extractedTitle)
		}

		if extractedDoi != "" {
			verbosePrint(verbose, fmt.Sprintf("Detected DOI: %s\n", extractedDoi), os.Stdout)
			return extractedDoi, nil
		}
		if extractedTitle != "" {
			queryTxt = extractedTitle
		}
	}

	if arxiv {
		host = "https://api.datacite.org/dois"
		verbosePrint(verbose, "[Info] Retrieving data from DataCite", os.Stdout)
		extractedDoi, err := request.DoiDataCite(host, queryTxt)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		verbosePrint(verbose, fmt.Sprintf("\nDetected DOI: %s\n ", extractedDoi), os.Stdout)
		fmt.Println(extractedDoi)
		return extractedDoi, nil
	}
	verbosePrint(verbose, "[Info] Retrieving data from CrossRef", os.Stdout)
	extractedDoi, err := request.DoiCrossRef(host, queryTxt)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	verbosePrint(verbose, fmt.Sprintf("\nDetected DOI: %s\n ", extractedDoi), os.Stderr)
	fmt.Println(extractedDoi)
	return extractedDoi, nil
}
