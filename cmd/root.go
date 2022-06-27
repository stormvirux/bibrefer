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
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     `bibrefer`,
	Version: version,
	Short:   "Search and retrieve doi or references for research articles",
	Long: `Bibrefer is a CLI application that can retrieve DOI and reference of a scholarly publication.
The DOI is obtained either from a pdf or the name of the article. 
The reference is fetched with DOI, pdf, or name of the article.  
`,

	Run: func(cmd *cobra.Command, args []string) {},
	Example: `  bibrefer ref 10.1016/j.jnca.2016.10.019
  bibrefer doi An Analysis of fault detection strategies in wsn
  bibrefer doi article1.pdf
`,
}

var (
	License bool
	verbose bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func rootUsage() string {
	return `
Usage: {{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command] {{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

func helpUsage() string {
	return `  ____  _ _     ____       __            
 | __ )(_) |__ |  _ \ ___ / _| ___ _ __  
 |  _ \| | '_ \| |_) / _ \ |_ / _ \ '__| 
 | |_) | | |_) |  _ <  __/  _|  __/ |    
 |____/|_|_.__/|_| \_\___|_|  \___|_|    

{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}
{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bibrefer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.Flags().BoolVarP(&License, "license", "l", false,
		"show license information")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false,
		"show verbose information")
	rootCmd.SetUsageTemplate(rootUsage())
	rootCmd.SetHelpTemplate(helpUsage())
}
