package streamdeck

import (
	"bytes"
	"image"
	"image/jpeg"

	"golang.org/x/image/draw"
)

const streamDeckProductID = 0x0060

var (
	streamDeckRev1Firmware      = []byte{0x04}
	streamDeckRev1Reset         = []byte{0x0b, 0x63}
	streamDeckRev1SetBrightness = []byte{0x05, 0x55, 0xd1, 0x01}
	streamDeckRev2Firmware      = []byte{0x05}
	streamDeckRev2Reset         = []byte{0x03, 0x02}
	streamDeckRev2SetBrightness = []byte{0x03, 0x08}
)

func streamDeckRev1PageHeader(pageIndex, keyIndex, payloadLength int, lastPage bool) []byte {
	var lastPageByte byte
	if lastPage {
		lastPageByte = 0x01
	}
	return []byte{
		0x02, 0x01,
		byte(pageIndex + 1), 0x00,
		lastPageByte,
		byte(keyIndex + 1),
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

func streamDeckRev2PageHeader(pageIndex, keyIndex, payloadLength int, lastPage bool) []byte {
	var lastPageByte byte
	if lastPage {
		lastPageByte = 0x01
	}
	return []byte{
		0x02, 0x07,
		byte(keyIndex),
		lastPageByte,
		byte(payloadLength),
		byte(payloadLength >> 8),
		byte(pageIndex),
		byte(pageIndex >> 8),
	}
}

func toRGBA(i image.Image) *image.RGBA {
	switch i := i.(type) {
	case *image.RGBA:
		return i
	}
	o := image.NewRGBA(i.Bounds())
	draw.Copy(o, image.Point{}, i, i.Bounds(), draw.Src, nil)
	return o
}

var bmpHeader = []byte{
	0x42, 0x4d, 0xf6, 0x3c, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x36, 0x00, 0x00, 0x00, 0x28, 0x00,
	0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x48, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xc0, 0x3c, 0x00, 0x00, 0xc4, 0x0e,
	0x00, 0x00, 0xc4, 0x0e, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func toBMP(i image.Image) ([]byte, error) {
	var (
		r = i.Bounds()
		b = make([]byte, len(bmpHeader)+r.Dx()*r.Dy())
		s = toRGBA(i)
	)
	copy(b, bmpHeader)

	o := len(bmpHeader)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		// flip image horizontally
		for x := r.Max.X - 1; x >= r.Min.X; x-- {
			c := s.RGBAAt(x, y)
			b[o+0] = c.B
			b[o+1] = c.G
			b[o+2] = c.R
			o += 3
		}
	}

	return b, nil
}

func toJPEG(i image.Image) ([]byte, error) {
	// flip image horizontally and vertically
	var (
		f  = image.NewRGBA(i.Bounds())
		r  = i.Bounds()
		dx = r.Dx()
		dy = r.Dy()
	)
	draw.Copy(f, image.Point{}, i, r, draw.Src, nil)
	for y := 0; y < dy/2; y++ {
		yy := r.Max.Y - y - 1
		for x := 0; x < dx; x++ {
			xx := r.Max.X - x - 1
			c := f.RGBAAt(x, y)
			f.SetRGBA(x, y, f.RGBAAt(xx, yy))
			f.SetRGBA(xx, yy, c)
		}
	}

	var b bytes.Buffer
	if err := jpeg.Encode(&b, f, &jpeg.Options{Quality: 100}); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
