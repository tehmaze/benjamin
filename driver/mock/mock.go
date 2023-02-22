package mock

import (
	"image"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/driver"
)

var (
	ErrOpen          error
	ErrClose         error
	ErrReset         error
	ErrClear         error
	ErrSetBrightness error
	Displays         int
	Encoders         int
	Buttons             int
	ButtonLayout        image.Point
	ButtonSize          image.Point
)

// Mock interface
type Mock struct {
}

func New() benjamin.Device {
	return new(Mock)
}

func (Mock) Manufacturer() string           { return "maze.io" }
func (Mock) Product() string                { return "mock" }
func (Mock) Serial() string                 { return "2342" }
func (Mock) Open() error                    { return ErrOpen }
func (Mock) Close() error                   { return ErrClose }
func (Mock) Reset() error                   { return ErrReset }
func (Mock) Clear() error                   { return ErrClear }
func (Mock) Display(int) benjamin.Display   { return nil }
func (Mock) Displays() int                  { return Displays }
func (Mock) Encoder(int) benjamin.Encoder   { return nil }
func (Mock) Encoders() int                  { return Encoders }
func (Mock) Button(int) benjamin.Button           { return nil }
func (Mock) ButtonAt(image.Point) benjamin.Button { return nil }
func (Mock) ButtonLayout() image.Point         { return ButtonLayout }
func (Mock) ButtonSize() image.Point           { return ButtonSize }
func (Mock) Buttons() int                      { return Buttons }
func (Mock) SetBrightness(float64) error    { return ErrSetBrightness }

func (Mock) Events() <-chan benjamin.Event {
	c := make(chan benjamin.Event)
	close(c)
	return c
}

func init() {
	driver.Register(func() bool { return true }, New)
}
