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

// refCmd represents the getref command
var refCmd = &cobra.Command{
	Use:   "ref [flags] <query>",
	Short: "Returns the reference of the given DOI, name, or pdf file",
	Long: `ref returns the the reference of a given DOI, name, or pdf file
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			fmt.Printf("\n")
			return fmt.Errorf("expected exactly one query. Provide a doi\n ")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		app := ref.Ref{BibKey: update, FullJournal: journal, FullAuthor: author, NoNewline: noNewline,
			Verbose: verbose, Output: output}
		ref, err := app.Run(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			return err
		}
		fmt.Printf("%s", ref)
		return nil
	},
	Example: `bibrefer ref -a 10.1016/j.jnca.2016.10.019
  bibrefer ref -o json 10.1016/j.jnca.2016.10.019
  bibrefer ref -ujaV 10.1016/j.jnca.2016.10.019
`,
}

var (
	update    bool
	journal   bool
	author    bool
	noNewline bool
	output    string
)

func refUsage() string {
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
}

func init() {
	rootCmd.AddCommand(refCmd)

	refCmd.Flags().BoolVarP(&update, "keep-key", "k", false,
		"keep the bib entry key format from doi.org")
	refCmd.Flags().BoolVarP(&journal, "full-journal", "j", false,
		"use full  journal name in reference")
	refCmd.Flags().BoolVarP(&author, "full-author", "a", false,
		"use full author names in reference")
	refCmd.Flags().BoolVarP(&noNewline, "no-newline", "n", false,
		"suppress trailing newline but prepend with newline")
	refCmd.Flags().StringVarP(&output, "output", "o", "bibtex",
		"sets the output format. Supported values are json, bibtex, and rdf-xml.")
	refCmd.SetUsageTemplate(refUsage())
}
