//go:build exclude

package main

import (
	"embed"
	"image"
	"image/gif"
	"log"
	"time"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/device"
	_ "github.com/tehmaze/benjamin/device/streamdeck" // Stream Deck support
	"golang.org/x/image/draw"
)

//go:embed go.gif
var content embed.FS

func main() {
	f, err := content.Open("go.gif")
	if err != nil {
		log.Fatalln(err)
	}

	g, err := gif.DecodeAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("gif with", len(g.Image), "frames")

	d, err := device.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = d.Close() }()

	if err = d.Reset(); err != nil {
		log.Fatalln(err)
	}

	go func(d device.Device, g *gif.GIF) {
		var (
			keySize = d.KeySize()
			dim     = d.Dim()
			p       = benjamin.NewPanel(d, 1, 15)
			t       []*benjamin.Tile
		)
		for y := 0; y < dim.Y; y++ {
			for x := 0; x < dim.X; x++ {
				w := &benjamin.Tile{
					Base:    benjamin.Base{Rect: image.Rect(x, y, x+1, y+1)},
					Texture: image.NewRGBA(image.Rectangle{Max: keySize}),
				}
				t = append(t, w)
				p.Layers[0].AddWidget(w)
			}
		}
		log.Printf("animate: %d tiles created", len(t))

		var (
			canvas    = benjamin.ImageForSurface(nil, d)
			canvasMax = canvas.Bounds().Max
			frameMax  = g.Image[0].Rect.Bounds().Max
			fx        = float64(frameMax.X) / float64(canvasMax.X)
			fy        = float64(frameMax.Y) / float64(canvasMax.Y)
		)
		for {
			// Reset canvas
			draw.Draw(canvas, canvas.Rect, image.Transparent, image.Point{}, draw.Src)

			// Render frames
			log.Println("animate", len(g.Image), "frames")
			for j, frame := range g.Image {
				var (
					dmin = image.Pt(
						int(float64(frame.Rect.Min.X)/fx),
						int(float64(frame.Rect.Min.Y)/fy),
					)
					dmax = image.Pt(
						int(float64(frame.Rect.Max.X)/fx),
						int(float64(frame.Rect.Max.Y)/fy),
					)
					dr = image.Rectangle{Min: dmin, Max: dmax}
				)
				//log.Printf("render %s frame to %s onto %s canvas", frame.Rect, dr, canvas.Rect)
				draw.CatmullRom.Scale(canvas, dr, frame, frame.Rect, draw.Src, nil)

				var i int
				for y := 0; y < canvas.Rect.Max.Y; y += keySize.Y {
					for x := 0; x < canvas.Rect.Max.X; x += keySize.X {
						r := image.Rect(x, y, x+keySize.X, y+keySize.Y)
						w := t[i]
						draw.Copy(w.Texture.(*image.RGBA), image.Point{}, canvas, r, draw.Over, nil)
						i++
					}
				}

				time.Sleep(time.Duration(g.Delay[j]) * 10 * time.Millisecond)
			}
		}
	}(d, g)

	log.Println("press (0,0) to quit")
	for ev := range d.Events() {
		log.Println(ev)
		if ev.Type == device.KeyPressed && ev.Pos.X == 0 && ev.Pos.Y == 0 {
			log.Println("bye!")
			return
		}
	}
}
