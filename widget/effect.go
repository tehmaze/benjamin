package widget

import (
	"image"
	"image/color"
	"time"

	"github.com/disintegration/imaging"
)

// Effect on a Widget.
type Effect interface {
	UpdateChecker

	Apply(*image.NRGBA, time.Time) *image.NRGBA
}

// Effects are zero or more Effects.
type Effects []Effect

func (fxs Effects) IsUpdated(t time.Time) bool {
	for _, fx := range fxs {
		if fx.IsUpdated(t) {
			return true
		}
	}
	return false
}

func (fxs Effects) Apply(i *image.NRGBA, t time.Time) *image.NRGBA {
	for _, fx := range fxs {
		i = fx.Apply(i, t)
	}
	return i
}

type Rotate struct {
	// Speed in degrees per second
	Speed      float64
	Background color.Color

	first time.Time
	last  time.Time
}

func (fx *Rotate) IsUpdated(t time.Time) bool {
	if fx.first.Equal(never) {
		return false
	}

	return t.Sub(fx.last).Seconds() >= fx.Speed
}

func (fx *Rotate) Apply(i *image.NRGBA, t time.Time) *image.NRGBA {
	if fx.first.Equal(never) {
		fx.first = t
		fx.last = t
		return i
	}

	fx.last = t
	delta := t.Sub(fx.first).Seconds()
	for ; delta > 3600; delta -= 3600 {
		fx.first = fx.first.Add(-time.Hour)
	}
	angle := fx.Speed * delta
	return imaging.Rotate(i, angle, fx.Background)
}
