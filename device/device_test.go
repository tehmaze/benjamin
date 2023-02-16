package device_test

import (
	"testing"

	"github.com/tehmaze/benjamin/device"
	_ "github.com/tehmaze/benjamin/device/streamdeck" // Stream Deck drivers
)

func TestDiscover(t *testing.T) {
	devices := device.Discover()
	t.Log("found", len(devices), "supported devices")
	for i, device := range devices {
		t.Logf("device %d:", i+1)
		t.Log("  manufacturer:", device.Manufacturer())
		t.Log("  product:     ", device.Product())
		t.Log("  serial:      ", device.SerialNumber())
	}
}
