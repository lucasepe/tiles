package cmd

import (
	"os"
	"strings"

	"github.com/lucasepe/tiles/tilemap"
	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Use:                   "render <tilemap URL or PATH>",
	Short:                 "Render a square static tilemap",
	Example:               renderCmdExample(),
	RunE: func(cmd *cobra.Command, args []string) error {
		tm, err := tilemap.Load(args[0])
		if err != nil {
			return err
		}

		return tm.Render(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}

func renderCmdExample() string {
	tpl := `  {{APP}} render https://github.com/lucasepe/tiles/examples/ark.yml
  {{APP}} render /path/to/my_map.yml
  {{APP}} render /path/to/my_map.yml | viu -`

	return strings.Replace(tpl, "{{APP}}", appName(), -1)
}
