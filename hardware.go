package benjamin

import "image"

type Device interface {
	Manufacturer() string
	Product() string
	Serial() string

	Open() error
	Close() error
	Reset() error

	// Clear all displays and buttons to black.
	Clear() error

	Events() <-chan Event

	Surface
}

type Surface interface {
	Display(int) Display
	Displays() int

	Encoder(int) Encoder
	Encoders() int

	Key(int) Key
	KeyAt(image.Point) Key
	Keys() int
	KeyLayout() image.Point

	SetBrightness(float64) error
}

type Peripheral interface {
	// Surface the pheripheral is connected to.
	Surface() Surface

	// Index of the pheripheral.
	Index() int
}

type Drawable interface {
	Size() image.Point
	SetImage(image.Image) error
}

type Display interface {
	Peripheral
	Drawable
}

// Encoder is a rotary encoder.
type Encoder interface {
	Peripheral

	// Display linked to the encoder, returns nil if the encoder has no display.
	Display() Display
}

type Key interface {
	Peripheral
	Drawable

	// Position on the key matrix.
	Position() image.Point
}
