package infinitton

import (
	"image"

	"github.com/tehmaze/benjamin/deck"
	"github.com/tehmaze/benjamin/internal/hid"
)

const (
	vendorID             = 0xffff
	iDisplayProductID    = 0x1f40
	iDisplayProductIDAlt = 0x1f41
)

func driver(info hid.DeviceInfo) deck.Deck {
	return NewIDisplay(info)
}

type iDisplay struct {
	info hid.DeviceInfo
	dev  *hid.Device
	key  [15]*key
}

func NewIDisplay(info hid.DeviceInfo) *iDisplay {
	d := &iDisplay{
		info: info,
	}
	for i := range d.key {
		d.key[i] = newKey(d, i)
	}
	return d
}

func (d *iDisplay) Open() (err error) {
	d.dev, err = d.info.Open()
	return
}

func (d *iDisplay) Close() error {
	return d.dev.Close()
}

func (d *iDisplay) Reset() error {
	return nil // TODO(maze): not implemented
}

func (d *iDisplay) Manufacturer() string     { return "Infinitton" }
func (d *iDisplay) Product() string          { return "iDisplay" }
func (d *iDisplay) SerialNumber() string     { return d.info.Serial }
func (d *iDisplay) Button(int) deck.Button   { return nil }
func (d *iDisplay) Buttons() int             { return 0 }
func (d *iDisplay) Display(int) deck.Display { return nil }
func (d *iDisplay) Displays() int            { return 0 }
func (d *iDisplay) DisplaySize() image.Point { return image.Point{} }
func (d *iDisplay) Encoder(int) deck.Encoder { return nil }
func (d *iDisplay) Encoders() int            { return 0 }
func (d *iDisplay) Dim() image.Point         { return image.Pt(3, 5) }

func (d *iDisplay) Key(p image.Point) deck.Key {
	if p.X < 0 || p.X >= 3 || p.Y < 0 || p.Y >= 5 {
		return nil
	}

	i := p.Y*3 + p.X
	return d.key[i]
}

func (iDisplay) Keys() int            { return 15 }
func (iDisplay) KeySize() image.Point { return image.Pt(72, 72) }
func (iDisplay) Margin() image.Point  { return image.Point{} }

func (d *iDisplay) Events() <-chan deck.Event {
	c := make(chan deck.Event)

	go func(c chan<- deck.Event) {
		defer close(c)
		for {

		}
	}(c)

	return c
}

func (d *iDisplay) SetBrightness(v float64) error {
	if v < 0.0 {
		v = 0.0
	} else if v > 1.0 {
		v = 1.0
	}

	b := []byte{0x00, 0x11, uint8(v * 100)}
	_, err := d.dev.SendFeatureReport(b)
	return err
}

func init() {
	deck.RegisterUSB(vendorID, iDisplayProductID, driver)
	deck.RegisterUSB(vendorID, iDisplayProductIDAlt, driver)
}
