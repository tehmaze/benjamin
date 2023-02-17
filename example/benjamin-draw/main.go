package main

import (
	"image"
	"image/color"
	"log"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/device"
	_ "github.com/tehmaze/benjamin/device/streamdeck" // Stream Deck support
	"github.com/tehmaze/benjamin/widget"
)

func main() {
	d, err := device.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = d.Close() }()

	if err := d.SetBrightness(100); err != nil {
		log.Fatal(err)
	}

	p := benjamin.NewPanel(d, 1, 30)
	for i := 0; i < d.Keys(); i++ {
		w := widget.Color(color.Black)
		if i == 0 {
			w = widget.Color(color.RGBA{R: 0xff, A: 0xff})
		}
		w.Move(image.Pt(i%d.Dim().X, i/d.Dim().X))
		p.Layers[0].AddWidget(w)
	}

	log.Println("press (0,0) to quit")
	on := make(map[image.Point]bool)
	for ev := range d.Events() {
		log.Println(ev)
		if ev.Type == device.KeyPressed && ev.Pos.X == 0 && ev.Pos.Y == 0 {
			log.Println("bye!")
			_ = d.SetBrightness(20)
			_ = d.Reset()
			return
		}
		switch ev.Type {
		case device.KeyPressed:
			if k := d.Key(ev.Pos); k != nil {
				c := color.White
				if on[ev.Pos] {
					c = color.Black
				}
				k.Update(image.NewUniform(c))
				on[ev.Pos] = !on[ev.Pos]
			}
		}
	}
}
