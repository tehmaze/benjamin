package main

import (
	"embed"
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/driver"
	"github.com/tehmaze/benjamin/widget"

	_ "github.com/tehmaze/benjamin/driver/all" // All hardware drivers
)

func main() {
	serial := flag.String("serial", "", "use device with this serial number")
	fps := flag.Int("fps", 25, "maximum frame rate")
	brightness := flag.Float64("brightness", 60, "brightness percentage")
	flag.Parse()

	d, err := newDeck(*serial)
	if err != nil {
		log.Fatal(err)
	}

	defer d.Close()
	if err = d.Reset(); err != nil {
		log.Fatal(err)
	}

	if *brightness < 5 {
		*brightness = 5
	} else if *brightness > 100 {
		*brightness = 100
	}

	if err = d.SetBrightness(*brightness / 100); err != nil {
		log.Fatal(err)
	}

	r := make(benjamin.Router)
	widgets := addKeys(d, r)

	var (
		render = time.NewTicker(time.Second / time.Duration(*fps))
		events = d.Events()
	)
	for {
		select {
		case event := <-events:
			log.Println(event)
			r.Handle(event)
		case t := <-render.C:
			for _, w := range widgets {
				if w.IsUpdated(t) {
					w.Drawable().SetImage(w.Frame(t))
				}
			}
		}
	}
}

func newDeck(serial string) (d benjamin.Device, err error) {
	if serial == "" {
		return driver.Open()
	}

	for _, d = range driver.Scan() {
		if d.Serial() == serial {
			return d, d.Open()
		}
	}

	return nil, driver.ErrNotFound
}

func addKeys(d benjamin.Device, r benjamin.Router) (widgets []widget.Widget) {
	dim := d.KeyLayout()

	log.Println("test: adding", d.Keys(), "keys")
	for y := 0; y < dim.Y; y++ {
		for x := 0; x < dim.X; x++ {
			// log.Printf("test: key (%d,%d)", x, y)
			k := d.KeyAt(image.Pt(x, y))
			if x == 0 && y == 0 {
				w := widget.KeyIcon(k, door)
				r.On(k, benjamin.TypeKeyPress, benjamin.EventHandlerFunc(func(_ benjamin.Event) {
					d.Reset()
					d.SetBrightness(0.2)
					d.Close()
					os.Exit(0)
				}))
				widgets = append(widgets, w)
			} else {
				w := widget.KeyIcon(k, relievedFace)
				r.On(k, benjamin.TypeKeyPress, benjamin.EventHandlerFunc(func(_ benjamin.Event) {
					w.Set(starStruck)
				}))
				r.On(k, benjamin.TypeKeyRelease, benjamin.EventHandlerFunc(func(_ benjamin.Event) {
					w.Set(relievedFace)
				}))
				widgets = append(widgets, w)
			}
		}
	}

	log.Println("test: adding", d.Encoders(), "encoders")
	for i, l := 0, d.Encoders(); i < l; i++ {
		e := d.Encoder(i)
		o := d.Display(i)
		w := widget.NewProgress(o, &widget.ProgressOptions{
			Label: progressNames[i%len(progressIcons)],
			Icon:  progressIcons[i%len(progressIcons)],
		})

		updateColor := func(value float64) {
			if value < 1 {
				w.SetColor(color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0xff})
			} else if value < 50 {
				w.SetColor(color.NRGBA{R: 0x7f, G: 0x7f, B: 0xff, A: 0xff})
			} else {
				w.SetColor(color.NRGBA{R: 0xff, G: 0x7f, B: 0x7f, A: 0xff})
			}
		}

		value := math.Floor(rand.Float64() * 100)
		updateColor(value)
		w.Set(value)

		r.On(e, benjamin.TypeEncoderChange, benjamin.EventHandlerFunc(func(event benjamin.Event) {
			data := event.Data.(benjamin.EncoderChange)
			value := w.Value + float64(data.Change)
			updateColor(value)
			w.Set(value)
			if event.Peripheral.(benjamin.Encoder).Index() == 0 {
				log.Printf("set brightness to %f%%", value)
				if err := d.SetBrightness(float64(value) / 100); err != nil {
					log.Println("set brightness error:", err)
				}
			}
		}))

		prev := w.Value
		r.On(e, benjamin.TypeEncoderPress, benjamin.EventHandlerFunc(func(_ benjamin.Event) {
			if w.Value == 0 {
				updateColor(prev)
				w.Set(prev)
			} else {
				prev = w.Value
				updateColor(0)
				w.Set(0)
			}
		}))
		r.On(o, benjamin.TypeDisplayPress, benjamin.EventHandlerFunc(func(_ benjamin.Event) {
			if w.Value == 0 {
				updateColor(prev)
				w.Set(prev)
			} else {
				prev = w.Value
				updateColor(0)
				w.Set(0)
			}
		}))
		r.On(o, benjamin.TypeDisplayLongPress, benjamin.EventHandlerFunc(func(_ benjamin.Event) {
			updateColor(100)
			w.Set(100)
		}))

		widgets = append(widgets, w)
	}

	return widgets
}

//go:embed data/*.png data/*.ttf
var content embed.FS

var (
	door          = mustImage("door.png")
	starStruck    = mustImage("star-struck.png")
	relievedFace  = mustImage("relieved-face.png")
	gopher        = mustImage("gopher.png")
	home          = mustImage("home.png")
	light         = mustImage("light.png")
	volume        = mustImage("volume.png")
	headset       = mustImage("headset.png")
	progressIcons = []image.Image{light, home, volume, headset}
	progressNames = []string{"Light", "Home", "Volume", "Microphone"}
)

func mustImage(name string) image.Image {
	f, err := content.Open("data/" + name)
	if err != nil {
		panic(err)
	}
	i, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return i
}
