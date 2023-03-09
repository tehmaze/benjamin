package streamdeck

import (
	"image"

	"github.com/tehmaze/benjamin/driver"
)

var V2 = Properties{
	Model:               "Stream Deck V2",
	ProductID:           0x006d,
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
	driver.RegisterUSB(VendorID, V2.ProductID, driverFor(V2))
}
