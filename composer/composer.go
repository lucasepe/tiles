package composer

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lucasepe/tiles/binpack"
	"github.com/lucasepe/tiles/data"
	"github.com/lucasepe/tiles/tileset"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Do generates a tileset from the image
// list and print the result to the specified writer.
func Do(il []string, wr io.Writer) error {
	items, err := decodeImageList(il)
	if err != nil {
		return err
	}
	sort.Sort(byMaxOfWidthAndHeight(items))

	bl := blockList{blocks: items}
	bl.width, bl.height = binpack.Pack(&bl)

	if err := bl.createPNG(); err != nil {
		return err
	}

	return bl.dump(wr)
}

// block holds tile position,
// dimensions and source image.
type block struct {
	src  string
	x, y int
	w, h int
	id   string
}

type blockList struct {
	blocks        []*block
	width, height int
	data          []byte
}

// createImage assembles all tiles in a
// bigger bin-packed image.
func (bl *blockList) createPNG() error {
	sheet := image.NewRGBA(image.Rect(0, 0, bl.width, bl.height))

	for _, el := range bl.blocks {
		fp, err := os.Open(el.src)
		if err != nil {
			return err
		}
		defer fp.Close()

		img, _, err := image.Decode(fp)
		if err != nil {
			return err
		}

		r := image.Rect(el.x, el.y, el.x+el.w, el.y+el.h)
		draw.Draw(sheet, r, img, image.Point{}, draw.Over)
	}

	buf := new(bytes.Buffer)
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(buf, sheet); err != nil {
		return err
	}

	bl.data = buf.Bytes()

	return nil
}

func (bl *blockList) dump(wr io.Writer) error {

	res := tileset.Tileset{
		Width:  bl.width,
		Height: bl.height,
		Tiles:  make([]*tileset.Tile, len(bl.blocks)),
		Data:   data.Wrap(base64.StdEncoding.EncodeToString(bl.data), 76),
	}

	for i, el := range bl.blocks {
		res.Tiles[i] = &tileset.Tile{
			MinX: el.x, MinY: el.y,
			MaxX: el.x + el.w, MaxY: el.y + el.h,
			ID: el.id,
		}
	}

	dat, err := yaml.Marshal(&res)
	if err != nil {
		return err
	}
	_, err = wr.Write(dat)
	return err
}

// Len returns the number of blocks in total.
func (bl *blockList) Len() int {
	return len(bl.blocks)
}

// Size returns the width and height of the block n.
func (bl *blockList) Size(n int) (width, height int) {
	el := bl.blocks[n]
	return el.w, el.h
}

// Place places the block n, at the position [x, y].
func (bl *blockList) Place(n, x, y int) {
	el := bl.blocks[n]
	el.x = x
	el.y = y
}

// byMaxOfWidthAndHeight implements sort.Interface based on the max(width, height).
type byMaxOfWidthAndHeight []*block

func (a byMaxOfWidthAndHeight) Len() int { return len(a) }
func (a byMaxOfWidthAndHeight) Less(i, j int) bool {
	m1 := maxInt(a[i].w, a[i].h)
	m2 := maxInt(a[j].w, a[j].h)
	return m1 < m2
}
func (a byMaxOfWidthAndHeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// decodeImageList load images from the specified list
// and decodes the dimensions of each image.
// Returns a map which keys are the source image file
// and values are tiles with derived ID and decoded dimensions.
func decodeImageList(list []string) ([]*block, error) {
	makeID := func(filename string) string {
		base := filepath.Base(filename)
		ext := filepath.Ext(filename)
		return strings.TrimSuffix(base, ext)
	}

	res := []*block{}

	for _, el := range list {
		r, err := os.Open(el)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		im, _, err := image.DecodeConfig(r)
		if err != nil {
			return nil, errors.Wrapf(err, "image <%s>", el)
		}

		res = append(res, &block{
			id:  makeID(el),
			src: el,
			w:   im.Width,
			h:   im.Height,
		})
	}

	return res, nil
}

func maxInt(a, b int) int {
	if b > a {
		return b
	}
	return a
}
