package main

import (
	"embed"
	"flag"
	"image"
	"image/png"
	"log"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/deck"
	"github.com/tehmaze/benjamin/widget"

	_ "github.com/tehmaze/benjamin/deck/all" // All hardware drivers
)

func main() {
	serial := flag.String("serial", "", "use device with this serial number")
	fps := flag.Int("fps", 20, "maximum frame rate")
	flag.Parse()

	d, err := newDeck(*serial)
	if err != nil {
		log.Fatal(err)
	}

	defer d.Close()
	if err = d.Reset(); err != nil {
		log.Fatal(err)
	}

	if err = d.SetBrightness(1); err != nil {
		log.Fatal(err)
	}

	p := benjamin.NewPanel(d, 1, *fps)
	addKeys(p)

	for event := range d.Events() {
		log.Println(event)
		switch event.Type {
		case deck.EventTypeKeyPressed:
			j := event.Data.(deck.KeyPress).Position()
			l := p.Layers[0]
			if w := l.Widget(j.X, j.Y); w != nil {
				w.(*widget.TileWidget).SetImage(starStruck)
			}
		case deck.EventTypeKeyReleased:
			j := event.Data.(deck.KeyRelease).Position()
			l := p.Layers[0]
			if w := l.Widget(j.X, j.Y); w != nil {
				w.(*widget.TileWidget).SetImage(relievedFace)
			}
		case deck.EventTypeEncoderPressed:
		case deck.EventTypeEncoderReleased:
		}
	}
}

func newDeck(serial string) (d deck.Deck, err error) {
	if serial == "" {
		return deck.Open()
	}

	for _, d = range deck.Discover() {
		if d.SerialNumber() == serial {
			return d, d.Open()
		}
	}

	return nil, deck.ErrNotFound
}

func addKeys(p *benjamin.Panel) {
	dim := p.Dim()
	for y := 0; y < dim.Y; y++ {
		for x := 0; x < dim.X; x++ {
			r := image.Rectangle{Min: image.Pt(x, y)}
			r.Max = r.Min.Add(image.Pt(1, 1))
			w := widget.Tile(relievedFace)
			w.Rect = r
			p.Layers[0].AddWidget(w)
		}
	}

	for i, l := 0, p.Encoders(); i < l; i++ {
		//e := p.Encoder(i)
		d := p.Display(i)
		d.Update(starStruck)
	}
}

//go:embed *.png
var content embed.FS

var (
	starStruck   = mustImage("star-struck.png")
	relievedFace = mustImage("relieved-face.png")
	gopher       = mustImage("gopher.png")
)

func mustImage(name string) image.Image {
	f, err := content.Open(name)
	if err != nil {
		panic(err)
	}
	i, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return i
}
