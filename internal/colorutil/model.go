package colorutil

import "image/color"

var (
	BGRModel color.Model = color.ModelFunc(bgrModel)
)

type BGR struct {
	B, G, R uint8
}

func (c BGR) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = 0xffff
	return
}

func bgrModel(c color.Color) color.Color {
	if _, ok := c.(BGR); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return BGR{uint8(b >> 8), uint8(g >> 8), uint8(r >> 8)}
}
