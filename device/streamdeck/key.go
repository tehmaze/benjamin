package streamdeck

import (
	"fmt"
	"image"

	"golang.org/x/image/draw"

	"github.com/tehmaze/benjamin/device"
)

type key struct {
	x, y   int
	index  int
	device *StreamDeck
}

func newKey(d *StreamDeck, x, y int) *key {
	return &key{
		x:      x,
		y:      y,
		index:  y*d.cols + x,
		device: d,
	}
}

func (b key) Surface() device.Surface {
	return b.device
}

func (b key) Device() device.Device {
	return b.device
}

func (b key) Position() image.Point {
	return image.Pt(b.x, b.y)
}

func (b key) Size() image.Point {
	return image.Pt(b.device.pixels, b.device.pixels)
}

func (b key) Update(i image.Image) error {
	r := i.Bounds()
	if r.Dx() != b.device.pixels || r.Dy() != b.device.pixels {
		if _, is := i.(*image.Uniform); !is {
			// Resize with Lanczos resampling.
			//log.Printf("streamdeck: resize %s image to %dx%d", r, b.device.pixels, b.device.pixels)
			o := image.NewRGBA(image.Rectangle{Max: b.Size()})
			draw.BiLinear.Scale(o, o.Bounds(), i, i.Bounds(), draw.Src, nil)
			i = o
		}
	}

	p, err := b.device.toImageFormat(i)
	if err != nil {
		return err
	}

	var (
		imageData = streamDeckImageData{
			data:     p,
			pageSize: b.device.imagePageSize - b.device.imagePageHeaderSize,
		}
		data   = make([]byte, b.device.imagePageSize)
		last   bool
		header []byte
	)

	// NB: We need to execute this in locked context, to prevent concurrent writes!
	b.device.mu.Lock()
	defer b.device.mu.Unlock()

	for page := 0; !last; page++ {
		var p []byte
		p, last = imageData.Page(page)
		header = b.device.imagePageHeader(page, b.index, len(p), last)

		copy(data, header)
		copy(data[len(header):], p)

		if _, err = b.device.dev.Write(data); err != nil {
			return fmt.Errorf(logPrefix+": image transfer to button %d failed: %w", b.index, err)
		}
	}

	return nil
}
