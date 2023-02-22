package widget

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/tehmaze/benjamin/internal/fontutil"
)

// Image widget, can be used for keys or displays.
type Image struct {
	Base
	Image *image.NRGBA
}

func (w *Image) IsUpdated(t time.Time) bool {
	return w.canvas == nil || !w.canvas.Rect.Eq(w.Image.Rect) || w.Effects.IsUpdated(t)
}

func (w *Image) Frame(t time.Time) *image.NRGBA {
	return w.Base.Frame(w.Image, t)
}

func (w *Image) Set(i image.Image) {
	size := w.Base.ConnectedTo.Size()
	w.Image = imaging.Resize(i, size.X, size.Y, imaging.Lanczos)
	w.canvas = nil
}

type Progress struct {
	Base
	Value float64

	opts *ProgressOptions
	icon *image.NRGBA
	face font.Face
	last float64
}

type ProgressOptions struct {
	Label      string
	Font       *truetype.Font
	FontSize   float64
	Fill       color.Color
	Color      color.Color
	Background image.Image
	Icon       image.Image
}

var DefaultProgressOptions = ProgressOptions{
	Font:       fontutil.RobotoBold,
	FontSize:   24,
	Fill:       color.White,
	Color:      color.NRGBA{R: 0x7f, G: 0x7f, B: 0xff, A: 0xff},
	Background: image.Transparent,
}

func (o *ProgressOptions) Defaults() {
	if o.Font == nil {
		o.Font = DefaultProgressOptions.Font
	}
	if o.FontSize <= 0 {
		o.FontSize = DefaultProgressOptions.FontSize
	}
	if o.Fill == nil {
		o.Fill = DefaultProgressOptions.Fill
	}
	if o.Color == nil {
		o.Color = DefaultProgressOptions.Color
	}
	if o.Background == nil {
		o.Background = DefaultProgressOptions.Background
	}
	if o.Icon == nil {
		o.Icon = DefaultProgressOptions.Icon
	}
}

func NewProgress(p DrawablePeripheral, opts *ProgressOptions, effects ...Effect) *Progress {
	if opts == nil {
		opts = new(ProgressOptions)
	}
	opts.Defaults()

	var (
		w = &Progress{
			Base:  MakeBase(p, effects...),
			opts:  opts,
			Value: math.Inf(+1),
			last:  math.Inf(-1),
			face:  truetype.NewFace(opts.Font, &truetype.Options{Size: opts.FontSize}),
		}
	)
	if opts.Icon != nil {
		size := p.Size()
		switch {
		case size.X > size.Y: // Landscape
			size.Y /= 2
			size.X = size.Y
		case size.X < size.Y: // Portrait
			fallthrough
		default: // Square
			size.X /= 2
			size.Y = size.X
		}
		log.Printf("progress icon resize %s->%s", opts.Icon.Bounds().Max, size)
		w.icon = imaging.Resize(opts.Icon, size.X, size.Y, imaging.Lanczos)
	}

	return w
}

func (w *Progress) IsUpdated(t time.Time) bool {
	return !almostEqual(w.Value, w.last, 1e-5) || w.Base.IsUpdated(t)
}

func (w *Progress) Set(value float64) {
	if value < 0 {
		value = 0
	} else if value > 100 {
		value = 100
	}
	w.Value = value
}

func (w *Progress) SetColor(c color.Color) {
	w.opts.Color = c
}

func (w *Progress) SetFill(c color.Color) {
	w.opts.Fill = c
}

func (w *Progress) SetBackground(i image.Image) {
	w.opts.Background = i
}

func (w *Progress) Frame(t time.Time) *image.NRGBA {
	if w.canvas == nil {
		w.canvas = image.NewNRGBA(image.Rectangle{Max: w.ConnectedTo.Size()})
	}

	draw.Draw(w.canvas, w.canvas.Rect, w.opts.Background, image.Point{}, draw.Src)

	var (
		dx = w.canvas.Rect.Dx()
		dy = w.canvas.Rect.Dy()
		r  = image.Rectangle{
			Min: image.Pt(4, dy-7),
			Max: image.Pt(dx-4, dy-2),
		}
	)
	switch {
	case dx > dy: // Landscape
		if w.icon != nil {
			draw.Copy(w.canvas, image.Pt(2, w.canvas.Rect.Max.Y/2-w.icon.Rect.Dy()/2), w.icon, w.icon.Rect, draw.Over, nil)
		}
	case dx < dy: // Portrait
		fallthrough
	default: // Square
		if w.icon != nil {
			draw.Copy(w.canvas, image.Pt(2, w.canvas.Rect.Max.Y/2-w.icon.Rect.Dy()/2), w.icon, w.icon.Rect, draw.Over, nil)
		}
	}

	v := int((float64(dx) / 100) * w.Value)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if (y == r.Min.Y || y == r.Max.Y-1) && (x == r.Min.X || x == r.Max.X-1) {
				continue
			}
			if x < v {
				w.canvas.Set(x, y, w.opts.Fill)
			} else {
				w.canvas.Set(x, y, w.opts.Color)
			}
		}
	}

	var (
		// h = w.face.Metrics().Height.Ceil()
		h = w.face.Metrics().Ascent.Floor()
		d = &font.Drawer{
			Dst:  w.canvas,
			Src:  image.NewUniform(w.opts.Fill),
			Face: w.face,
		}
		l = fmt.Sprintf("%d%%", int(math.Ceil(w.Value)))
	)

	// Draw label (if any)
	if w.opts.Label != "" {
		d.Src = image.NewUniform(w.opts.Color)
		d.Dot = fixed.P(w.canvas.Rect.Max.X-2, w.canvas.Rect.Max.Y/2-w.face.Metrics().Height.Ceil()/2-2)
		d.Dot.X -= d.MeasureString(w.opts.Label)
		d.DrawString(w.opts.Label)
	}

	// Draw percentage
	d.Src = image.NewUniform(w.opts.Fill)
	d.Dot = fixed.P(w.canvas.Rect.Max.X-2, w.canvas.Rect.Max.Y/2+h)
	d.Dot.X -= d.MeasureString(l)
	d.DrawString(l)

	//log.Printf("progress: %s: %f->%f", w.canvas.Rect, w.last, w.Value)
	w.last = w.Value
	return w.canvas
}

func almostEqual(a, b, E float64) bool {
	return math.Abs(a-b) < E
}

var (
	_ Widget = (*Image)(nil)
)
