package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	banner = ` 
 ______  ____  _        ___   _____ 
|      ||    || |      /  _] / ___/ 
|      | |  | | |     /  [_ (   \_ 
|_|  |_| |  | | |___ |    _] \__  |
  |  |   |  | |     ||   [_  /  \ |  
  |  |   |  | |     ||     | \    |  
  |__|  |____||_____||_____|  \___| `

	appSummary = "Create and inspect a tile set from multiple PNG images."

	optID = "id"
)

// rootCmd represents the base command when called without any subcommands
var (
	version = "dev"
	rootCmd = &cobra.Command{
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Version:               version,
		Use:                   fmt.Sprintf("%s <COMMAND>", appName()),
		Short:                 appSummary,
		Long:                  banner,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}} - Luca Sepe <luca.sepe@gmail.com>
`)
}

func appName() string {
	return filepath.Base(os.Args[0])
}
