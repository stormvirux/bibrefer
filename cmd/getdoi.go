/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stormvirux/bibrefer/internal/getdoi"
)

var (
	clip    bool
	verbose bool
	arxiv   bool
)

// getdoiCmd represents the getdoi command
var getdoiCmd = &cobra.Command{
	Use:   "getdoi [flags] query",
	Short: "Returns the publication DOI for a given publication",
	Long: `getdoi returns the publication DOI for a given publication from either a pdf file, CrossRef, or DataCite (for ArXiv).
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			fmt.Printf("\n")
			return fmt.Errorf("missing query. Expected atleast one query. Enter the file name or title of publication\n ")
		}
		return nil
	},
	//SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := getdoi.App{}
		if err := app.Run(args, []bool{arxiv, clip, verbose}); err != nil {
			// _, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			return err
		}
		return nil
	},
}

const doiUsage = `getdoi returns the publication DOI for a given publication from either a pdf file, CrossRef, or DataCite (for ArXiv).

Usage:
  bibrefer getdoi [flags] <query>

Flags:
  -c, --clip       copy the DOI to clipboard
  -d, --datacite   retrieve the DOI from DataCite (for ArXiv)
  -h, --help       help for getdoi
  -V, --verbose    show verbose information

Query: 
  can consist of a publication title or a PDF file. For publication title (not file) from arXiv use -d/--datacite flag to fetch the doi from DataCite.
`

func init() {
	rootCmd.AddCommand(getdoiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getdoiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getdoiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	getdoiCmd.Flags().BoolVarP(&clip, "clip", "c", false,
		"copy the DOI to clipboard")
	getdoiCmd.Flags().BoolVarP(&verbose, "verbose", "V", false,
		"show verbose information")
	getdoiCmd.Flags().BoolVarP(&arxiv, "datacite", "d", false,
		"retrieve the DOI from DataCite (for ArXiv)")

	// getdoiCmd.SetUsageTemplate(doiUsage)
	// getdoiCmd.SetHelpTemplate(getdoiCmd.UsageTemplate())
}

/*func getUserInput() ([]string, error) {
	// Use bufferio to get the input
	var query string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s\n", doiUsage)
		fmt.Printf("Provide a single file name or a publication name here:\n")
		scanner.Scan()
		query := scanner.Text()
		if query != "" {
			break
		}
	}
	return []string{query}, nil
}*/

// TODO: Update template from Cobra go template
