package loupedeck

import "github.com/tehmaze/benjamin/deck"

func init() {
	var LoupedeckLiveS = &DeviceType{
		vendorID:  loupedeckVendorID,
		productID: 0x0006,
		name:      "Loupedeck Live S",
		cols:      5,
		rows:      3,
		buttons:   4,
	}

	deck.RegisterUSB(
		LoupedeckLiveS.productID,
		LoupedeckLiveS.vendorID,
		LoupedeckLiveS.Driver)
}
