package window

import (
	"image"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/tehmaze/benjamin/device"
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

func (Window) Dim() image.Point     { return image.Pt(dimX, dimY) }
func (Window) Key(int) device.Key   { return nil }
func (Window) Keys() int            { return dimX * dimY }
func (Window) KeySize() image.Point { return image.Pt(keySize, keySize) }

func (d *Window) Events() <-chan device.Event {
	e := make(chan device.Event, 8)
	return e
}

func init() {
	device.Register(func() bool {
		return true
	}, func() device.Device {
		return New()
	})
}

var _ device.Device = (*Window)(nil)
