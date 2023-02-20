package loupedeck

import "github.com/tehmaze/benjamin/deck"

func init() {
	var LoupedeckLive = &DeviceType{
		vendorID:  loupedeckVendorID,
		productID: 0x0004,
		name:      "Loupedeck Live",
		cols:      4,
		rows:      3,
		buttons:   8,
	}

	deck.RegisterUSB(
		LoupedeckLive.productID,
		LoupedeckLive.vendorID,
		LoupedeckLive.Driver)
}
