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
	//log.Printf("streamdeck: display %d image %s", d.index, i.Bounds())

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
