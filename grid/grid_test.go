package grid

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
)

func TestNewGrid(t *testing.T) {
	grid, err := NewGrid(12, 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 12, grid.rows)
	assert.Equal(t, 10, grid.cols)
	assert.Equal(t, 64, grid.cellSize)

	assert.NotNil(t, grid.Context())
	assert.Equal(t, float64(64), grid.CellSize())
}

func TestGridCellCenter(t *testing.T) {
	tests := []struct {
		row  int
		col  int
		want gg.Point
	}{
		{2, 2, gg.Point{X: 160, Y: 160}},
		{5, 4, gg.Point{X: 288, Y: 352}},
		{6, 7, gg.Point{X: 480, Y: 416}},
		{9, 3, gg.Point{X: 224, Y: 608}},
	}

	grid, err := NewGrid(12, 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := grid.CellCenter(tt.row, tt.col)
			if got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func TestGridVerifyInBounds(t *testing.T) {
	tests := []struct {
		row  int
		col  int
		want string
	}{
		{2, 2, ""},
		{5, 4, "cell (5, 4) is out of bounds"},
		{6, 7, "cell (6, 7) is out of bounds"},
		{1, 3, ""},
	}

	grid, err := NewGrid(5, 5, 64)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := grid.VerifyInBounds(tt.row, tt.col)
			if got != nil && got.Error() != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func TestGridLayout(t *testing.T) {
	grid, err := NewGrid(4, 4, 12)
	if err != nil {
		t.Fatal(err)
	}

	grid.DrawGrid()
	grid.DrawCoords()
	grid.DrawBorder()

	var data bytes.Buffer
	if err := grid.EncodePNG(&data); err != nil {
		t.Fatal(err)
	}

	str := base64.StdEncoding.EncodeToString(data.Bytes())
	//t.Logf(str)
	assert.True(t, strings.HasPrefix(str, "iVBORw0KGgoAAAANSUhEUgAAAGAAAABgCAIAAABt+uBvAAAErElEQVR4Aeyb3U7qShiG"))
}
