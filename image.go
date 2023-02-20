package benjamin

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/tehmaze/benjamin/deck"
)

// ImageForKey resizes an image to match the Key pixels.
func ImageForKey(i image.Image, k deck.Key) *image.RGBA {
	if k == nil {
		return nil
	}

	o := image.NewRGBA(image.Rectangle{Max: k.Size()})
	if i == nil {
		draw.Draw(o, o.Bounds(), image.Transparent, image.Point{}, draw.Src)
	} else {
		Resize(o, i)
	}

	return o
}

// ImageForSurface resizes an image to match the Surface pixels.
func ImageForSurface(i image.Image, s deck.Surface) *image.RGBA {
	if s == nil {
		return nil
	}

	var (
		dim     = s.Dim()
		keySize = s.KeySize()
		size    = dim.Mul(keySize.X)
	)

	o := image.NewRGBA(image.Rectangle{Max: size})
	if i == nil {
		draw.Draw(o, o.Bounds(), image.Transparent, image.Point{}, draw.Src)
	} else {
		Resize(o, i)
	}

	return o
}

func area(r image.Rectangle) int {
	return r.Dx() * r.Dy()
}

// Resize helper, uses nearest neighbor for scaling down, cat mull rom for scaling up.
func Resize(dst *image.RGBA, src image.Image) {
	var i draw.Interpolator
	if area(dst.Rect) < area(src.Bounds()) {
		// Scaling down, use simple kernel.
		i = draw.NearestNeighbor
	} else {
		i = draw.CatmullRom
	}
	i.Scale(dst, dst.Rect, src, src.Bounds(), draw.Src, nil)
}
