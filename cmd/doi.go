/*
Copyright © 2022 Thaha Mohammed <thaha.mohammed@aalto.fi>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stormvirux/bibrefer/internal/getdoi"
	"os"
)

var (
	clip    bool
	verbose bool
	arxiv   bool
)

// doiCmd represents the doi command
var doiCmd = &cobra.Command{
	Use:   "doi [flags] <query>",
	Short: "Returns the publication DOI for a given publication",
	Long: `doi returns the publication DOI for a given publication from either a pdf file, CrossRef, or DataCite (for ArXiv).
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
		doi, err := app.Run(args, []bool{arxiv, clip, verbose})
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			return err
		}
		fmt.Println(doi)
		return nil
	},
}

/*const doiUsage = `doi returns the publication DOI for a given publication from either a pdf file, CrossRef, or DataCite (for ArXiv).

Usage:
  bibrefer getdoi [flags] <query>

Flags:
  -c, --clip       copy the DOI to clipboard
  -d, --datacite   retrieve the DOI from DataCite (for ArXiv)
  -h, --help       help for getdoi
  -V, --verbose    show verbose information


`*/

func doiUsage(cmd *cobra.Command) string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command] <query> {{end}}{{if gt (len .Aliases) 0}}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}
{{if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

Query:
  can consist of a publication title or a PDF file. For publication title (not file) from arXiv use -d/--datacite flag to fetch the doi from DataCite.{{if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

func init() {
	rootCmd.AddCommand(doiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getdoiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getdoiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	doiCmd.Flags().BoolVarP(&clip, "clip", "c", false,
		"copy the DOI to clipboard")
	doiCmd.Flags().BoolVarP(&verbose, "verbose", "V", false,
		"show verbose information")
	doiCmd.Flags().BoolVarP(&arxiv, "datacite", "d", false,
		"retrieve the DOI from DataCite (for ArXiv)")

	doiCmd.SetUsageTemplate(doiUsage(doiCmd))
	// doiCmd.SetHelpTemplate(doiCmd.UsageTemplate())
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
