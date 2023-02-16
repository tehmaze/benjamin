// Package dummy contains a dummy device for testing/mocking.
package dummy

import (
	"image"

	"github.com/tehmaze/benjamin/device"
)

// Defaults
var (
	Dim          = image.Pt(5, 3)
	KeySize      = image.Pt(72, 72)
	Manufacturer = "maze.io"
	Product      = "Benjamin"
	SerialNumber = "42"
)

type Dummy struct {
	ErrClose   error
	ErrOpen    error
	ErrReset   error
	ChanEvents chan device.Event
}

func (d Dummy) Close() error                { return d.ErrClose }
func (d Dummy) Open() error                 { return d.ErrOpen }
func (d Dummy) Reset() error                { return d.ErrReset }
func (d Dummy) Dim() image.Point            { return Dim }
func (d Dummy) Key(i int) device.Key        { return DummyKey{device: d, index: i} }
func (d Dummy) Keys() int                   { return Dim.X * Dim.Y }
func (d Dummy) KeySize() image.Point        { return KeySize }
func (d Dummy) Manufacturer() string        { return Manufacturer }
func (d Dummy) Product() string             { return Product }
func (d Dummy) SerialNumber() string        { return SerialNumber }
func (d Dummy) Events() <-chan device.Event { return d.ChanEvents }

type DummyKey struct {
	ErrUpdate error
	device    device.Device
	index     int
}

func (d DummyKey) Device() device.Device      { return d.device }
func (d DummyKey) Position() image.Point      { return positionOf(d.device.Dim(), d.index) }
func (d DummyKey) Size() image.Point          { return KeySize }
func (d DummyKey) Update(_ image.Image) error { return d.ErrUpdate }

func init() {
	device.Register(func() bool { return true }, func() device.Device {
		return Dummy{}
	})
}

func positionOf(dim image.Point, i int) image.Point {
	x, y := i%dim.X, i/dim.X
	return image.Pt(x, y)
}
