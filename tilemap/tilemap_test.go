package tilemap

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/lucasepe/tiles/tileset"
)

func TestFetchFromURI(t *testing.T) {
	tm, err := Load("../examples/demo.yml")
	if err != nil {
		t.Fatal(err)
	}

	ts, err := tileset.Load(tm.atlasList...)
	if err != nil {
		t.Fatal(err)
	}

	img, err := tm.Render(ts)
	if err != nil {
		t.Fatal(err)
	}

	savePNG(img, "delme.png")
}

func savePNG(im image.Image, filename string) error {
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	enc := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}

	return enc.Encode(fp, im)
}
