package streamdeck

import (
	"encoding/binary"
	"image"
	"log"
	"time"

	"github.com/tehmaze/benjamin"
	"github.com/tehmaze/benjamin/driver"
)

var Plus = Properties{
	Model:               "Stream Deck Plus",
	ProductID:           0x0084,
	model:               plus,
	displays:            4,
	displaySize:         image.Point{200, 100},
	displayLayout:       image.Pt(4, 1),
	encoders:            4,
	keys:                8,
	keyLayout:           image.Point{4, 2},
	keySize:             image.Point{120, 120},
	keyDataOffset:       3,
	keyTranslate:        translateLTR(),
	imagePageSize:       1024,
	imagePageHeaderSize: 8,
}

type plusModel struct {
	*Device
	*baseModel
}

func plus(device *Device) model {
	return &plusModel{
		Device:    device,
		baseModel: gen2(device).(*baseModel),
	}
}

func (m *plusModel) Handle(p []byte, c chan<- benjamin.Event) {
	switch p[1] {
	case 0x00: // key
		m.baseModel.handleButton(p[1:], c)
	case 0x02: // display
		m.handleDisplay(p[1:], c)
	case 0x03: // encoder
		m.handleEncoder(p[4:], c)
	}
}

func (m *plusModel) handleDisplay(p []byte, c chan<- benjamin.Event) {
	var (
		x  = binary.LittleEndian.Uint16(p[5:])
		y  = binary.LittleEndian.Uint16(p[7:])
		at = image.Pt(int(x), int(y))
	)
	var (
		index   = x / uint16(m.prop.displaySize.X)
		display = m.display[index]
	)
	switch p[3] {
	case 0x01: // short press
		c <- benjamin.NewDisplayPress(m, display, at)
	case 0x02: // long press
		c <- benjamin.NewDisplayLongPress(m, display, at)
	case 0x03: // swipe
		x = binary.LittleEndian.Uint16(p[9:])
		y = binary.LittleEndian.Uint16(p[11:])
		to := image.Pt(int(x), int(y))
		c <- benjamin.NewDisplaySwipe(m, display, at, to)
	default:
		log.Print("display: unknown command", p[3])
	}
}

func (m *plusModel) handleEncoder(p []byte, c chan<- benjamin.Event) {
	switch p[0] {
	case 0x00: // press/release
		state := p[1:]
		for i := 0; i < len(state) && i < m.prop.encoders; i++ {
			var (
				press   = state[i] != 0
				encoder = m.encoder[i]
			)
			if encoder.state != state[i] {
				encoder.state = state[i]
				if press {
					encoder.press = time.Now()
					c <- benjamin.NewEncoderPress(m, encoder)
				} else {
					c <- benjamin.NewEncoderRelease(m, encoder, time.Since(encoder.press))
				}
			}
		}

	case 0x01: // change
		state := p[1:]
		for i := 0; i < len(state) && i < m.prop.encoders; i++ {
			if change := int8(state[i]); change != 0 {
				encoder := m.encoder[i]
				c <- benjamin.NewEncoderChange(m, encoder, int(change), 8)
			}
		}
	}
}

func init() {
	driver.RegisterUSB(VendorID, Plus.ProductID, driverFor(Plus))
}
