package benjamin

import (
	"image"
	"image/color"
	"image/draw"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/tehmaze/benjamin/device"
)

const fontDPI = 72

type Widget interface {
	// Size of the widget in number of keys in X and Y direction.
	Size() image.Point

	// Handle an Event on the Widget.
	Handle(device.Event)

	// Render the Widget.
	Render(device.Device) image.Image

	// IsDirty indicates the Widget has updated and the Key should be rendered.
	IsDirty() bool
}

type Grid struct {
	Widgets []Widget
	Stride  int
}

// Base widget, does nothing.
type Base struct {
	OnKeyPress   func(device.Device)
	OnKeyRelease func(device.Device, time.Duration)
	Clean        bool
}

func (Base) Size() image.Point { return image.Pt(1, 1) }

func (w *Base) Handle(event device.Event) {
	switch event.Type {
	case device.KeyPressed:
		if w.OnKeyPress != nil {
			w.OnKeyPress(event.Device)
		}
	case device.KeyReleased:
		if w.OnKeyRelease != nil {
			w.OnKeyRelease(event.Device, event.Duration)
		}
	}
}

func (Base) Render(_ device.Device) image.Image {
	return image.Transparent
}

func (w *Base) Update(_ image.Image) error {
	w.Clean = true
	return nil
}

func (w *Base) IsDirty() bool { return !w.Clean }

// Icon is an image.
type Icon struct {
	Base
	image.Image
}

func (Icon) Size() image.Point { return image.Pt(1, 1) }

func (w Icon) Render(_ device.Device) image.Image {
	return w.Image
}

func (w *Icon) Update(i image.Image) error {
	w.Image = i
	w.Clean = false
	return nil
}

// Solid color widget.
type Solid struct {
	Base
	Color color.Color
}

func (w *Solid) Render(_ device.Device) image.Image {
	return image.NewUniform(w.Color)
}

// Text is a center/middle aligned text.
type Text struct {
	Base
	Text       string
	Font       *truetype.Font
	FontSize   float64
	Color      color.Color
	Background image.Image
	tile       *image.RGBA
}

func (w *Text) Render(d device.Device) image.Image {
	r := image.Rectangle{Max: d.KeySize()}
	if w.tile == nil || !w.tile.Bounds().Eq(r) {
		w.tile = image.NewRGBA(r)
	}

	if w.Background != nil {
		draw.Draw(w.tile, r, w.Background, image.Point{}, draw.Src)
	}

	var (
		fontSize  = w.FontSize
		fontColor = image.NewUniform(w.Color)
	)
	if fontSize == 0 {
		fontSize = 16
	}
	var (
		fontFace = truetype.NewFace(w.Font, &truetype.Options{
			Size: fontSize,
			DPI:  fontDPI,
		})
		fontDraw = &font.Drawer{
			Dst:  w.tile,
			Src:  fontColor,
			Face: fontFace,
		}
		fontInfo   = fontFace.Metrics()
		textHeight = fontInfo.Ascent + fontInfo.Descent
		y          = (fixed.I(r.Dy()) - textHeight.Mul(fixed.I(strings.Count(w.Text, "\n"))))
	)
	if y < fixed.I(0) {
		y = fixed.I(0)
	}

	for _, line := range strings.Split(w.Text, "\n") {
		x := (fixed.I(r.Dy()) - fontDraw.MeasureString(line)) / 2
		if x < fixed.I(0) {
			x = fixed.I(0)
		}
		fontDraw.Dot = fixed.Point26_6{X: x, Y: y}
		fontDraw.DrawString(line)
		y += textHeight
	}

	return w.tile
}

// Compile-time interface compliance checks.
var (
	_ Widget = (*Base)(nil)
	_ Widget = (*Icon)(nil)
	_ Widget = (*Solid)(nil)
	_ Widget = (*Text)(nil)
)
