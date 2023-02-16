package main

import (
	"fmt"

	"github.com/tehmaze/benjamin/device"

	_ "github.com/tehmaze/benjamin/device/streamdeck"
	_ "github.com/tehmaze/benjamin/device/window"
)

func main() {
	devices := device.Discover()
	if len(devices) == 0 {
		fmt.Println("no compatible devices found")
		return
	}

	fmt.Println(len(devices), "compatible devices found:")
	for i, device := range devices {
		fmt.Println("device", i+1)
		fmt.Println("  +- manufacturer:", device.Manufacturer())
		fmt.Println("  +- product:     ", device.Product())
		fmt.Println("  `- serial:      ", device.SerialNumber())
	}
}
