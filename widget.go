package benjamin

import (
	"image"

	"github.com/tehmaze/benjamin/deck"
)

const fontDPI = 72

type Widget interface {
	// Bounds is the bounding box of the Widget in number of keys.
	Bounds() image.Rectangle

	// Move updates the widget position.
	Move(image.Point)

	// Handle an Event on the Widget.
	Handle(deck.Event)

	// ImageFor returns the Widget image, can return nil if it doesn't have an image.
	ImageFor(deck.Key) image.Image

	// Render the Widget on the surface.
	//Render(device.Surface) error

	// UpdateRequired indicates the Widget has updated and the Key should be rendered.
	UpdateRequired() bool

	// IsVisible indicates the Widget is visible or hidden.
	IsVisible() bool
}

type Grid struct {
	Widgets []Widget
	Stride  int
}

/*
func (w *Tile) Render(s device.Surface) error {
	if k := device.KeyAt(s, w.Position.X, w.Position.Y); k != nil {
		w.IsClean = true
		return k.Update(w.Texture)
	}
	return nil
}
*/
