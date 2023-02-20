package window

import (
	"image"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/tehmaze/benjamin/deck"
)

const (
	dimX    = 5
	dimY    = 3
	keySize = 96
)

// Title of the window.
var Title = "Benjamin"

var (
	a      = app.New()
	window fyne.Window
	once   sync.Once
)

type Window struct {
}

func New() *Window {
	return new(Window)
}

func (d *Window) Open() error {
	if a != nil {
		return nil
	}

	once.Do(func() {
		window = a.NewWindow(Title)
		window.Resize(fyne.NewSize(dimX*keySize, dimY*keySize))
		go window.ShowAndRun()
	})

	return nil
}

func (d *Window) Close() error {
	a.Quit()
	return nil
}

func (d *Window) Reset() error {
	return nil
}

func (Window) Dim() image.Point          { return image.Pt(dimX, dimY) }
func (Window) Key(image.Point) deck.Key  { return nil }
func (Window) Keys() int                 { return dimX * dimY }
func (Window) KeySize() image.Point      { return image.Pt(keySize, keySize) }
func (Window) Margin() image.Point       { return image.Pt(0, 0) }
func (Window) SetBrightness(uint8) error { return nil }

func (d *Window) Events() <-chan deck.Event {
	e := make(chan deck.Event, 8)
	return e
}

func init() {
	deck.Register(func() bool {
		return true
	}, func() deck.Deck {
		return New()
	})
}

var _ deck.Deck = (*Window)(nil)
