package benjamin

import (
	"image"
	"log"
	"time"

	"github.com/tehmaze/benjamin/deck"
)

const DefaultMaxFPS = 60

// Panel can have multiple layers and refreshes at a fixed frame rate.
type Panel struct {
	deck.Deck
	Background image.Image
	Layers     Layers
	stop       chan struct{}
}

func NewPanel(deck deck.Deck, layers, maxFPS int) *Panel {
	p := &Panel{
		Deck:   deck,
		Layers: make(Layers, layers),
		stop:   make(chan struct{}, 1),
	}
	for i := 0; i < layers; i++ {
		p.Layers[i] = NewLayer(deck)
	}
	go p.update(maxFPS)
	return p
}

func (p *Panel) Stop() {
	select {
	case p.stop <- struct{}{}:
	default:
	}
}

func (p *Panel) update(maxFPS int) {
	if maxFPS < 1 {
		maxFPS = DefaultMaxFPS
	}

	t := time.NewTicker(time.Second / time.Duration(maxFPS))
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := p.Refresh(false); err != nil {
				log.Println("panel: refresh error:", err)
			}
		case <-p.stop:
			return
		}

	}
}

// Refresh all widgets.
func (p *Panel) Refresh(force bool) error {
	if force || p.Layers.UpdateRequired() {
		if p.Background != nil {
			return p.Layers.Refresh(&deck.BackgroundSurface{
				Surface:    p,
				Background: p.Background,
			}, force)
		}
		return p.Layers.Refresh(p, force)
	}
	return nil
}
