package streamdeck

import (
	"image"

	"github.com/tehmaze/benjamin/driver"
)

var MK2 = Properties{
	Model:               "Stream Deck MK.2",
	ProductID:           0x0080,
	model:               gen2,
	keys:                15,
	keyLayout:           image.Point{5, 3},
	keySize:             image.Point{72, 72},
	keyDataOffset:       3,
	keyTranslate:        translateLTR(),
	keyImageTransform:   transform(rotate180),
	imagePageSize:       1024,
	imagePageHeaderSize: 8,
}

func init() {
	driver.RegisterUSB(VendorID, MK2.ProductID, driverFor(MK2))
}
