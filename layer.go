package benjamin

import (
	"image"
	"image/draw"
	"log"

	"github.com/tehmaze/benjamin/device"
)

type Layer struct {
	Device  device.Device
	Widgets []Widget
	dim     image.Point
}

func NewLayer(device device.Device) *Layer {
	dim := device.Dim()
	return &Layer{
		Device:  device,
		Widgets: make([]Widget, dim.X*dim.Y),
		dim:     dim,
	}
}

func (l *Layer) AddWidget(w Widget, at image.Point) (ok bool) {
	if at.X < 0 || at.X >= l.dim.X || at.Y < 0 || at.Y >= l.dim.Y {
		// Not visible on this layer.
		return false
	}

	dim := w.Size()
	for y := at.Y; y < at.Y+dim.Y && y < l.dim.Y; y++ {
		for x := at.X; x < at.X+dim.X && x < l.dim.X; x++ {
			i := y*l.dim.X + x
			l.Widgets[i] = w
			ok = true
		}
	}
	return
}

func (l *Layer) RemoveWidget(w Widget) (has bool) {
	for i, o := range l.Widgets {
		if o == w {
			has = true
			l.Widgets[i] = nil
		}
	}
	return
}

func (l *Layer) Widget(x, y int) Widget {
	if x < 0 || x >= l.dim.X || y < 0 || y >= l.dim.Y {
		return nil
	}

	i := y*l.dim.X + x
	return l.Widgets[i]
}

// Layers are zero or more stacked layers.
type Layers []*Layer

func (ls Layers) Dirty(x, y int) bool {
	if len(ls) == 0 {
		return false
	}

	i := y*ls[0].dim.X + x
	return ls.IsDirtyIndex(i)
}

func (ls Layers) IsDirtyIndex(i int) bool {
	for _, l := range ls {
		if w := l.Widgets[i]; w != nil && w.IsDirty() {
			return true
		}
	}
	return false
}

func (ls Layers) Render(x, y int) image.Image {
	if len(ls) == 0 {
		return nil
	}

	i := y*ls[0].dim.X + x
	return ls.RenderIndex(i)
}

func (ls Layers) RenderIndex(i int) image.Image {
	var o *image.RGBA
	for n, l := range ls {
		w := l.Widgets[i]
		if w == nil {
			continue // Drop to next layer
		}

		log.Println("benjamin: render image layer", n, "index", i)
		if r := w.Render(l.Device); r != nil {
			if o == nil {
				// First image we encountered, create same-size image for compositing
				o = image.NewRGBA(r.Bounds())
			}

			// Blend
			draw.Draw(o, o.Bounds(), r, image.Point{}, draw.Over)
		}
	}
	return o
}
