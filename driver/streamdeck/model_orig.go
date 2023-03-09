package streamdeck

import (
	"image"

	"github.com/tehmaze/benjamin/driver"
)

var Orig = Properties{
	Model:               "Stream Deck",
	ProductID:           0x0060,
	model:               gen1,
	keys:                15,
	keyLayout:           image.Point{5, 3},
	keySize:             image.Point{72, 72},
	keyDataOffset:       1,
	keyTranslate:        translateRTL(5),
	keyImageTransform:   transform(rotate180),
	imageBytes:          convertBMP,
	imagePageSize:       8191,
	imagePageHeaderSize: 16,
}

func init() {
	driver.RegisterUSB(VendorID, Orig.ProductID, driverFor(Orig))
}
