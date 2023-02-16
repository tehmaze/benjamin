package device

import (
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
