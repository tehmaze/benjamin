package benjamin

import (
	"github.com/tehmaze/benjamin/device"
)

type Panel struct {
	device.Device

	Layers
}

func NewPanel(device device.Device, layers int) *Panel {
	p := &Panel{
		Device: device,
		Layers: make(Layers, layers),
	}
	for i := 0; i < layers; i++ {
		p.Layers[i] = NewLayer(device)
	}
	return p
}

func (p *Panel) Refresh(force bool) error {
	for i, n := 0, p.Keys(); i < n; i++ {
		if force || p.Layers.IsDirtyIndex(i) {
			if err := p.Key(i).Update(p.Layers.RenderIndex(i)); err != nil {
				return err
			}
		}
	}
	return nil
}
