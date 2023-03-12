package infinitton

import (
	"image"

	"github.com/karalabe/hid"
	"golang.org/x/image/draw"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/driver"
	"github.com/tehmaze/benjamin/internal/imageutil"
)

const (
	vendorID             = 0xffff
	iDisplayProductID    = 0x1f40
	iDisplayProductIDAlt = 0x1f41
)

func New(info hid.DeviceInfo) benjamin.Device {
	return NewIDisplay(info)
}

type iDisplay struct {
	info         hid.DeviceInfo
	dev          *hid.Device
	button       [15]*button
	buttonCanvas *imageutil.BGR
	canvas       *imageutil.BGR
}

func NewIDisplay(info hid.DeviceInfo) *iDisplay {
	d := &iDisplay{
		info:         info,
		buttonCanvas: imageutil.NewBGR(image.Rect(0, 0, 72, 72)),
		canvas:       imageutil.NewBGR(image.Rect(0, 0, 3*72, 5*72)),
	}
	for i := range d.button {
		d.button[i] = newButton(d, i)
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

func (d *iDisplay) DeviceInfo() hid.DeviceInfo   { return d.info }
func (d *iDisplay) Path() string                 { return d.info.Path }
func (d *iDisplay) Manufacturer() string         { return "Infinitton" }
func (d *iDisplay) Product() string              { return d.info.Product }
func (d *iDisplay) Serial() string               { return d.info.Serial }
func (d *iDisplay) Button(int) benjamin.Button   { return nil }
func (d *iDisplay) Buttons() int                 { return 15 }
func (d *iDisplay) ButtonLayout() image.Point    { return image.Pt(3, 5) }
func (d *iDisplay) ButtonSize() image.Point      { return image.Pt(72, 72) }
func (d *iDisplay) Display(int) benjamin.Display { return nil }
func (d *iDisplay) Displays() int                { return 0 }
func (d *iDisplay) DisplaySize() image.Point     { return image.Point{} }
func (d *iDisplay) Encoder(int) benjamin.Encoder { return nil }
func (d *iDisplay) Encoders() int                { return 0 }

func (d *iDisplay) ButtonAt(p image.Point) benjamin.Button {
	if p.X < 0 || p.X >= 3 || p.Y < 0 || p.Y >= 5 {
		return nil
	}

	i := p.Y*3 + p.X
	return d.button[i]
}

func (d *iDisplay) Events() <-chan benjamin.Event {
	c := make(chan benjamin.Event)

	go func(c chan<- benjamin.Event) {
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

func (d *iDisplay) Clear() error {
	i := image.NewNRGBA(image.Rectangle{Max: d.ButtonSize()})
	draw.Draw(i, i.Rect, image.Black, image.Point{}, draw.Src)
	for _, k := range d.button {
		if err := k.SetImage(i); err != nil {
			return err
		}
	}
	return nil
}

func (d *iDisplay) DisplayArea() benjamin.Screen { return nil }
func (d *iDisplay) ButtonArea() benjamin.Screen  { return d }
func (d *iDisplay) Index() int                   { return -1 }
func (d *iDisplay) Size() image.Point            { return d.canvas.Rect.Max }
func (d *iDisplay) Surface() benjamin.Surface    { return d }

func (d *iDisplay) SetImage(i image.Image) error {
	if i == nil {
		return nil
	} else if d.canvas.Rect.Eq(i.Bounds()) {
		draw.Copy(d.canvas, image.Point{}, i, i.Bounds(), draw.Src, nil)
	} else {
		draw.CatmullRom.Scale(d.canvas, d.canvas.Rect, i, i.Bounds(), draw.Src, nil)
	}

	o := 0
	for y := 0; y < 5; y++ {
		r := image.Rectangle{
			Min: image.Point{
				Y: (y + 0) * 72,
			},
			Max: image.Point{
				X: 72,
				Y: (y + 1) * 72,
			},
		}
		for x := 0; x < 3; x++ {
			draw.Copy(d.buttonCanvas, image.Point{}, d.canvas, r, draw.Src, nil)
			r.Min.X += 72
			r.Max.X += 72
			if err := d.button[o].SetImage(d.buttonCanvas); err != nil {
				return err
			}
			o++
		}
	}
	return nil
}

func init() {
	driver.RegisterUSB(vendorID, iDisplayProductID, New)
	driver.RegisterUSB(vendorID, iDisplayProductIDAlt, New)
}
