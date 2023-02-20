package widget

import (
	"image"
	"image/color"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/deck"
)

// Solid color widget.
type ColorWidget struct {
	Base
	Color color.Color
}

func Color(c color.Color) *ColorWidget {
	return &ColorWidget{
		Base:  makeBase(),
		Color: c,
	}
}

func (w *ColorWidget) ImageFor(k deck.Key) image.Image {
	w.IsClean = true
	return solidColorImage{w.Color}
}

type solidColorImage struct {
	Color color.Color
}

func (i solidColorImage) At(x, y int) color.Color           { return i.Color }
func (i solidColorImage) Convert(_ color.Color) color.Color { return i.Color }
func (i solidColorImage) ColorModel() color.Model           { return i }
func (i solidColorImage) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Point{X: 1, Y: 1}}
}

var _ benjamin.Widget = (*ColorWidget)(nil)
