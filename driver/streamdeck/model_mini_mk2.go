package streamdeck

import (
	"image"

	"github.com/tehmaze/benjamin/driver"
)

var MiniMK2 = Properties{
	Model:               "Stream Deck Mini MK.2",
	ProductID:           0x0090,
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
	driver.RegisterUSB(VendorID, MiniMK2.ProductID, driverFor(MiniMK2))
}
