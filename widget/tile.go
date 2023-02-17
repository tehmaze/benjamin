package widget

import (
	"image"

	"github.com/tehmaze/benjamin/device"
)

// TileWidget is a single button image, it can be used across multiple keys in which
// case all of the keys will have the same texture.
type TileWidget struct {
	Base
	Texture image.Image
}

func Tile(i image.Image) *TileWidget {
	return &TileWidget{
		Base:    makeBase(),
		Texture: i,
	}
}

func (TileWidget) Dim() image.Point { return image.Pt(1, 1) }

func (w *TileWidget) ImageFor(key device.Key) image.Image {
	if key.Position().In(w.Bounds()) {
		w.IsClean = true
		return w.Texture
	}
	return nil
}

func (w *TileWidget) SetImage(i image.Image) error {
	w.Texture = i
	w.IsClean = false
	return nil
}
