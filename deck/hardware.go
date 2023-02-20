package deck

import (
	"image"
	"image/color"
)

// Pheripheral connected to a deck Surface.
type Pheripheral interface {
	// Surface the pheripheral is connected to.
	Surface() Surface
}

// Drawable allows drawing an image.
type Drawable interface {
	// Size of the display in pixels.
	Size() image.Point

	// Update the display graphics.
	Update(image.Image) error
}

// Display is a (touch) display.
type Display interface {
	Pheripheral
	Drawable

	// Index of the display.
	Index() int
}

// Button that can be pressed, and maybe has a (color) LED.
type Button interface {
	Pheripheral

	// Position opn the Surface.
	Position() image.Point

	// Update the button color.
	UpdateColor(color.Color) error
}

// Key with a screen on a Surface.
type Key interface {
	Pheripheral
	Drawable

	// Position on the Surface.
	Position() image.Point
}

// Encoder is a rotary encoder.
type Encoder interface {
	Pheripheral

	// Index of the encoder.
	Index() int
}
