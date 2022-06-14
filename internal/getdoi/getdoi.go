package getdoi

import (
	"fmt"
	"path/filepath"
	"strings"
)

type App struct{}

// TODO: Write tests

func (a *App) Run(query []string, flags []bool) error {
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
			return err
		}

		if verbose {
			fmt.Printf("\nDetected Title: %s\n", extractedTitle)
		}
		// TODO: Maybe use Log.Print

		if extractedDoi != "" {
			if verbose {
				fmt.Printf("Detected DOI: %s\n", extractedDoi)
				return nil
			}
			fmt.Println(extractedDoi)
			return nil
		}
		if extractedTitle != "" {
			queryTxt = extractedTitle
		}
	}

	if arxiv {
		host = "https://api.datacite.org/dois"
		verbosePrint(verbose, "[Info] Retrieving data from DataCite")
		extractedDoi, err := doiDataCite(host, queryTxt)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if verbose {
			fmt.Printf("\nDetected DOI: %s\n ", extractedDoi)
			return nil
		}
		fmt.Println(extractedDoi)
		return nil
	}
	verbosePrint(verbose, "[Info] Retrieving data from CrossRef")
	extractedDoi, err := doiCrossRef(host, queryTxt)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Printf("\n")
	if verbose {
		fmt.Printf("\nDetected DOI: %s\n ", extractedDoi)
		return nil
	}

	fmt.Println(extractedDoi)
	return nil
}
