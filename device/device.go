package device

import (
	"fmt"
	"image"

	"github.com/karalabe/hid"
)

// Device hardware.
type Device interface {
	// Open the underlying hardware interface.
	Open() error

	// Close the underlying hardware interface.
	Close() error

	// Reset the device.
	Reset() error

	// Product name.
	Product() string

	// Manufacturer of the device.
	Manufacturer() string

	// SerialNumber of the device.
	SerialNumber() string

	// Dim is the number of key columns and rows.
	Dim() image.Point

	// Key returns the Key at index. Returns nil if the key doesn't exist.
	Key(int) Key

	// Keys is the number of keys on the device.
	Keys() int

	// KeySize of the key in pixels.
	KeySize() image.Point

	// Events returns a channel that contains all Key events.
	Events() <-chan Event
}

// Key on a device.
type Key interface {
	// Device the key is connected to.
	Device() Device

	// Position opn the Device.
	Position() image.Point

	// Size of the key in pixels.
	Size() image.Point

	// Update the key graphics.
	Update(image.Image) error
}

// Driver returns a new Device.
type Driver func() Device

type deviceDriver struct {
	Detect func() bool
	Driver Driver
}

var drivers []deviceDriver

// Register a new driver.
func Register(detect func() bool, driver Driver) {
	drivers = append(drivers, deviceDriver{
		Detect: detect,
		Driver: driver,
	})
}

// USBDriver returns the hardware USB driver for a Device.
type USBDriver func(hid.DeviceInfo) Device

// drivers holds all registered device drivers.
var usbDrivers = make(map[uint32]USBDriver)

// RegisterUSB registers a new USB device Driver.
func RegisterUSB(vendorID, productID uint16, driver USBDriver) {
	id := uint32(vendorID)<<16 | uint32(productID)
	if d, dupe := usbDrivers[id]; dupe {
		panic(fmt.Sprintf("device driver for %04x:%04x is already registered as %T: %+v", vendorID, productID, d, d))
	}
	usbDrivers[id] = driver
}

// Discover connected devices that we support.
func Discover() []Device {
	var vendorIDs = make(map[uint16]bool)
	for id := range usbDrivers {
		vendorIDs[uint16(id>>16)] = true
	}

	var devices []Device
	for vendorID := range vendorIDs {
		for _, deviceInfo := range hid.Enumerate(vendorID, 0) {
			id := uint32(deviceInfo.VendorID)<<16 | uint32(deviceInfo.ProductID)
			if driver, ok := usbDrivers[id]; ok {
				devices = append(devices, driver(deviceInfo))
			}
		}
	}

	for _, driver := range drivers {
		if driver.Detect() {
			devices = append(devices, driver.Driver())
		}
	}

	return devices
}
