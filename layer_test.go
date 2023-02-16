package benjamin_test

import (
	"image"
	"testing"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/device/dummy"
)

func TestNewLayers(t *testing.T) {
	d := new(dummy.Dummy)
	p := benjamin.NewPanel(d, 3)

	if n := len(p.Layers); n != 3 {
		t.Errorf("expected 3 layers, got %d", n)
	}

	if i := p.Render(1, 1); i != nil {
		t.Errorf("expected nil image for empty panel, got %T", i)
	}

	p.Layers[1].Widgets[6] = &benjamin.Icon{
		Image: image.NewRGBA(image.Rectangle{Max: dummy.KeySize}),
	}
	if i := p.Render(1, 1); i == nil {
		t.Errorf("expected Icon for (1, 1), got %T", i)
	}
}
