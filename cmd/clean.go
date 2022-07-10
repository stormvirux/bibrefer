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
	"github.com/stormvirux/bibrefer/internal/ref"
	"os"

	"github.com/spf13/cobra"
)

var (
	texFile []string
	bibFile string
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean [flags] [query]",
	Short: "Cleans and updates the bibliography",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 && bibFile == "" {
			fmt.Printf("\n")
			return fmt.Errorf("expected either a bib file or a bib Entry, not both. Provide one.\n ")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		app := ref.Clean{BibKey: update, FullJournal: journal, FullAuthor: author,
			NoNewline: noNewline, Verbose: verbose, Output: output,
			TexFile: texFile, BibFile: bibFile}
		ref, err := app.Run(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			return err
		}
		if ref != "" {
			fmt.Printf("%s", ref)
		}
		return nil
	},
	Example: `bibrefer clean -b testBib.bib -t test.tex
bibrefer clean -b testBib.bib
bibrefer clean -vjk -b testBib.bib -t test.tex, test2.tex
bibrefer clean -vjk -b testBib.bib -t test.tex -t test2.tex
bibrefer clean -k "@Article{andrews2011tractable,
	Title                    = {A tractable approach to coverage and rate in cellular networks},
	Author                   = {G.J. Andrews and F. Baccelli and R.K. Ganti},
	Journal                  = {{IEEE} Trans. Commun.},
	Year                     = {2011},
	Month                    = {Nov.},
	Number                   = {11},
	Pages                    = {3122--3134},
	Volume                   = {59}
}"
`,
}

/*func cleanUsage() string {
	return `
Usage: {{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command] <query> {{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

Query:
  {{with $x := "DOI of the publication"}}{{$x}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}*/

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().StringSliceVarP(&texFile, "tex", "t",
		[]string{}, "tex file(s) to update")
	cleanCmd.Flags().StringVarP(&bibFile, "bib", "b",
		"", "bib file to update")
	cleanCmd.Flags().BoolVarP(&update, "keep-key", "k",
		false, "keep the bib entry key format from doi.org")
	cleanCmd.Flags().BoolVarP(&journal, "full-journal", "j",
		false, "use full  journal name in reference")
	cleanCmd.Flags().BoolVarP(&author, "full-author", "a",
		false, "use full author names in reference")
}
