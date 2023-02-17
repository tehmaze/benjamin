package benjamin_test

import (
	"testing"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/device/dummy"
)

func TestNewLayers(t *testing.T) {
	d := new(dummy.Dummy)
	p := benjamin.NewPanel(d, 3, 0)

	if n := len(p.Layers); n != 3 {
		t.Errorf("expected 3 layers, got %d", n)
	}
}
