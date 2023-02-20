package main

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/deck"
	_ "github.com/tehmaze/benjamin/deck/all" // Stream Deck support
	"github.com/tehmaze/benjamin/widget"
)

func main() {
	d, err := deck.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = d.Close() }()

	r := rainbow(d.Keys() << 1)
	p := benjamin.NewPanel(d, 1, 30)
	for i := 0; i < d.Keys(); i++ {
		w := widget.Color(r[i])
		w.Move(image.Pt(i%d.Dim().X, i/d.Dim().X))
		p.Layers[0].AddWidget(w)
	}

	go func() {
		var (
			n = d.Keys()
			l = len(r)
			w = d.Dim().X
			t = time.NewTicker(time.Second / 20)
			o int
		)
		defer t.Stop()
		for {
			<-t.C
			for i := 0; i < n; i++ {
				var (
					x = i % w
					y = i / w
					w = p.Layers[0].Widgets[i].(*widget.ColorWidget)
				)
				w.Color = r[(x+y+o)%l]
				w.IsClean = false
			}
			o++
		}
	}()

	log.Println("press (0,0) to quit")
	for ev := range d.Events() {
		log.Println(ev)
		switch ev.Type {
		case deck.EventTypeKeyPressed:
			d := ev.Data.(deck.KeyPress)
			if d.Position().X == 0 && d.Position().Y == 0 {
				log.Println("bye!")
				return
			}
		}
	}
}

func rainbow(n int) color.Palette {
	var (
		p = make(color.Palette, n)
		s = 1 / float64(n)
		h float64
	)
	for i := range p {
		p[i] = hslToRGBA(h, 1, .5)
		h += s
	}
	return p
}

func hslToRGBA(h, s, l float64) color.RGBA {
	if s <= 0 {
		// it's gray
		g := uint8(l * 255)
		return color.RGBA{g, g, g, 0xff}
	}

	var v1, v2 float64
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}

	v1 = 2*l - v2

	r := hueToRGB(v1, v2, h+(1.0/3.0))
	g := hueToRGB(v1, v2, h)
	b := hueToRGB(v1, v2, h-(1.0/3.0))

	return color.RGBA{uint8(r * 0xff), uint8(g * 0xff), uint8(b * 0xff), 0xff}
}

func hueToRGB(v1, v2, h float64) float64 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	switch {
	case 6*h < 1:
		return (v1 + (v2-v1)*6*h)
	case 2*h < 1:
		return v2
	case 3*h < 2:
		return v1 + (v2-v1)*((2.0/3.0)-h)*6
	}
	return v1
}
