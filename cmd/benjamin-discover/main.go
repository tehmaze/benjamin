package main

import (
	"fmt"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/driver"
	_ "github.com/tehmaze/benjamin/driver/all" // All hardware drivers
)

func main() {
	devices := driver.Scan()
	if len(devices) == 0 {
		fmt.Println("no compatible devices found")
		return
	}

	fmt.Println(len(devices), "compatible devices found:")
	for i, device := range devices {
		if usb, ok := device.(benjamin.USBDevice); ok {
			info := usb.DeviceInfo()
			fmt.Printf("device %d (usb id %04x:%04x)\n", i+1, info.VendorID, info.ProductID)
			fmt.Println("  +- path:        ", info.Path)
		} else {
			fmt.Println("device", i+1)
		}
		fmt.Println("  +- manufacturer:", device.Manufacturer())
		fmt.Println("  +- product:     ", device.Product())
		fmt.Println("  `- serial:      ", device.Serial())
	}
}
