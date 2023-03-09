package streamdeck

import (
	"image"

	"github.com/tehmaze/benjamin/driver"
)

var Mini = Properties{
	Model:               "Stream Deck Mini",
	ProductID:           0x0063,
	model:               gen1,
	keys:                6,
	keyLayout:           image.Point{3, 2},
	keySize:             image.Point{80, 80},
	keyDataOffset:       1,
	keyTranslate:        translateRTL(5),
	keyImageTransform:   transform(rotate180),
	imageBytes:          convertBMP,
	imagePageSize:       1024,
	imagePageHeaderSize: 16,
}

func init() {
	driver.RegisterUSB(VendorID, Mini.ProductID, driverFor(Mini))
}
