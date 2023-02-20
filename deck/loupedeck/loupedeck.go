package loupedeck

import (
	"encoding/binary"
	"image"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/karalabe/hid"

	"github.com/tehmaze/benjamin/deck"
)

var endian = binary.BigEndian

type DeviceType struct {
	vendorID     uint16
	productID    uint16
	name         string
	cols         int
	rows         int
	buttons      int
	encoders     int
	encoderScale int
}

func (t DeviceType) Driver(info hid.DeviceInfo) deck.Deck {
	return New(info, t)
}

type LoupeDeck struct {
	DeviceType
	dev      *hid.Device
	ws       *websocket.Conn
	info     hid.DeviceInfo
	url      string
	key      []*key
	keyPress []time.Time
	encoder  []*encoder
}

func New(info hid.DeviceInfo, t DeviceType) *LoupeDeck {
	d := &LoupeDeck{
		DeviceType: t,
		info:       info,
		key:        make([]*key, t.buttons),
		keyPress:   make([]time.Time, t.buttons),
		encoder:    make([]*encoder, t.encoders),
	}
	for i := range d.key {
		d.key[i] = newKey(d, i)
	}
	return d
}

func (d *LoupeDeck) Open() (err error) {
	if d.info.VendorID > 0 {
		// USB-serial connection
		d.dev, err = d.info.Open()
	} else {
		d.ws, _, err = websocket.DefaultDialer.Dial(d.url, nil)
	}
	return
}

func (d *LoupeDeck) Close() error {
	switch {
	case d.dev != nil:
		return d.dev.Close()
	case d.ws != nil:
		return d.ws.Close()
	default:
		return nil
	}
}

func (d *LoupeDeck) Manufacturer() string {
	if d.info.Manufacturer == "" {
		return "Loupedeck"
	}
	return d.info.Manufacturer
}

func (d *LoupeDeck) Product() string {
	if d.info.Product == "" {
		return d.name
	}
	return d.info.Product
}

func (d *LoupeDeck) Dim() image.Point {
	return image.Point{d.cols, d.rows}
}

func (d *LoupeDeck) Events() <-chan deck.Event {
	c := make(chan deck.Event)
	switch {
	case d.dev != nil:
		go d.usbEvents(d.dev, c)
	case d.ws != nil:
		go d.wsEvents(d.ws, c)
	default:
		close(c)
	}
	return c
}

func (d *LoupeDeck) usbEvents(r *hid.Device, c chan<- deck.Event) {
	defer close(c)

	b := make([]byte, 1024)
	for {
		n, err := r.Read(b)
		if err != nil {
			log.Println("loupedeck: read error:", err)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return
			}
			continue
		}
		d.handle(c, b[:n])
	}
}

func (d *LoupeDeck) wsEvents(r *websocket.Conn, c chan<- deck.Event) {
	defer close(c)
	for {
		_, b, err := r.ReadMessage()
		if err != nil {
			c <- deck.ErrorEvent(err)
			return
		}

		d.handle(c, b)
	}
}

func (d *LoupeDeck) handle(c chan<- deck.Event, b []byte) {
	if len(b) < 2 {
		return
	}

	switch endian.Uint16(b) {
	case resConfirm:
	case resButton:
		var (
			i     = endian.Uint16(b[2:])
			state = b[4]
		)
		if pressed := state == 0; pressed {
			d.keyPress[i] = time.Now()
			c <- deck.KeyPressEvent(d.key[i])
		} else {
			c <- deck.KeyReleaseEvent(d.key[i], time.Since(d.keyPress[i]))
		}
	case resEncoder:
		var (
			index = int(endian.Uint16(b[2:]))
			state = int(b[4])
		)
		c <- deck.EncoderChangeEvent(index, state, d.encoderScale)
	case resTouch:
		var (
			x = int(endian.Uint16(b[4:]))
			y = int(endian.Uint16(b[6:]))
		)
		c <- deck.TouchEvent(image.Pt(x, y))
	case resTouchEnd:
		var (
			x = int(endian.Uint16(b[4:]))
			y = int(endian.Uint16(b[6:]))
		)
		c <- deck.TouchEndEvent(image.Pt(x, y))
	}
}

func (d *LoupeDeck) Key(p image.Point) deck.Key {
	if p.Y == 0 && p.X >= 0 && p.X < d.buttons {
		return newKey(d, p.X)
	}
	return nil
}

type key struct {
	deck  *LoupeDeck
	index int
}

func newKey(deck *LoupeDeck, index int) *key {
	return &key{
		deck:  deck,
		index: index,
	}
}

func (k *key) Position() image.Point {
	return image.Point{X: k.index}
}

var (
	_ deck.Deck    = (*LoupeDeck)(nil)
	_ deck.Surface = (*LoupeDeck)(nil)
	_ deck.Key     = (*key)(nil)
)
