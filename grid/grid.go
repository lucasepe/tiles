package grid

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// Grid represents the grid structure
type Grid struct {
	cellSize          int
	rows              int
	cols              int
	margin            int
	lineDashes        float64
	lineColor         string
	lineStrokeWidth   float64
	borderDashes      float64
	borderColor       string
	borderStrokeWidth float64
	backgroundColor   string

	watermark string

	canvasWidth  int
	canvasHeight int

	imageWidth  int
	imageHeight int

	font *truetype.Font
	ctx  *gg.Context
}

// NewGrid creates a new grid and sets it up with its configuration
func NewGrid(rows, cols int, cellSize int, opts ...func(*Grid)) (*Grid, error) {
	if rows == 0 {
		return nil, fmt.Errorf("no rows provided")
	}

	if cols == 0 {
		return nil, fmt.Errorf("no columns provided")
	}

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	res := Grid{
		rows: rows, cols: cols,
		cellSize:        cellSize,
		margin:          24,
		lineColor:       "#b8b8a7",
		backgroundColor: "#ffffff",
		borderColor:     "#161615",
		font:            font,
	}

	for _, opt := range opts {
		opt(&res)
	}

	res.canvasWidth = res.cellSize * res.cols
	res.canvasHeight = res.cellSize * res.rows
	res.imageWidth = res.canvasWidth + 2*res.margin
	res.imageHeight = res.canvasHeight + 2*res.margin

	max := res.imageWidth
	if max < res.imageHeight {
		max = res.imageHeight
	}
	res.borderColor = "#ffffff00"
	res.borderStrokeWidth = 0.002 * float64(max)
	res.lineStrokeWidth = 0.001 * float64(max)

	res.ctx = gg.NewContext(res.imageWidth, res.imageHeight)
	res.ctx.Translate(float64(res.margin), float64(res.margin))
	res.ctx.SetHexColor(res.backgroundColor)
	res.ctx.Clear()

	return &res, nil
}

// Context returns the grid drawing context
func (g *Grid) Context() *gg.Context {
	return g.ctx
}

// EncodePNG encodes the final image as PNG
func (g *Grid) EncodePNG(w io.Writer) error {
	// specify compression level
	enc := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
	if err := enc.Encode(w, g.ctx.Image()); err != nil {
		return err
	}

	return nil
}

// SavePNG saves the grid as PNG image.
func (g *Grid) SavePNG(filename string) error {
	if filename == "" {
		currentTime := time.Now()
		ctf := currentTime.Format("200601021504")
		filename = fmt.Sprintf("GRID%s.png", ctf)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return g.EncodePNG(f)
}

// DrawBorder draws a border around the grid.
func (g *Grid) DrawBorder() {
	canvasWidth := float64(g.cellSize * g.cols)
	canvasHeight := float64(g.cellSize * g.rows)

	g.ctx.Push()
	g.ctx.MoveTo(0, 0)
	g.ctx.LineTo(0, canvasHeight)
	g.ctx.LineTo(canvasWidth, canvasHeight)
	g.ctx.LineTo(canvasWidth, 0)
	g.ctx.LineTo(0, 0)

	if g.borderDashes > 0 {
		g.ctx.SetDash(g.borderDashes)
	} else {
		g.ctx.SetDash()
	}

	g.ctx.SetLineWidth(g.borderStrokeWidth)
	g.ctx.SetHexColor(g.borderColor)
	g.ctx.Stroke()
	g.ctx.Pop()
}

// DrawWatermark draws the watermark.
func (g *Grid) DrawWatermark() {
	const lineSpacing = 1.2

	w, h := float64(g.ctx.Width()), float64(g.ctx.Height())

	fontSize := 0.15 * math.Min(w, h)
	face := truetype.NewFace(g.font, &truetype.Options{Size: fontSize})

	g.ctx.Push()
	g.ctx.SetHexColor("#c0c0c088")
	g.ctx.SetFontFace(face)

	s := math.Min(g.ctx.MeasureMultilineString(g.watermark, lineSpacing))
	x := 0.5 * (w - 2*float64(g.margin))
	y := 0.5 * (h - 2*float64(g.margin))
	g.ctx.DrawStringWrapped(g.watermark, x, y, 0.5, 0.5, s, lineSpacing, gg.AlignCenter)
	g.ctx.Pop()
}

// DrawGrid draws the grid.
func (g *Grid) DrawGrid() {
	g.ctx.Push()
	for i := 1; i < g.cols; i++ {
		x := float64(i * g.cellSize)
		g.ctx.MoveTo(x, 0)
		g.ctx.LineTo(x, float64(g.canvasHeight))
	}

	for i := 1; i < g.rows; i++ {
		y := float64(i * g.cellSize)
		g.ctx.MoveTo(0, y)
		g.ctx.LineTo(float64(g.canvasWidth), y)
	}

	if g.lineDashes > 0 {
		g.ctx.SetDash(g.lineDashes)
	} else {
		g.ctx.SetDash()
	}
	if g.lineColor != "" {
		g.ctx.SetHexColor(g.lineColor)
	}

	g.ctx.SetLineWidth(g.lineStrokeWidth)
	g.ctx.Stroke()
	g.ctx.Pop()
}

// DrawImage draws the image at row and col.
// If the image size is greater then the gtid cell size
// it will be shrinked.
func (g *Grid) DrawImage(img image.Image, row, col int) error {
	if err := g.VerifyInBounds(row, col); err != nil {
		return err
	}

	b := img.Bounds()
	size := b.Max.X
	if b.Max.Y > size {
		size = b.Max.Y
	}

	if g.cellSize < size {
		img = imaging.Resize(img, g.cellSize, g.cellSize, imaging.Lanczos)
	}

	center := g.CellCenter(row, col)

	dc := g.Context()
	dc.Push()
	//g.ctx.RotateAbout(gg.Radians(alpha), center.X, center.Y)
	dc.DrawImageAnchored(img, int(center.X), int(center.Y), 0.5, 0.5)
	dc.Pop()

	return nil
}

// DrawCoords draws all cells locations
func (g *Grid) DrawCoords() {
	cs := g.CellSize()
	fontSize := 0.3 * cs

	face := truetype.NewFace(g.font, &truetype.Options{Size: fontSize})

	g.ctx.Push()
	g.ctx.SetFontFace(face)
	g.ctx.SetHexColor("#00000099")
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			txt := fmt.Sprintf("%d,%d", i, j)
			center := g.CellCenter(i, j)
			sw, sh := g.ctx.MeasureString(txt)

			g.ctx.Push()
			g.ctx.SetHexColor("#00000022")
			g.ctx.DrawRoundedRectangle(center.X-0.5*sw-4, center.Y-0.5*sh-4, sw+8, sh+8, 4)
			g.ctx.Fill()
			g.ctx.Pop()

			g.ctx.DrawStringAnchored(txt, center.X, center.Y, 0.5, 0.35)
		}
	}
	g.ctx.Pop()
}

// CellSize returns the cell dimension
func (g *Grid) CellSize() float64 {
	return float64(g.cellSize)
}

// CellCenter retuns the cell coordinates in the grid
func (g *Grid) CellCenter(row, col int) gg.Point {
	size := g.CellSize()

	x := 0.5*size + float64(col)*size
	y := 0.5*size + float64(row)*size

	return gg.Point{X: x, Y: y}
}

// VerifyInBounds verify that the coordinates
// belongs to the grid
func (g *Grid) VerifyInBounds(row, col int) error {
	if row < 0 || row >= g.rows || col < 0 || col >= g.cols {
		return fmt.Errorf("cell (%d, %d) is out of bounds", row, col)
	}
	return nil
}

// Background sets the grid background color
func Background(hex string) func(*Grid) {
	return func(g *Grid) {
		g.backgroundColor = hex
	}
}

// Margin sets the grid margin in pixels.
func Margin(val int) func(*Grid) {
	return func(g *Grid) {
		g.margin = val
	}
}

// Watermark sets the grid watermark.
func Watermark(text string) func(*Grid) {
	return func(g *Grid) {
		g.watermark = text
	}
}
