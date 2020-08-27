package cmd

import (
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/lucasepe/tiles/tileset"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Use:                   "pull <tileset PATH or URL>",
	Example:               pullCmdExample(),
	Short:                 "Extracts the tile with the specified identifier from the tileset",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetString(optID)
		if err != nil {
			return err
		}

		ts, err := tileset.Load(args[0])
		if err != nil {
			return err
		}

		tile, ok := ts[0].Get(id)
		if !ok {
			return fmt.Errorf("tile with id = %s not found", id)
		}

		img, err := ts[0].Image(tile)
		if err != nil {
			return err
		}

		// TODO save image to stdout
		enc := png.Encoder{
			CompressionLevel: png.BestSpeed,
		}

		return enc.Encode(os.Stdout, img)
	},
}

func init() {
	pullCmd.Flags().String(optID, "", "the tile identifier in the specified tileset")
	pullCmd.MarkFlagRequired(optID)

	rootCmd.AddCommand(pullCmd)
}

func pullCmdExample() string {
	tpl := `  {{APP}}  pull --id aws_waf ../examples/aws_tileset.yml
  {{APP}}  pull --id aws_waf ../examples/aws_tileset.yml > aws_waf.png
  {{APP}}  pull --id aws_waf ../examples/aws_tileset.yml | viu -`

	return strings.Replace(tpl, "{{APP}}", appName(), -1)
}
