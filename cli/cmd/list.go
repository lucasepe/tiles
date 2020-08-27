package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/lucasepe/tiles/tileset"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Use:                   "list <tileset PATH or URL>",
	Short:                 "Lists all tiles identifiers contained in the specified tilset",
	Example:               listCmdExample(),
	RunE: func(cmd *cobra.Command, args []string) error {
		sets, err := tileset.Load(args[0])
		if err != nil {
			return err
		}

		for _, ts := range sets {
			for _, el := range ts.Tiles {
				if _, err := fmt.Fprintln(os.Stdout, el.ID); err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listCmdExample() string {
	tpl := `  {{APP}} list https://github.com/lucasepe/tiles/TBD
  {{APP}} list /path/to/tileset.yml`

	return strings.Replace(tpl, "{{APP}}", appName(), -1)
}
