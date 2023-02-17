package device

import (
	"image"

	"golang.org/x/image/draw"
)

// Surface with keys
type Surface interface {
	// Dim is the number of key columns and rows.
	Dim() image.Point

	// Key returns the Key at index. Returns nil if the key doesn't exist.
	Key(image.Point) Key

	// Keys is the number of keys on the device.
	Keys() int

	// KeySize of the key in pixels.
	KeySize() image.Point

	// Margin is the number of pixels between keys.
	Margin() image.Point

	// SetBrightness updates the brightness, range 0-100.
	SetBrightness(uint8) error
}

type BackgroundSurface struct {
	Surface
	Background image.Image
}

// WithBackground returns a Surface that has a default background.
func WithBackground(surface Surface, background image.Image) *BackgroundSurface {
	return &BackgroundSurface{
		Surface:    surface,
		Background: background,
	}
}

func (s *BackgroundSurface) Key(p image.Point) Key {
	if k := s.Surface.Key(p); k != nil {
		return &backgroundKey{
			Key:     k,
			surface: s,
		}
	}
	return nil
}

type backgroundKey struct {
	Key
	surface *BackgroundSurface
	buf     [2]*image.RGBA
}

func newBackgroundKey(key Key, s *BackgroundSurface) *backgroundKey {
	var (
		size = key.Size()
		pos  = key.Position()
		k    = &backgroundKey{
			Key:     key,
			surface: s,
			buf: [2]*image.RGBA{
				image.NewRGBA(image.Rectangle{Max: size}),
				image.NewRGBA(image.Rectangle{Max: size}),
			},
		}
		origin image.Point
	)

	// Calculate origin in backgroundSurface.Background
	origin = origin.Add(size)
	origin.X *= pos.X
	origin.Y *= pos.Y
	draw.Copy(k.buf[0], image.Point{}, s.Background, image.Rectangle{
		Min: origin,
		Max: origin.Add(size),
	}, draw.Src, nil)

	return k
}

func (k *backgroundKey) Surface() Surface {
	return k.surface
}

func (k *backgroundKey) Update(i image.Image) error {
	draw.Copy(k.buf[1], image.Point{}, k.buf[0], k.buf[0].Rect, draw.Src, nil)
	draw.Copy(k.buf[1], image.Point{}, i, i.Bounds(), draw.Over, nil)
	return k.Key.Update(k.buf[1])
}