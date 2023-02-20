package streamdeck

import (
	"fmt"
	"image"

	"github.com/tehmaze/benjamin/deck"
	"golang.org/x/image/draw"
)

type display struct {
	index  int
	device *StreamDeck
}

func newDisplay(d *StreamDeck, index int) *display {
	return &display{
		index:  index,
		device: d,
	}
}

func (d *display) Index() int {
	return d.index
}

func (d *display) Size() image.Point {
	return d.device.displaySize
}

func (d *display) Surface() deck.Surface {
	return d.device
}

func (d *display) Update(i image.Image) error {
	o := image.NewRGBA(image.Rectangle{Max: d.device.displaySize})
	draw.BiLinear.Scale(o, o.Rect, i, i.Bounds(), draw.Src, nil)

	p, err := toJPEGVerbatim(o)
	if err != nil {
		return err
	}

	var (
		imageData = streamDeckImageData{
			data:     p,
			pageSize: 1024 - 16,
		}
		data   = make([]byte, 1024) // d.device.imagePageSize)
		r      = image.Rect(d.index*d.device.displaySize.X, 0, (d.index+1)*d.device.displaySize.X, d.device.displaySize.Y)
		last   bool
		header []byte
	)

	// NB: We need to execute this in locked context, to prevent concurrent writes!
	d.device.mu.Lock()
	defer d.device.mu.Unlock()

	for page := 0; !last; page++ {
		var p []byte
		p, last = imageData.Page(page)
		header = displayImagePageHeader(page, r, len(p), last)

		copy(data, header)
		copy(data[len(header):], p)

		if _, err = d.device.dev.Write(data); err != nil {
			return fmt.Errorf(logPrefix+": image transfer to display %d failed: %w", d.index, err)
		}
	}

	return nil
}

func displayImagePageHeader(page int, r image.Rectangle, size int, lastPage bool) []byte {
	var lastPageByte byte
	if lastPage {
		lastPageByte = 0x01
	}
	var (
		x = uint16(r.Min.X)
		y = uint16(r.Min.Y)
		w = uint16(r.Dx())
		h = uint16(r.Dy())
	)
	return []byte{
		0x02, 0x03,
		uint8(x), uint8(x >> 8),
		uint8(y), uint8(y >> 8),
		uint8(w), uint8(w >> 8),
		uint8(h), uint8(h >> 8),
		lastPageByte,
		uint8(page), uint8(page >> 8),
		uint8(size), uint8(size >> 8),
		0x00,
	}
}
