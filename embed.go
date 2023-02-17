package benjamin

import (
	"embed"
	"io"

	"github.com/golang/freetype/truetype"
)

// Embedded fonts
var (
	Roboto     = embeddedFont("Roboto-Regular.ttf")
	RobotoBold = embeddedFont("Roboto-Bold.ttf")
)

//go:embed data/font/*.ttf
var content embed.FS

func embeddedFont(name string) *truetype.Font {
	r, err := content.Open("data/font/" + name)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	_ = r.Close()

	f, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}

	return f
}
