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

type USBDevice interface {
	USBID() (vendorID, productID uint16)
	Path() string
}

type Surface interface {
	Display(int) Display
	Displays() int
	DisplayArea() Screen

	Encoder(int) Encoder
	Encoders() int

	Button(int) Button
	ButtonAt(image.Point) Button
	Buttons() int
	ButtonLayout() image.Point
	ButtonArea() Screen

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

type Button interface {
	Peripheral
	Drawable

	// Position on the key matrix.
	Position() image.Point
}

type Screen interface {
	Peripheral
	Drawable
}
