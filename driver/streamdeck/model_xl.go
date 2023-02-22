package streamdeck

import (
	"image"

	"github.com/tehmaze/benjamin/driver"
)

var XL = Properties{
	Model:               "Stream Deck XL",
	ProductID:           0x006c,
	model:               gen2,
	keys:                32,
	keyLayout:           image.Point{8, 4},
	keySize:             image.Point{96, 96},
	keyDataOffset:       3,
	keyTranslate:        translateLTR(),
	keyImageTransform:   transform(rotate180),
	imagePageSize:       1024,
	imagePageHeaderSize: 8,
}

func init() {
	driver.RegisterUSB(VendorID, XL.ProductID, driverFor(XL))
}
