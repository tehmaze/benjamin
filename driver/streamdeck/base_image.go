package streamdeck

import (
	"bytes"
	"image"
	"image/jpeg"
)

var bmpHeader = []byte{
	0x42, 0x4d, 0xf6, 0x3c, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x36, 0x00, 0x00, 0x00, 0x28, 0x00,
	0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x48, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xc0, 0x3c, 0x00, 0x00, 0xc4, 0x0e,
	0x00, 0x00, 0xc4, 0x0e, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func convertBMP(i *image.NRGBA) ([]byte, error) {
	var b bytes.Buffer
	b.Write(bmpHeader)
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		for x := i.Rect.Min.X; x < i.Rect.Max.X; x++ {
			c := i.NRGBAAt(x, y)
			b.Write([]byte{c.B, c.G, c.R})
		}
	}
	return b.Bytes(), nil
}

func convertJPEG(i image.Image) ([]byte, error) {
	var b bytes.Buffer
	err := jpeg.Encode(&b, i, &jpeg.Options{Quality: 100})
	return b.Bytes(), err
}

type imageTransform interface {
	Transform(*image.NRGBA)
}

type imageTransformFunc func(*image.NRGBA)

func (f imageTransformFunc) Transform(i *image.NRGBA) { f(i) }

type imageTransforms []imageTransform

func (ts imageTransforms) Transform(i *image.NRGBA) {
	for _, t := range ts {
		t.Transform(i)
	}
}

func transform(fs ...imageTransformFunc) imageTransform {
	ts := make(imageTransforms, len(fs))
	for i, f := range fs {
		ts[i] = f
	}
	return ts
}

func flipH(i *image.NRGBA) {
	l := i.Rect.Dx() * 4
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		j := y * i.Stride
		reverse(i.Pix[j : j+l])
	}
}

func flipV(i *image.NRGBA) {
	l := i.Rect.Dx() * 4
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y/2; y++ {
		j := y * i.Stride
		k := (i.Rect.Max.Y - y - 1) * i.Stride
		flip(i.Pix[j:j+l], i.Pix[k:k+l])
	}
}

func rotate180(i *image.NRGBA) {
	l := i.Rect.Dx() * 4
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y/2; y++ {
		j := y * i.Stride
		k := (i.Rect.Max.Y - y - 1) * i.Stride
		flip(reverse(i.Pix[j:j+l]), reverse(i.Pix[k:k+l]))
	}
}

func reverse(pix []uint8) []uint8 {
	if len(pix) <= 4 {
		return pix
	}
	i := 0
	j := len(pix) - 4
	for i < j {
		pi := pix[i : i+4 : i+4]
		pj := pix[j : j+4 : j+4]
		pi[0], pj[0] = pj[0], pi[0]
		pi[1], pj[1] = pj[1], pi[1]
		pi[2], pj[2] = pj[2], pi[2]
		pi[3], pj[3] = pj[3], pi[3]
		i += 4
		j -= 4
	}
	return pix
}

func flip(a, b []uint8) {
	if len(a) != len(b) || len(a) < 4 {
		return
	}
	i := 0
	j := len(a) - 4
	for i < j {
		pa := a[i : i+4 : i+4]
		pb := b[i : i+4 : i+4]
		pa[0], pb[0] = pb[0], pa[0]
		pa[1], pb[1] = pb[1], pa[1]
		pa[2], pb[2] = pb[2], pa[2]
		pa[3], pb[3] = pb[3], pa[3]
		i += 4
	}
}

type imageData struct {
	Data     []byte
	PageSize int
}

func (d imageData) Page(index int) ([]byte, bool) {
	o := index * d.PageSize
	if o >= len(d.Data) {
		return nil, true
	}

	l := d.PageLength(index)
	if o+l >= len(d.Data) {
		l = len(d.Data) - o
	}

	return d.Data[o : o+l], index == d.PageCount()-1
}

func (d imageData) PageLength(index int) int {
	r := len(d.Data) - index*d.PageSize
	if r > d.PageSize {
		return d.PageSize
	}
	if r > 0 {
		return r
	}
	return 0
}

func (d imageData) PageCount() int {
	c := len(d.Data) / d.PageSize
	if len(d.Data)%d.PageSize > 0 {
		c++
	}
	return c
}
