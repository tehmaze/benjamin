package widget

import (
	"image"
	"time"

	"github.com/tehmaze/benjamin"
)

var never time.Time

type Widget interface {
	UpdateChecker

	// Drawable returns the underlying Peripheral.
	Drawable() benjamin.Drawable

	// Frame is the Widget image at time.
	Frame(time.Time) *image.NRGBA
}

type UpdateChecker interface {
	IsUpdated(t time.Time) bool
}

type DrawablePeripheral interface {
	benjamin.Drawable
	benjamin.Peripheral
}

// Base widget
type Base struct {
	ConnectedTo DrawablePeripheral
	Effects     Effects
	canvas      *image.NRGBA
}

func MakeBase(peripheral DrawablePeripheral, effects ...Effect) Base {
	return Base{
		ConnectedTo: peripheral,
		Effects:     effects,
	}
}

func (w *Base) Peripheral() benjamin.Peripheral {
	return w.ConnectedTo
}

func (w *Base) Drawable() benjamin.Drawable {
	return w.ConnectedTo
}

func (w *Base) IsUpdated(t time.Time) bool {
	return w.Effects.IsUpdated(t)
}

func (w *Base) Frame(i *image.NRGBA, t time.Time) *image.NRGBA {
	if len(w.Effects) == 0 {
		return i
	}

	if w.canvas == nil || !w.canvas.Rect.Eq(i.Rect) {
		w.canvas = image.NewNRGBA(i.Rect)
	}
	copy(w.canvas.Pix, i.Pix)

	return w.Effects.Apply(w.canvas, t)
}
