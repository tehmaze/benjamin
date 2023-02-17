package widget

import (
	"image"
	"time"

	"github.com/tehmaze/benjamin/device"
)

// Base widget, does nothing.
type Base struct {
	Rect         image.Rectangle
	OnKeyPress   func(device.Device)
	OnKeyRelease func(device.Device, time.Duration)
	IsClean      bool
	IsHidden     bool
}

func makeBase() Base {
	return Base{
		Rect: image.Rect(0, 0, 1, 1),
	}
}

func (w *Base) Bounds() image.Rectangle {
	return w.Rect
}

func (w *Base) Move(p image.Point) {
	w.Rect = image.Rect(
		p.X,
		p.Y,
		p.X+w.Rect.Dx(),
		p.Y+w.Rect.Dy(),
	)
}

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

func (w *Base) ImageFor(_ device.Key) image.Image {
	w.IsClean = true
	return nil
}

func (w *Base) UpdateRequired() bool { return !w.IsClean }
func (w *Base) IsVisible() bool      { return !w.IsHidden }
func (w *Base) Dirty()               { w.IsClean = false }
func (w *Base) Hide()                { w.IsHidden = true }
func (w *Base) Show()                { w.IsHidden = false }
