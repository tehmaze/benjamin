package streamdeck

import (
	"encoding/binary"
	"fmt"
	"image"
	"log"
	"math"
	"sync"
	"time"

	"github.com/tehmaze/benjamin/deck"
	"github.com/tehmaze/benjamin/internal/hid"
)

const (
	elgatoVendorID = 0x0fd9
	logPrefix      = "benjamin.deck.streamdeck"
)

type deviceType struct {
	productID            uint16
	name                 string
	dim                  image.Point
	cols, rows           int
	keys                 int
	pixels               int
	buttons              int
	encoders             int
	displays             int
	displaySize          image.Point
	margin               int
	dpi                  int
	padding              int
	featureReportSize    int
	firmwareOffset       int
	keyStateOffset       int
	translateKey         func(index, cols uint8) uint8
	handleInput          func(*StreamDeck, []byte, chan<- deck.Event)
	imagePageSize        int
	imagePageHeaderSize  int
	toImageFormat        func(image.Image) ([]byte, error)
	imagePageHeader      func(pageIndex, keyIndex, payloadLength int, lastPage bool) []byte
	commandFirmware      []byte
	commandReset         []byte
	commandSetBrightness []byte
}

func (t deviceType) Driver(info hid.DeviceInfo) deck.Deck {
	return New(info, t)
}

type StreamDeck struct {
	deviceType
	mu           sync.Mutex
	dev          *hid.Device
	info         hid.DeviceInfo
	display      []*display
	encoder      []*encoder
	encoderState []byte
	encoderPress []time.Time
	key          []*key
	keyState     []byte
	keyPress     []time.Time
	keyTrigger   []time.Time
}

func New(info hid.DeviceInfo, t deviceType) *StreamDeck {
	d := &StreamDeck{
		deviceType:   t,
		info:         info,
		display:      make([]*display, t.displays),
		encoder:      make([]*encoder, t.encoders),
		encoderState: make([]byte, t.encoders),
		encoderPress: make([]time.Time, t.encoders),
		key:          make([]*key, t.keys),
		keyState:     make([]byte, t.keys),
		keyPress:     make([]time.Time, t.keys),
		keyTrigger:   make([]time.Time, t.keys),
	}
	for i := range d.display {
		d.display[i] = newDisplay(d, i)
	}
	for i := range d.encoder {
		d.encoder[i] = newEncoder(d, i)
	}
	for i := range d.key {
		d.key[i] = newKey(d, i%t.cols, i/t.cols)
	}
	return d
}

func (d *StreamDeck) Open() (err error) {
	d.dev, err = d.info.Open()
	return
}

func (d *StreamDeck) Close() (err error) {
	if d.dev != nil {
		err = d.dev.Close()
	}
	return
}

func (d *StreamDeck) Name() string {
	return d.name
}

func (d *StreamDeck) Path() string {
	return d.info.Path
}

func (d *StreamDeck) Dim() image.Point {
	return image.Pt(d.cols, d.rows)
}

func (StreamDeck) Button(int) deck.Button {
	// TODO
	return nil
}

func (d *StreamDeck) Buttons() int {
	return d.buttons
}

func (d *StreamDeck) Display(i int) deck.Display {
	if i < 0 || i >= d.displays {
		return nil
	}
	return d.display[i]
}

func (d *StreamDeck) Displays() int {
	return d.displays
}

func (d *StreamDeck) DisplaySize() image.Point {
	return d.displaySize
}

func (d *StreamDeck) Key(p image.Point) deck.Key {
	i := p.Y*d.cols + p.X
	if i < 0 || i >= len(d.key) {
		return nil
	}
	return d.key[i]
}

func (d *StreamDeck) KeySize() image.Point {
	return image.Point{X: d.pixels, Y: d.pixels}
}

func (d *StreamDeck) Encoder(i int) deck.Encoder {
	if i < 0 || i >= d.encoders {
		return nil
	}
	return d.encoder[i]
}

func (d *StreamDeck) Encoders() int {
	return d.encoders
}

func (d *StreamDeck) Margin() image.Point {
	return image.Point{X: d.margin, Y: d.margin}
}

func (d *StreamDeck) Reset() error {
	return d.sendFeatureReport(d.commandReset)
}

func (d *StreamDeck) Version() string {
	r, err := d.getFeatureReport(d.commandFirmware)
	if err != nil {
		return ""
	}
	return string(r[d.firmwareOffset:])
}

func (d *StreamDeck) Product() string {
	return d.info.Product
}

func (d *StreamDeck) Manufacturer() string {
	return d.info.Manufacturer
}

func (d *StreamDeck) ID() string {
	return fmt.Sprintf("%04x:%04x", d.info.VendorID, d.info.ProductID)
}

func (d *StreamDeck) SerialNumber() string {
	return d.info.Serial
}

func (d *StreamDeck) Keys() int {
	return d.keys
}

func (d *StreamDeck) Events() <-chan deck.Event {
	var (
		c = make(chan deck.Event)
		b = make([]byte, d.keyStateOffset+d.keys)
	)
	if d.encoders > 0 && len(b) < d.keyStateOffset+12 {
		b = make([]byte, d.keyStateOffset+12)
	}
	go func(c chan<- deck.Event) {
		defer close(c)

		// Read events from device
		for {
			copy(d.keyState, b[d.keyStateOffset:])

			var err error
			if _, err = d.dev.Read(b); err == nil {
				if d.handleInput != nil {
					d.handleInput(d, b, c)
				} else {
					d.handleKeyInput(b, c)
				}
			}

			if err != nil {
				select {
				case c <- deck.ErrorEvent(err):
				default:
				}
				close(c)
				return
			}
		}
	}(c)

	return c
}

func (d *StreamDeck) handleKeyInput(b []byte, c chan<- deck.Event) {
	for i, l := d.keyStateOffset, d.keyStateOffset+d.keys; i < l && i < len(b); i++ {
		j := uint8(i - d.keyStateOffset)
		if d.translateKey != nil {
			j = d.translateKey(j, uint8(d.cols))
		}
		log.Printf("key %d at index %d", j, i)
		if b[i] != d.keyState[j] {
			if pressed := b[i] == 1; pressed {
				// Press action immediately triggers a press
				d.keyPress[j] = time.Now()
				d.keyTrigger[j] = d.keyPress[j]
				c <- deck.KeyPressEvent(d.key[j])
			} else {
				c <- deck.KeyReleaseEvent(d.key[j], time.Since(d.keyPress[j]))
			}
		}
	}
}

func (d *StreamDeck) handleLCDInput(b []byte, c chan<- deck.Event) {
	var (
		pos = image.Point{
			X: int(binary.LittleEndian.Uint16(b[6:])),
			Y: int(binary.LittleEndian.Uint16(b[8:])),
		}
		index = pos.X / d.displaySize.X
	)
	log.Printf("lcd input index %d, pos %s", index, pos)
	switch b[4] {
	case 0x01: // Short press
		c <- deck.TouchEvent(d.display[index], pos)
	case 0x02: // Long press
		c <- deck.TouchEvent(d.display[index], pos)
	case 0x03: // Swipe
		to := image.Point{
			X: int(binary.LittleEndian.Uint16(b[10:])),
			Y: int(binary.LittleEndian.Uint16(b[12:])),
		}
		c <- deck.SwipeEvent(d.display[index], pos, to)
	}
}

func (d *StreamDeck) handleEncoderInput(b []byte, c chan<- deck.Event) {
	switch b[4] {
	case 0x00: // Press/release
		for i := 0; i < d.encoders; i++ {
			j := i + d.keyStateOffset + 1
			if b[j] != d.encoderState[i] {
				if pressed := b[j] == 1; pressed {
					// Press action immediately triggers a press
					d.encoderPress[i] = time.Now()
					c <- deck.EncoderPressEvent(d.encoder[i])
				} else {
					c <- deck.EncoderReleaseEvent(d.encoder[i], time.Since(d.encoderPress[i]))
				}
			}
		}
		copy(d.encoderState, b[d.keyStateOffset+1:])
	case 0x01: // Change
		for i := 0; i < d.encoders; i++ {
			j := i + d.keyStateOffset + 1
			if v := int8(b[j]); v != 0 {
				c <- deck.EncoderChangeEvent(d.encoder[i], int(v), math.MaxInt8)
			}
		}
	}
}

func (d *StreamDeck) SetBrightness(v float64) error {
	if v < 0.0 {
		v = 0.0
	} else if v > 1.0 {
		v = 1.0
	}
	return d.sendFeatureReport(append(d.commandSetBrightness, uint8(v*100)))
}

func (d *StreamDeck) getFeatureReport(p []byte) ([]byte, error) {
	b := make([]byte, d.featureReportSize)
	copy(b, p)
	d.mu.Lock()
	_, err := d.dev.GetFeatureReport(b)
	d.mu.Unlock()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (d *StreamDeck) sendFeatureReport(p []byte) error {
	b := make([]byte, d.featureReportSize)
	copy(b, p)
	d.mu.Lock()
	_, err := d.dev.SendFeatureReport(b)
	d.mu.Unlock()
	return err
}

type streamDeckImageData struct {
	data     []byte
	pageSize int
}

func (d streamDeckImageData) Page(index int) ([]byte, bool) {
	o := index * d.pageSize
	if o >= len(d.data) {
		return nil, true
	}

	l := d.pageLength(index)
	if o+l > len(d.data) {
		l = len(d.data) - o
	}

	return d.data[o : o+l], index == d.pageCount()-1
}

func (d streamDeckImageData) pageLength(index int) int {
	r := len(d.data) - index*d.pageSize
	if r > d.pageSize {
		return d.pageSize
	}
	if r > 0 {
		return r
	}
	return 0
}

func (d streamDeckImageData) pageCount() int {
	c := len(d.data) / d.pageSize
	if len(d.data)%d.pageSize > 0 {
		return c + 1
	}
	return c
}

func translateRightToLeft(index, cols uint8) uint8 {
	keyCol := index % cols
	return (index - keyCol) + (cols + 1) - keyCol
}

func init() {
	deck.RegisterUSB(elgatoVendorID, streamDeckProductID, deviceType{
		productID:            streamDeckProductID,
		name:                 "Stream Deck",
		cols:                 5,
		rows:                 3,
		keys:                 15,
		pixels:               72, //
		margin:               24, // 0.75"
		dpi:                  124,
		padding:              16,
		featureReportSize:    17,
		firmwareOffset:       5,
		keyStateOffset:       1,
		translateKey:         translateRightToLeft,
		imagePageSize:        7819,
		imagePageHeaderSize:  16,
		imagePageHeader:      streamDeckRev1PageHeader,
		toImageFormat:        toBMP,
		commandFirmware:      streamDeckRev1Firmware,
		commandReset:         streamDeckRev1Reset,
		commandSetBrightness: streamDeckRev1SetBrightness,
	}.Driver)
}
