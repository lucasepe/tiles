package cmd

import (
	"os"
	"strings"

	"github.com/lucasepe/tiles/composer"
	"github.com/lucasepe/tiles/imagelist"
	"github.com/spf13/cobra"
)

// composeCmd represents the compose command
var composeCmd = &cobra.Command{
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Use:                   "compose <PNG_IMAGES_FOLDER>",
	Short:                 "Generate a tileset from all PNG images in the specified directory",
	Example:               composeCmdExample(),
	RunE: func(cmd *cobra.Command, args []string) error {
		images, err := imagelist.Load(args[0])
		if err != nil {
			return err
		}

		return composer.Do(images, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}

func composeCmdExample() string {
	tpl := `  {{APP}} compose /path/to/png/images/ > my_tileset.yml`
	return strings.Replace(tpl, "{{APP}}", appName(), -1)
}
