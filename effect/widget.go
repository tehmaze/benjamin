package effect

import (
	"image"
	"time"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/device"
)

type blink struct {
	benjamin.Widget
	ticker    *time.Ticker
	isVisible bool
	isClean   bool
}

// Blink creates a blinking widget.
func Blink(w benjamin.Widget, interval time.Duration) benjamin.Widget {
	x := &blink{
		Widget: w,
		ticker: time.NewTicker(interval / 2),
	}

	go func(t <-chan time.Time, w *blink) {
		w.isVisible = !w.isVisible
	}(x.ticker.C, x)

	return x
}

func (w *blink) Close() error {
	w.ticker.Stop()
	return nil
}

func (w *blink) IsDirty() bool {
	return !w.isClean
}

func (w *blink) Render(d device.Device) image.Image {
	if w.isVisible {
		return w.Widget.Render(d)
	}
	return image.Transparent
}
