package benjamin

import (
	"image"
	"sync"

	"github.com/tehmaze/benjamin/deck"
)

type Layer struct {
	Device  deck.Deck
	Widgets []Widget
	mu      sync.RWMutex
}

func NewLayer(device deck.Deck) *Layer {
	return &Layer{
		Device: device,
	}
}

func (l *Layer) AddWidget(w Widget) (ok bool) {
	l.mu.Lock()
	l.Widgets = append(l.Widgets, w)
	l.mu.Unlock()
	return
}

func (l *Layer) RemoveWidget(w Widget) (has bool) {
	l.mu.Lock()
	for i, o := range l.Widgets {
		if o == w {
			has = true
			l.Widgets[i] = nil
			l.Widgets = append(l.Widgets[:i], l.Widgets[i+1:]...)
		}
	}
	l.mu.Unlock()
	return
}

func (l *Layer) Widget(x, y int) Widget {
	p := image.Pt(x, y)
	for _, w := range l.Widgets {
		if p.In(w.Bounds()) {
			return w
		}
	}
	return nil
}

func (l *Layer) UpdateRequired() bool {
	for _, w := range l.Widgets {
		if w.UpdateRequired() {
			// log.Println("layer: dirty", w)
			return true
		}
	}
	return false
}

func (l *Layer) Refresh(s deck.Surface, force bool) error {
	seen := make(map[Widget]bool)
	for _, w := range l.Widgets {
		if seen[w] {
			continue
		}

		if force || w.UpdateRequired() {
			r := w.Bounds()
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					k := s.Key(image.Pt(x, y))
					if k == nil {
						continue
					}
					if i := w.ImageFor(k); i != nil {
						if err := k.Update(i); err != nil {
							return err
						}
					}
				}
			}
		}

		seen[w] = true
	}
	return nil
}

// Layers are zero or more stacked layers.
type Layers []*Layer

func (ls Layers) UpdateRequired() bool {
	for _, layer := range ls {
		if layer.UpdateRequired() {
			return true
		}
	}
	return false
}

func (ls Layers) Refresh(s deck.Surface, force bool) error {
	for _, l := range ls {
		if err := l.Refresh(s, force); err != nil {
			return err
		}
	}
	return nil
}
