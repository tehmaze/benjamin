package streamdeck

import (
	"image"
	"time"

	"github.com/tehmaze/benjamin"
	"golang.org/x/image/draw"
)

var (
	// DisplayImageInterpolator is the default interpolator for Display images.
	DisplayImageInterpolator draw.Interpolator = draw.CatmullRom

	// ButtonImageInterpolator is the default interpolator for Button images.
	ButtonImageInterpolator draw.Interpolator = draw.BiLinear
)

var (
	// blank image, used for clearing the button/display.
	blank = image.NewNRGBA(image.Rect(0, 0, 800, 100))
)

type display struct {
	device *Device
	index  int
	image  *image.NRGBA
}

func newDisplay(device *Device, index int) *display {
	return &display{
		device: device,
		index:  index,
		image:  image.NewNRGBA(image.Rectangle{Max: device.prop.displaySize}),
	}
}

func (d *display) Surface() benjamin.Surface {
	return d.device
}

func (d *display) Index() int {
	return d.index
}

func (d *display) Size() image.Point {
	return d.device.prop.displaySize
}

func (d *display) Position() image.Point {
	return image.Pt(0, d.index)
}

func (d *display) SetImage(i image.Image) error {
	if i == nil {
		i = blank
	}

	// Fill our key image with the new image.
	if o, ok := i.(*image.RGBA); ok && o.Rect.Eq(d.image.Rect) {
		// Fast path, copy pixels.
		copy(d.image.Pix, o.Pix)
	} else {
		// Interpolate image into key image.
		ButtonImageInterpolator.Scale(d.image, d.image.Rect, i, i.Bounds(), draw.Src, nil)
	}

	// Copy to general display buffer area.
	draw.Copy(d.device.displayImage, image.Pt(d.device.prop.displaySize.X*d.index, 0), d.image, d.image.Rect, draw.Src, nil)

	b, err := convertJPEG(d.device.displayImage)
	if err != nil {
		return err
	}

	return d.device.SetDisplayImage(b)
}

type encoder struct {
	device *Device
	index  int
	state  byte
	press  time.Time
}

func newEncoder(device *Device, index int) *encoder {
	return &encoder{
		device: device,
		index:  index,
	}
}

func (e *encoder) Surface() benjamin.Surface {
	return e.device
}

func (e *encoder) Index() int {
	return e.index
}

func (e *encoder) Display() benjamin.Display {
	if e.index >= e.device.prop.displays {
		return nil
	}
	return e.device.display[e.index]
}

type key struct {
	device *Device
	index  int
	pos    image.Point
	state  byte
	press  time.Time
	image  *image.NRGBA
}

func newButton(device *Device, index int) *key {
	return &key{
		device: device,
		pos:    image.Pt(index%device.prop.keyLayout.X, index/device.prop.keyLayout.X),
		index:  index,
		image:  image.NewNRGBA(image.Rectangle{Max: device.prop.keySize}),
	}
}

func (k *key) Surface() benjamin.Surface {
	return k.device
}

func (k *key) Index() int {
	return k.index
}

func (k *key) Position() image.Point {
	return k.pos
}

func (k *key) Size() image.Point {
	return k.device.prop.keySize
}

func (k *key) SetImage(i image.Image) error {
	// Fill our key image with the new image.
	if i == nil {
		i = blank
	} else if i != k.image {
		if o, ok := i.(*image.RGBA); ok && o.Rect.Eq(k.image.Rect) {
			// Fast path, copy pixels.
			copy(k.image.Pix, o.Pix)
		} else {
			// Interpolate image into key image.
			ButtonImageInterpolator.Scale(k.image, k.image.Rect, i, i.Bounds(), draw.Src, nil)
		}

		// Apply transformations
		if k.device.prop.keyImageTransform != nil {
			k.device.prop.keyImageTransform.Transform(k.image)
		}
	}

	var (
		b   []byte
		err error
	)
	if k.device.prop.imageBytes != nil {
		b, err = k.device.prop.imageBytes(k.image)
	} else {
		b, err = convertJPEG(k.image)
	}
	if err != nil {
		return err
	}
	return k.device.SetButtonImage(k.index, b)
}

// keyArea is a virtual screen that renders to all buttons
type keyArea struct {
	device *Device
	canvas *image.NRGBA
}

func newKeyArea(device *Device) *keyArea {
	var (
		l = device.prop.keyLayout
		s = device.prop.keySize
		w = l.X * s.X
		h = l.Y * s.Y
	)
	return &keyArea{
		device: device,
		canvas: image.NewNRGBA(image.Rect(0, 0, w, h)),
	}
}

func (s *keyArea) Surface() benjamin.Surface {
	return s.device
}

func (s *keyArea) Index() int {
	return -1
}

func (s *keyArea) Size() image.Point {
	return s.canvas.Rect.Max
}

func (s *keyArea) SetImage(i image.Image) error {
	if i == nil {
		i = blank
	}

	if o, ok := i.(*image.NRGBA); ok && o.Rect.Eq(s.canvas.Rect) {
		copy(s.canvas.Pix, o.Pix)
	} else {
		DisplayImageInterpolator.Scale(s.canvas, s.canvas.Rect, i, i.Bounds(), draw.Src, nil)
	}

	for y := 0; y < s.device.prop.keyLayout.Y; y++ {
		r := image.Rect(0, y*s.device.prop.keySize.Y, s.device.prop.keySize.X, (y+1)*s.device.prop.keySize.Y)
		for x := 0; x < s.device.prop.keyLayout.X; x++ {
			i := y*s.device.prop.keyLayout.X + x
			k := s.device.key[i]
			draw.Copy(k.image, image.Point{}, s.canvas, r, draw.Src, nil)
			if err := k.SetImage(k.image); err != nil {
				return err
			}
			r.Min.X += s.device.prop.keySize.X
			r.Max.X += s.device.prop.keySize.X
		}
	}
	return nil
}

// displayArea is a virtual display that renders to all displays
type displayArea struct {
	device *Device
	canvas *image.NRGBA
}

func newDisplayArea(device *Device) *displayArea {
	var (
		l = device.prop.displayLayout
		s = device.prop.displaySize
		w = l.X * s.X
		h = l.Y * s.Y
	)
	return &displayArea{
		device: device,
		canvas: image.NewNRGBA(image.Rect(0, 0, w, h)),
	}
}

func (s *displayArea) Surface() benjamin.Surface {
	return s.device
}

func (s *displayArea) Index() int {
	return -1
}

func (s *displayArea) Size() image.Point {
	return s.canvas.Rect.Max
}

func (s *displayArea) SetImage(i image.Image) error {
	if s == nil {
		panic("drawing to nil display area!")
	}

	if i == nil {
		i = blank
	}

	// Fill our key image with the new image.
	if o, ok := i.(*image.RGBA); ok && o.Rect.Eq(s.canvas.Rect) {
		// Fast path, copy pixels.
		copy(s.canvas.Pix, o.Pix)
	} else {
		// Interpolate image into key image.
		DisplayImageInterpolator.Scale(s.canvas, s.canvas.Rect, i, i.Bounds(), draw.Src, nil)
	}

	b, err := convertJPEG(s.canvas)
	if err != nil {
		return err
	}

	return s.device.SetDisplayImage(b)
}

var (
	_ benjamin.Screen = (*keyArea)(nil)
	_ benjamin.Screen = (*displayArea)(nil)
)
