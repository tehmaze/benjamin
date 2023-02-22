package driver

import (
	"errors"
	"fmt"

	"github.com/karalabe/hid"

	"github.com/tehmaze/benjamin"
)

// USBDriver returns a device driver for a USB device.
type USBDriver func(hid.DeviceInfo) benjamin.Device

// Driver returns a device drivers for a device.
type Driver func() benjamin.Device

type deviceDriver struct {
	Detect func() bool
	Driver Driver
}

var (
	ErrNotFound = errors.New("benjamin: no compatible device found")
	usbDrivers  = make(map[uint16]map[uint16]USBDriver)
	drivers     []deviceDriver
)

// Register a driver.
func Register(detect func() bool, driver Driver) {
	drivers = append(drivers, deviceDriver{
		Detect: detect,
		Driver: driver,
	})
}

// RegisterUSB registers a USB driver.
func RegisterUSB(vendorID, productID uint16, driver USBDriver) {
	if _, dupe := usbDrivers[vendorID][productID]; dupe {
		panic(fmt.Sprintf("USB driver for %04x:%04x already registered", vendorID, productID))
	}
	if _, ok := usbDrivers[vendorID]; !ok {
		usbDrivers[vendorID] = make(map[uint16]USBDriver)
	}
	usbDrivers[vendorID][productID] = driver
}

// Scan available devices.
func Scan() []benjamin.Device {
	var available []benjamin.Device

	// Enumerate the USB bus for known drivers.
	for vendorID, devices := range usbDrivers {
		for _, info := range hid.Enumerate(vendorID, 0) {
			if d, ok := devices[info.ProductID]; ok {
				available = append(available, d(info))
			}
		}
	}

	// Enumerate driver detections.
	for _, driver := range drivers {
		if driver.Detect() {
			available = append(available, driver.Driver())
		}
	}

	return available
}

// Open the first available device.
func Open() (benjamin.Device, error) {
	for _, device := range Scan() {
		err := device.Open()
		return device, err
	}
	return nil, ErrNotFound
}
