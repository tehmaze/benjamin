package effect

import (
	"time"

	"github.com/tehmaze/benjamin/effect/ease"
)

const second = float64(time.Second)

// Tweening sequence.
type Tweening struct {
	Overflow time.Duration
	t, d     float64 // time and duration
	begin    float64
	final    float64
	change   float64
	ease     ease.Func
	reverse  bool
}

// Tween function, transition from begin to final over the specified duration.
func Tween(begin, final float64, duration time.Duration, ease ease.Func) *Tweening {
	return &Tweening{
		begin:  begin,
		final:  final,
		change: final - begin,
		d:      float64(duration) / second,
		ease:   ease,
	}
}

func (tween *Tweening) Set(now time.Duration) (current float64, done bool) {
	t := float64(now) / second
	switch {
	case t <= 0:
		tween.Overflow = time.Duration(t * second)
		tween.t = 0
		current = tween.begin
	case t >= tween.d:
		tween.Overflow = time.Duration((t - tween.d) * second)
		tween.t = tween.d
		current = tween.final
	default:
		tween.Overflow = 0
		tween.t = t
		current = tween.ease(tween.t, tween.begin, tween.change, tween.d)
	}

	if tween.reverse {
		return current, tween.t <= 0
	}
	return current, tween.t >= tween.d
}

func (tween *Tweening) Reset() {
	if tween.reverse {
		tween.Set(time.Duration(tween.d * second))
	} else {
		tween.Set(0)
	}
}
