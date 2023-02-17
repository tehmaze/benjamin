package device

import (
	"fmt"
	"image"
	"time"
)

type EventType int

const (
	KeyPressed EventType = iota
	KeyReleased
	Error
)

// Event on a Device.
type Event struct {
	// Device the event was on.
	Device Device

	// Type of event.
	Type EventType

	// Pos is the position of the key.
	Pos image.Point

	// Duration of the event, or how long the key was pressed for at release.
	Duration time.Duration

	// Err on the device
	Err error
}

func (ev Event) String() string {
	switch ev.Type {
	case KeyPressed:
		return fmt.Sprintf("key press at %s", ev.Pos)
	case KeyReleased:
		return fmt.Sprintf("key release at %s after %s", ev.Pos, ev.Duration)
	case Error:
		return ev.Err.Error()
	default:
		return "invalid"
	}
}
