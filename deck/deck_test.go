package deck_test

import (
	"testing"

	"github.com/tehmaze/benjamin/deck"
	_ "github.com/tehmaze/benjamin/deck/streamdeck" // Stream Deck drivers
)

func TestDiscover(t *testing.T) {
	devices := deck.Discover()
	t.Log("found", len(devices), "supported decks")
	for i, device := range devices {
		t.Logf("device %d:", i+1)
		t.Log("  manufacturer:", device.Manufacturer())
		t.Log("  product:     ", device.Product())
		t.Log("  serial:      ", device.SerialNumber())
	}
}
