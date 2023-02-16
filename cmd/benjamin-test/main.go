package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/tehmaze/benjamin/device"

	_ "github.com/tehmaze/benjamin/device/streamdeck"
	_ "github.com/tehmaze/benjamin/device/window"
)

func main() {
	deviceManufacturer := flag.String("manufacturer", "", "filter device by manufacturer")
	deviceProduct := flag.String("product", "", "filter device by product")
	deviceSerial := flag.String("serial", "", "filter device by serial number")
	flag.Parse()

	var deviceFilter []func(d device.Device) bool

	if *deviceManufacturer != "" {
		deviceFilter = append(deviceFilter, func(d device.Device) bool {
			return strings.EqualFold(d.Manufacturer(), *deviceManufacturer)
		})
	}
	if *deviceProduct != "" {
		deviceFilter = append(deviceFilter, func(d device.Device) bool {
			return strings.EqualFold(d.Product(), *deviceProduct)
		})
	}
	if *deviceSerial != "" {
		deviceFilter = append(deviceFilter, func(d device.Device) bool {
			return strings.EqualFold(d.SerialNumber(), *deviceSerial)
		})
	}

	var devices []device.Device
discovering:
	for _, device := range device.Discover() {
		for _, filter := range deviceFilter {
			if !filter(device) {
				continue discovering
			}
		}
		devices = append(devices, device)
	}

	var dev device.Device
	switch l := len(devices); l {
	case 0:
		fmt.Println("no compatible devices found")
		return
	case 1:
		dev = devices[0]
		fmt.Println("using", dev.Manufacturer(), dev.Product(), "with serial", dev.SerialNumber())
	default:
		fmt.Println(l, "devices found, pick one using flags")
		return
	}

	if err := dev.Open(); err != nil {
		fmt.Println("error:", err)
		return
	}
}

func filterAny(_ device.Device) bool { return true }
