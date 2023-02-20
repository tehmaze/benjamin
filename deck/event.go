package deck

import (
	"fmt"
	"image"
	"time"
)

type EventType int

const (
	EventTypeKeyPressed EventType = iota
	EventTypeKeyReleased
	EventTypeEncoderPressed
	EventTypeEncoderReleased
	EventTypeEncoderChanged
	EventTypeTouch
	EventTypeTouchEnd
	EventTypeSwipe
	EventTypeError
)

// Event on a deck.
type Event struct {
	// Deck the event was on.
	Deck Deck

	// Type of event.
	Type EventType

	// Data for the event.
	Data EventData
}

type EventData interface{}

func (ev Event) String() string {
	switch ev.Type {
	case EventTypeKeyPressed:
		d := ev.Data.(KeyPress)
		return fmt.Sprintf("key press at %s", d.Key.Position())
	case EventTypeKeyReleased:
		d := ev.Data.(KeyRelease)
		return fmt.Sprintf("key release at %s after %s", d.Position(), d.Duration)
	case EventTypeEncoderPressed:
		d := ev.Data.(EncoderPress)
		return fmt.Sprintf("encoder %d press", d.Encoder.Index())
	case EventTypeEncoderReleased:
		d := ev.Data.(EncoderRelease)
		return fmt.Sprintf("encoder %d release", d.Encoder.Index())
	case EventTypeEncoderChanged:
		d := ev.Data.(EncoderChange)
		return fmt.Sprintf("encoder %d change %d/%d", d.Encoder.Index(), d.Value, d.Scale)
	case EventTypeTouch:
		d := ev.Data.(Touch)
		return fmt.Sprintf("touch %d at %s", d.Index(), d.Point)
	case EventTypeTouchEnd:
		d := ev.Data.(TouchEnd)
		return fmt.Sprintf("touch end %d at %s", d.Index(), d.Point)
	case EventTypeSwipe:
		d := ev.Data.(Swipe)
		return fmt.Sprintf("swipe from %s to %s", d.From, d.To)
	case EventTypeError:
		d := ev.Data.(Error)
		return d.Error()
	default:
		return "invalid"
	}
}

// KeyPress event data.
type KeyPress struct {
	// Key that was pressed.
	Key
}

func KeyPressEvent(key Key) Event {
	return Event{
		Type: EventTypeKeyPressed,
		Data: KeyPress{key},
	}
}

// KeyRelease event data.
type KeyRelease struct {
	// Key that was released.
	Key

	// Duration of the event, or how long the key was pressed for at release.
	Duration time.Duration
}

func KeyReleaseEvent(key Key, after time.Duration) Event {
	return Event{
		Type: EventTypeKeyReleased,
		Data: KeyRelease{
			Key:      key,
			Duration: after,
		},
	}
}

type EncoderPress struct {
	Encoder
}

func EncoderPressEvent(enc Encoder) Event {
	return Event{
		Type: EventTypeEncoderPressed,
		Data: EncoderPress{
			Encoder: enc,
		},
	}
}

type EncoderRelease struct {
	Encoder
	Duration time.Duration
}

func EncoderReleaseEvent(enc Encoder, after time.Duration) Event {
	return Event{
		Type: EventTypeEncoderReleased,
		Data: EncoderRelease{
			Encoder:  enc,
			Duration: after,
		},
	}
}

// EncoderChange event data.
type EncoderChange struct {
	Encoder

	// Value of the encoder.
	Value int

	// Scale of the encoder is the maximum value.
	Scale int
}

func EncoderChangeEvent(enc Encoder, value, scale int) Event {
	return Event{
		Type: EventTypeEncoderChanged,
		Data: EncoderChange{
			Encoder: enc,
			Value:   value,
			Scale:   scale,
		},
	}
}

// Touch event data.
type Touch struct {
	Display
	Point image.Point
}

func TouchEvent(d Display, at image.Point) Event {
	return Event{
		Type: EventTypeTouch,
		Data: Touch{
			Display: d,
			Point:   at,
		},
	}
}

// TouchEnd event data.
type TouchEnd struct {
	Display
	Point image.Point
}

func TouchEndEvent(d Display, at image.Point) Event {
	return Event{
		Type: EventTypeTouchEnd,
		Data: TouchEnd{
			Display: d,
			Point:   at,
		},
	}
}

// Swipe event data.
type Swipe struct {
	Display
	From, To image.Point
}

func SwipeEvent(d Display, from, to image.Point) Event {
	return Event{
		Type: EventTypeSwipe,
		Data: Swipe{
			Display: d,
			From:    from,
			To:      to,
		},
	}
}

// Error event data.
type Error struct {
	// Err is the error that was thrown.
	Err error
}

func ErrorEvent(err error) Event {
	return Event{
		Type: EventTypeError,
		Data: Error{Err: err},
	}
}

func (err Error) Error() string {
	return err.Err.Error()
}
