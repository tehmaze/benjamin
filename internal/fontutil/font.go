package fontutil

import (
	"io"
	"os"

	"github.com/golang/freetype/truetype"
)

// Builtin fonts
var (
	Roboto     = must(Load("data/Roboto-Regular.ttf"))
	RobotoBold = must(Load("data/Roboto-Bold.ttf"))
)

func Load(name string) (*truetype.Font, error) {
	var (
		r   io.ReadCloser
		b   []byte
		err error
	)
	if r, err = content.Open(name); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		r, err = os.Open(name)
	}
	if err != nil {
		return nil, err
	}
	if b, err = io.ReadAll(r); err != nil {
		_ = r.Close()
		return nil, err
	}
	_ = r.Close()
	return truetype.Parse(b)
}

func must(font *truetype.Font, err error) *truetype.Font {
	if err != nil {
		panic(err)
	}
	return font
}
