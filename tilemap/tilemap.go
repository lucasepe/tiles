package tilemap

import (
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"

	"github.com/lucasepe/tiles/data"
	"github.com/lucasepe/tiles/grid"
	"github.com/lucasepe/tiles/tileset"
	"gopkg.in/yaml.v2"
)

type TileMap struct {
	cols      int
	rows      int
	tileSize  int
	layout    []int
	margin    int
	watermark string
	bgColor   string
	mapping   map[int]string
	atlasList []string
}

func (tm *TileMap) Render(wr io.Writer) error {
	repo, err := tileset.Load(tm.atlasList...)
	if err != nil {
		return err
	}

	gr, err := grid.NewGrid(tm.rows, tm.cols, tm.tileSize,
		grid.Background(tm.bgColor),
		grid.Margin(tm.margin),
		grid.Watermark(tm.watermark))
	if err != nil {
		return err
	}

	gr.DrawBorder()

	for c := 0; c < tm.cols; c++ {
		for r := 0; r < tm.rows; r++ {
			// Grab the tile index
			pos := r*tm.cols + c
			if pos >= len(tm.layout) {
				return fmt.Errorf("invalid index [%d] with a grid length of %d", pos, len(tm.layout))
			}

			idx := tm.layout[pos]
			if idx <= 0 {
				continue
			}

			// Find the image for the tile id
			id, ok := tm.mapping[idx]
			if !ok {
				return fmt.Errorf("tile with index: %d not found in mapping", idx)
			}

			img, err := findImageByID(repo, id)
			if err != nil {
				return err
			}

			gr.DrawImage(img, r, c)
		}
	}

	gr.DrawWatermark()

	return gr.EncodePNG(wr)
}

// UnmarshalYAML implements the Unmarshaler interface of the yaml pkg.
func (tm *TileMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	aux := struct {
		Cols      int            `yaml:"cols"`
		Rows      int            `yaml:"rows"`
		TileSize  int            `yaml:"tile_size"`
		Margin    int            `yaml:"margin"`
		BgColor   string         `yaml:"bg_color"`
		Layout    string         `yaml:"layout"`
		Watermark string         `yaml:"watermark"`
		Mapping   map[int]string `yaml:"mapping"`
		AtlasList []string       `yaml:"atlas_list"`
	}{}

	err := unmarshal(&aux)
	if err != nil {
		return err
	}

	tm.cols = aux.Cols
	tm.rows = aux.Rows
	tm.tileSize = aux.TileSize
	tm.margin = aux.Margin
	tm.bgColor = aux.BgColor
	tm.watermark = aux.Watermark
	tm.mapping = make(map[int]string)
	for k, v := range aux.Mapping {
		tm.mapping[k] = v
	}

	tm.atlasList = make([]string, len(aux.AtlasList))
	for i, uri := range aux.AtlasList {
		tm.atlasList[i] = uri
	}

	layout := strings.Split(strings.Replace(aux.Layout, " ", ",", -1), ",")

	tm.layout = make([]int, len(layout))
	for i := 0; i < len(layout); i++ {
		num, err := strconv.Atoi(strings.TrimSpace(layout[i]))
		if err != nil {
			return err
		}
		tm.layout[i] = num
	}

	return nil
}

func Load(uri string) (TileMap, error) {
	dat, err := data.Fetch(uri, -1)
	if err != nil {
		return TileMap{}, err
	}

	res := TileMap{}
	if err := yaml.Unmarshal(dat, &res); err != nil {
		return TileMap{}, err
	}

	return res, nil
}

func findImageByID(repo []*tileset.Tileset, id string) (image.Image, error) {
	for _, el := range repo {
		if tile, ok := el.Get(id); ok {
			return el.Image(tile)
		}
	}

	return nil, fmt.Errorf("tile with id: %s not found", id)
}
