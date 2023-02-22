package streamdeck

import (
	"fmt"
	"image"
	"io"
	"sync"
	"time"

	"github.com/karalabe/hid"

	"github.com/tehmaze/benjamin"
)

// VendorID for Elgate (Corsair) Stream Decks
const VendorID = 0x0fd9

const (
	sendBufferSize = 16
)

type Properties struct {
	ProductID           uint16
	Model               string
	model               func(*Device) model
	displays            int         //
	displayLayout       image.Point // in cols x rows
	displaySize         image.Point // in pixels
	encoders            int         //
	keys                int         //
	keyLayout           image.Point // in cols x rows
	keySize             image.Point // in pixels
	keyDataOffset       int
	keyTranslate        func(int) int
	keyImageTransform   imageTransform
	imageBytes          func(*image.NRGBA) ([]byte, error)
	imagePageSize       int
	imagePageHeaderSize int
}

func driverFor(p Properties) func(hid.DeviceInfo) benjamin.Device {
	return func(info hid.DeviceInfo) benjamin.Device {
		return New(info, p)
	}
}

func New(info hid.DeviceInfo, prop Properties) *Device {
	d := &Device{
		prop:         prop,
		info:         info,
		display:      make([]*display, prop.displays),
		displayImage: image.NewNRGBA(image.Rect(0, 0, prop.displaySize.X*prop.displays, prop.displaySize.Y)),
		encoder:      make([]*encoder, prop.encoders),
		key:          make([]*key, prop.keys),
	}
	d.model = prop.model(d)

	for i := range d.display {
		d.display[i] = newDisplay(d, i)
	}
	for i := range d.encoder {
		d.encoder[i] = newEncoder(d, i)
	}
	for i := range d.key {
		d.key[i] = newButton(d, i)
	}

	return d
}

type Device struct {
	model
	prop         Properties
	info         hid.DeviceInfo
	dev          *hid.Device
	mu           sync.RWMutex
	display      []*display
	displayImage *image.NRGBA
	encoder      []*encoder
	key          []*key
	event        map[benjamin.EventType][]benjamin.EventHandler
}

func (d *Device) Manufacturer() string { return d.info.Manufacturer }
func (d *Device) Product() string      { return d.info.Product }
func (d *Device) Serial() string       { return d.info.Serial }

func (d *Device) Open() error {
	var err error
	if d.dev == nil {
		d.dev, err = d.info.Open()
	}
	return err
}

func (d *Device) Close() error {
	return d.dev.Close()
}

func (d *Device) Display(index int) benjamin.Display {
	if index < 0 || index >= d.prop.displays {
		return nil
	}
	return d.display[index]
}

func (d *Device) Displays() int              { return d.prop.displays }
func (d *Device) DisplayLayout() image.Point { return d.prop.displayLayout }
func (d *Device) Encoders() int              { return d.prop.encoders }
func (d *Device) ButtonAt(p image.Point) benjamin.Button {
	return d.Button(p.Y*d.prop.keyLayout.X + p.X)
}
func (d *Device) Buttons() int              { return d.prop.keys }
func (d *Device) ButtonLayout() image.Point { return d.prop.keyLayout }

func (d *Device) Encoder(index int) benjamin.Encoder {
	if index < 0 || index >= d.prop.encoders {
		return nil
	}
	return d.encoder[index]
}

func (d *Device) Button(index int) benjamin.Button {
	//log.Printf("streamdeck: key %d requested", index)
	if index < 0 || index >= d.prop.keys {
		return nil
	}
	return d.key[index]
}

func (d *Device) Events() <-chan benjamin.Event {
	c := make(chan benjamin.Event, 16)

	go func(c chan<- benjamin.Event) {
		defer close(c)

		p := make([]byte, 64)
		for {
			n, err := d.dev.Read(p)
			if err != nil {
				c <- benjamin.NewError(d, err)
				return
			}
			//log.Printf("read %d:\n%s", n, hex.Dump(p[:n]))
			d.Handle(p[:n], c)
		}
	}(c)

	return c
}

func (d *Device) Clear() error {
	for _, d := range d.display {
		if err := d.SetImage(image.Black); err != nil {
			return err
		}
	}
	for _, k := range d.key {
		if err := k.SetImage(image.Black); err != nil {
			return err
		}
	}
	return nil
}

func (d *Device) sendFeatureReport(p []byte) error {
	if d.dev == nil {
		return io.ErrClosedPipe
	}

	d.mu.RLock()
	_, err := d.dev.SendFeatureReport(p)
	d.mu.RUnlock()
	return err
}

type model interface {
	Reset() error
	SetBrightness(float64) error
	Handle(p []byte, c chan<- benjamin.Event)
	SetButtonImage(keyIndex int, imageData []byte) error
	SetDisplayImage(imageData []byte) error
}

type baseModel struct {
	*Device
	reset               func(*Device) error
	setBrightness       func(*Device, float64) error
	imagePageHeader     func(pageIndex, keyIndex, dataSize int, isLast bool) []byte
	imagePageHeaderSize int
	imagePageSize       int
}

func (m *baseModel) Reset() error {
	return m.reset(m.Device)
}

func (m *baseModel) SetBrightness(v float64) error {
	return m.setBrightness(m.Device, v)
}

func (m *baseModel) Handle(p []byte, c chan<- benjamin.Event) {
	switch p[1] {
	case 0x00: // key
		m.handleButton(p[1:], c)
	}
}

func (m *baseModel) handleButton(p []byte, c chan<- benjamin.Event) {
	state := p[m.prop.keyDataOffset:]
	for i := 0; i < len(state) && i < m.prop.keys; i++ {
		var (
			press = state[i] != 0
			index = m.prop.keyTranslate(i)
			key   = m.key[index]
		)
		if key.state != state[index] {
			key.state = state[i]
			if press {
				key.press = time.Now()
				c <- benjamin.NewButtonPress(m, key)
			} else {
				c <- benjamin.NewButtonRelease(m, key, time.Since(key.press))
			}
		}
	}
}

func (m *baseModel) SetButtonImage(index int, imageBytes []byte) error {
	var (
		data = imageData{
			Data:     imageBytes,
			PageSize: m.prop.imagePageSize - m.prop.imagePageHeaderSize,
		}
		buf       = make([]byte, m.prop.imagePageSize)
		header, b []byte
		last      bool
		err       error
	)
	m.mu.Lock()
	defer m.mu.Unlock()
	for page := 0; !last; page++ {
		b, last = data.Page(page)
		header = m.imagePageHeader(page, index, len(b), last)
		copy(buf, header)
		copy(buf[len(header):], b)
		//log.Printf("streamdeck: key %d image page %d of %d", index, page, len(buf))
		if _, err = m.dev.Write(buf); err != nil {
			return fmt.Errorf("streamdeck: image transfer to key %d failed: %w", index, err)
		}
	}
	return nil
}

func (m *baseModel) SetDisplayImage(imageBytes []byte) error {
	const (
		displayPageSize       = 1024
		displayPageHeaderSize = 16
	)
	var (
		data = imageData{
			Data:     imageBytes,
			PageSize: displayPageSize - displayPageHeaderSize,
		}
		buf       = make([]byte, displayPageSize)
		header, b []byte
		last      bool
		err       error
	)
	m.mu.Lock()
	defer m.mu.Unlock()
	for page := 0; !last; page++ {
		b, last = data.Page(page)
		header = m.displayPageHeader(page, m.displayImage.Rect, len(b), last)
		copy(buf, header)
		copy(buf[len(header):], b)
		//log.Printf("streamdeck: display image page %d of %d", page, len(buf))
		if _, err = m.dev.Write(buf); err != nil {
			return fmt.Errorf("streamdeck: image transfer to display failed: %w", err)
		}
	}
	return nil
}

func (m *baseModel) displayPageHeader(pageIndex int, area image.Rectangle, dataSize int, isLast bool) []byte {
	var last byte
	if isLast {
		last = 0x01
	}
	var (
		x = uint16(area.Min.X)
		y = uint16(area.Min.Y)
		w = uint16(area.Dx())
		h = uint16(area.Dy())
	)
	return []byte{
		0x02, 0x0c,
		byte(x), byte(x >> 8),
		byte(y), byte(y >> 8),
		byte(w), byte(w >> 8),
		byte(h), byte(h >> 8),
		last,
		byte(pageIndex),
		byte(pageIndex >> 8),
		byte(dataSize),
		byte(dataSize >> 8),
		0x00,
	}
}

func translateLTR() func(int) int      { return func(i int) int { return i } }
func translateRTL(o int) func(int) int { return func(i int) int { return o - i - 1 } }
