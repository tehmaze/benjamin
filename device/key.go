package device

import "image"

// Key on a device.
type Key interface {
	// Position opn the Device.
	Position() image.Point

	// Size of the key in pixels.
	Size() image.Point

	// Update the key graphics.
	Update(image.Image) error

	// Surface the key is connected to.
	Surface() Surface
}

// Resetter can be used to check if a Key allows resetting.
type Resetter interface {
	Reset() error
}
