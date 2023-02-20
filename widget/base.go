package widget

import (
	"image"
	"time"

	"github.com/tehmaze/benjamin/deck"
)

// Base widget, does nothing.
type Base struct {
	Rect         image.Rectangle
	OnKeyPress   func(deck.Key)
	OnKeyRelease func(deck.Key, time.Duration)
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

func (w *Base) Handle(event deck.Event) {
	switch event := event.Data.(type) {
	case deck.KeyPress:
		if w.OnKeyPress != nil {
			w.OnKeyPress(event.Key)
		}
	case deck.KeyRelease:
		if w.OnKeyRelease != nil {
			w.OnKeyRelease(event.Key, event.Duration)
		}
	}
}

func (w *Base) ImageFor(_ deck.Key) image.Image {
	w.IsClean = true
	return nil
}

func (w *Base) UpdateRequired() bool { return !w.IsClean }
func (w *Base) IsVisible() bool      { return !w.IsHidden }
func (w *Base) Dirty()               { w.IsClean = false }

func (w *Base) Hide() {
	w.IsHidden = true
	w.IsClean = false
}

func (w *Base) Show() {
	w.IsHidden = false
	w.IsClean = false
}
