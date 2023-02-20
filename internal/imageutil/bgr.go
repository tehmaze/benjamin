package imageutil

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/tehmaze/benjamin/internal/colorutil"
)

type BGR struct {
	// Pix holds the image's pixels, in B, G, R order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*3].
	Pix []byte
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func NewBGR(r image.Rectangle) *BGR {
	return &BGR{
		Pix:    make([]uint8, 3*r.Dx()*r.Dy()),
		Stride: 3 * r.Dx(),
		Rect:   r,
	}
}

func ToBGR(i image.Image) *BGR {
	if i, ok := i.(*BGR); ok {
		return i
	}
	r := i.Bounds()
	o := NewBGR(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			o.Set(x, y, i.At(x, y))
		}
	}
	return o
}

func (p *BGR) Bounds() image.Rectangle {
	return p.Rect
}

func (p *BGR) ColorModel() color.Model {
	return colorutil.BGRModel
}

func (p *BGR) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return colorutil.BGR{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3]
	return colorutil.BGR{B: s[0], G: s[1], R: s[2]}
}

func (p *BGR) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	b := colorutil.BGRModel.Convert(c).(colorutil.BGR)
	s := p.Pix[i : i+3 : i+3]
	s[0] = b.B
	s[1] = b.G
	s[2] = b.R
}

func (p *BGR) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

var (
	_ image.Image = (*BGR)(nil)
	_ draw.Image  = (*BGR)(nil)
)
