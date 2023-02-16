package benjamin

import (
	"image"
	"time"
)

type KeyEventType int

const (
	KeyPressed KeyEventType = iota
	KeyReleased
)

type KeyEvent struct {
	// Pos is the position of the key.
	Pos image.Point

	// Type of event.
	Type KeyEventType

	// Duration of the event, or how long the key was pressed for at release.
	Duration time.Duration
}
