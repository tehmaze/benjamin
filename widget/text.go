package widget

import (
	"image"
	"image/color"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/deck"
)

const (
	fontDPI = 72
)

// TextWidget defaults.
var (
	DefaultTextFont       = benjamin.RobotoBold
	DefaultTextFontSize   = 16.0
	DefaultTextColor      = color.White
	DefaultTextBackground = image.Transparent
)

// TextWidget is a center/middle aligned text.
type TextWidget struct {
	Base
	Text       string
	Font       *truetype.Font
	FontSize   float64
	Color      color.Color
	Background image.Image
	canvas     *image.RGBA
	tile       *image.RGBA
}

func Text(text string) *TextWidget {
	return &TextWidget{
		Base:       makeBase(),
		Text:       text,
		Font:       DefaultTextFont,
		FontSize:   DefaultTextFontSize,
		Color:      DefaultTextColor,
		Background: DefaultTextBackground,
	}
}

func (w *TextWidget) ImageFor(k deck.Key) image.Image {
	pos := k.Position()
	if !pos.In(w.Rect) {
		return nil
	}

	// Calculate our tile size.
	dim := k.Size()
	max := w.Rect.Size()
	if max.X <= 0 {
		dim.X = 1
	}
	if max.Y <= 0 {
		dim.Y = 1
	}
	max.X *= dim.X
	max.Y *= dim.Y

	r := image.Rectangle{Max: max}
	if w.canvas == nil || !w.canvas.Bounds().Eq(r) {
		w.canvas = image.NewRGBA(r)
	}
	if w.tile == nil || !w.tile.Bounds().Eq(image.Rectangle{Max: dim}) {
		w.tile = image.NewRGBA(image.Rectangle{Max: dim})
	}

	if w.Background != nil {
		draw.Draw(w.canvas, r, w.Background, image.Point{}, draw.Src)
	} else {
		draw.Draw(w.canvas, r, image.Transparent, image.Point{}, draw.Src)
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
			Dst:  w.canvas,
			Src:  fontColor,
			Face: fontFace,
		}
		fontInfo   = fontFace.Metrics()
		textHeight = fontInfo.Ascent + fontInfo.Descent
		textY      = (fixed.I(r.Dy()) - textHeight.Mul(fixed.I(strings.Count(w.Text, "\n"))))
	)
	if textY < fixed.I(0) {
		textY = fixed.I(0)
	}

	for _, line := range strings.Split(w.Text, "\n") {
		textX := (fixed.I(r.Dy()) - fontDraw.MeasureString(line)) / 2
		if textX < fixed.I(0) {
			textX = fixed.I(0)
		}
		fontDraw.Dot = fixed.Point26_6{X: textX, Y: textY}
		fontDraw.DrawString(line)
		textY += textHeight
	}

	draw.Copy(w.tile, image.Point{}, w.canvas, image.Rectangle{
		Min: image.Pt((w.Rect.Min.X-pos.X+0)*dim.X, (w.Rect.Min.Y-pos.Y+0)*dim.Y),
		Max: image.Pt((w.Rect.Min.X-pos.X+1)*dim.X, (w.Rect.Min.Y-pos.Y+1)*dim.Y),
	}, draw.Src, nil)

	return w.tile
}

var _ benjamin.Widget = (*TextWidget)(nil)
