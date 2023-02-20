package main

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/deck"
	"github.com/tehmaze/benjamin/effect"

	_ "github.com/tehmaze/benjamin/deck/streamdeck" // Stream Deck support
)

func main() {
	d, err := device.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = d.Close() }()

	r := rainbow(d.Keys() << 1)
	p := benjamin.NewPanel(d, 1, 30)
	for i := 0; i < d.Keys(); i++ {
		w := benjamin.NewSolid(r[i])
		w.Move(image.Pt(i%d.Dim().X, i/d.Dim().X))
		p.Layers[0].AddWidget(w)
	}

	for i, w := range p.Layers[0].Widgets {
		p.Layers[0].Widgets[i] = effect.Fade(w, 0.25, time.Second) //  ease.OutElastic)
	}

	log.Println("press (0,0) to quit")
	for ev := range d.Events() {
		log.Println(ev)
		if ev.Type == device.KeyPressed && ev.Pos.X == 0 && ev.Pos.Y == 0 {
			log.Println("bye!")
			return
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
