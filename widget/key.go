package widget

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/tehmaze/benjamin"
)

func ButtonIcon(key benjamin.Button, i image.Image, effects ...Effect) *Image {
	dim := key.Size()

	return &Image{
		Base:  MakeBase(key, effects...),
		Image: imaging.Resize(i, dim.X, dim.Y, imaging.Lanczos),
	}
}
