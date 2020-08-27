package tileset

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png" // load the PNG driver
	"io"
	"strings"
	"time"

	"github.com/lucasepe/tiles/cache"
	"github.com/lucasepe/tiles/data"
	"gopkg.in/yaml.v2"
)

// Tile describes a tile in the tiles set.
type Tile struct {
	ID   string `yaml:"id"`
	MinX int    `yaml:"minX"`
	MinY int    `yaml:"minY"`
	MaxX int    `yaml:"maxX"`
	MaxY int    `yaml:"maxY"`
}

// Rect returns the image rectangle fot the tile.
func (t *Tile) Rect() image.Rectangle {
	return image.Rect(t.MinX, t.MinY, t.MaxX, t.MaxY)
}

// Tileset describes a tile set.
type Tileset struct {
	Tiles  []*Tile `yaml:"tiles,omitempty"`
	Width  int     `yaml:"width"`
	Height int     `yaml:"height"`
	Data   string  `yaml:"data"`

	uri string
}

// Load fetches an array of tileset(s).
func Load(uri ...string) ([]*Tileset, error) {
	res := make([]*Tileset, len(uri))
	for i, u := range uri {
		el, err := loadOne(u)
		if err != nil {
			return nil, err
		}

		res[i] = el
	}

	return res, nil
}

// loadOne fetches a single tile set from the specified uri.
func loadOne(uri string) (*Tileset, error) {
	dat, err := data.Fetch(uri, -1)
	if err != nil {
		return nil, err
	}

	res := &Tileset{uri: uri}
	if err := yaml.Unmarshal(dat, &res); err != nil {
		return nil, err
	}
	res.Data = strings.Replace(res.Data, "\n", "", -1)

	return res, err
}

// Get returns the tile with the specified id.
func (ts *Tileset) Get(id string) (Tile, bool) {
	for _, el := range ts.Tiles {
		if strings.EqualFold(el.ID, id) {
			return Tile{
				ID:   el.ID,
				MinX: el.MinX, MinY: el.MinY,
				MaxX: el.MaxX, MaxY: el.MaxY,
			}, true
		}
	}

	return Tile{}, false
}

// List writers all the tiles id on the specified writer.
func (ts *Tileset) List(wr io.Writer) error {
	for _, el := range ts.Tiles {
		if _, err := fmt.Fprintln(wr, el.ID); err != nil {
			return err
		}
	}

	return nil
}

// Image returns the tile image.
func (ts *Tileset) Image(tile Tile) (image.Image, error) {
	img, err := ts.cachedImage()
	if err != nil {
		return nil, err
	}

	return img.(subImager).SubImage(tile.Rect()), nil
}

func (ts *Tileset) cachedImage() (image.Image, error) {
	var res image.Image
	if el, found := storage.Get(ts.uri); found {
		res = el.(image.Image)
		return res, nil
	}

	data, err := base64.StdEncoding.DecodeString(ts.Data)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	storage.Set(ts.uri, img, cache.DefaultExpiration)

	return img, nil
}

// subImager interface to return
// an image representing the portion of
// the tileset image visible through r.
type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

var storage *cache.Cache

func init() {
	storage = cache.New(5*time.Minute, 10*time.Minute)
}
