// Package dummy contains a dummy device for testing/mocking.
package dummy

import (
	"image"

	"github.com/tehmaze/benjamin/deck"
)

// Defaults
var (
	Dim          = image.Pt(5, 3)
	KeySize      = image.Pt(72, 72)
	Margin       = image.Pt(8, 8)
	Manufacturer = "maze.io"
	Product      = "Benjamin"
	SerialNumber = "42"
)

type Dummy struct {
	ErrClose         error
	ErrOpen          error
	ErrReset         error
	ErrSetBrightness error
	ChanEvents       chan deck.Event
}

func (d Dummy) Close() error               { return d.ErrClose }
func (d Dummy) Open() error                { return d.ErrOpen }
func (d Dummy) Reset() error               { return d.ErrReset }
func (d Dummy) Dim() image.Point           { return Dim }
func (d Dummy) Button(int) deck.Button     { return nil }
func (d Dummy) Buttons() int               { return 0 }
func (d Dummy) Display(int) deck.Display   { return nil }
func (d Dummy) Displays() int              { return 0 }
func (d Dummy) DisplaySize() image.Point   { return image.Point{} }
func (d Dummy) Encoder(int) deck.Encoder   { return nil }
func (d Dummy) Encoders() int              { return 0 }
func (d Dummy) Key(p image.Point) deck.Key { return DummyKey{deck: d, point: p} }
func (d Dummy) Keys() int                  { return Dim.X * Dim.Y }
func (d Dummy) KeySize() image.Point       { return KeySize }
func (d Dummy) Margin() image.Point        { return Margin }
func (d Dummy) Manufacturer() string       { return Manufacturer }
func (d Dummy) Product() string            { return Product }
func (d Dummy) SerialNumber() string       { return SerialNumber }
func (d Dummy) Events() <-chan deck.Event  { return d.ChanEvents }
func (d Dummy) SetBrightness(uint8) error  { return d.ErrSetBrightness }

type DummyKey struct {
	ErrUpdate error
	deck      deck.Deck
	point     image.Point
}

func (d DummyKey) Surface() deck.Surface      { return d.deck }
func (d DummyKey) Device() deck.Deck          { return d.deck }
func (d DummyKey) Position() image.Point      { return d.point }
func (d DummyKey) Size() image.Point          { return KeySize }
func (d DummyKey) Update(_ image.Image) error { return d.ErrUpdate }

func init() {
	deck.Register(func() bool { return true }, func() deck.Deck {
		return Dummy{}
	})
}

func positionOf(dim image.Point, i int) image.Point {
	x, y := i%dim.X, i/dim.X
	return image.Pt(x, y)
}
