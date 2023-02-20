package streamdeck

import (
	"github.com/tehmaze/benjamin/deck"
)

type encoder struct {
	index  int
	device *StreamDeck
}

func newEncoder(d *StreamDeck, index int) *encoder {
	return &encoder{
		index:  index,
		device: d,
	}
}

func (b encoder) Surface() deck.Surface {
	return b.device
}

func (b encoder) Index() int {
	return b.index
}
