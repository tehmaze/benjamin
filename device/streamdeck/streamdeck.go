package streamdeck

import (
	"fmt"
	"image"
	"time"

	"github.com/karalabe/hid"
	log "github.com/sirupsen/logrus"

	"github.com/tehmaze/benjamin/device"
)

const (
	elgatoVendorID = 0x0fd9
)

type deviceType struct {
	productID            uint16
	name                 string
	cols, rows           int
	keys                 int
	pixels               int
	dpi                  int
	padding              int
	featureReportSize    int
	firmwareOffset       int
	keyStateOffset       int
	translateKey         func(index, cols uint8) uint8
	imagePageSize        int
	imagePageHeaderSize  int
	toImageFormat        func(image.Image) ([]byte, error)
	imagePageHeader      func(pageIndex, keyIndex, payloadLength int, lastPage bool) []byte
	commandFirmware      []byte
	commandReset         []byte
	commandSetBrightness []byte
}

func (t deviceType) Driver() device.USBDriver {
	return func(info hid.DeviceInfo) device.Device {
		return New(info, t)
	}
}

type StreamDeck struct {
	deviceType
	dev        *hid.Device
	info       hid.DeviceInfo
	key        []*key
	keyState   []byte
	keyPress   []time.Time
	keyTrigger []time.Time
}

func New(info hid.DeviceInfo, t deviceType) *StreamDeck {
	d := &StreamDeck{
		deviceType: t,
		info:       info,
		key:        make([]*key, t.keys),
		keyState:   make([]byte, t.keys),
		keyPress:   make([]time.Time, t.keys),
		keyTrigger: make([]time.Time, t.keys),
	}
	for i := range d.key {
		d.key[i] = newKey(d, i%t.cols, i/t.cols)
	}
	return d
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

func (d *StreamDeck) Key(i int) device.Key {
	if i < 0 || i >= len(d.key) {
		return nil
	}
	return d.key[i]
}

func (d *StreamDeck) KeySize() image.Point {
	return image.Pt(d.pixels, d.pixels)
}

func (d *StreamDeck) Open() (err error) {
	if d.dev, err = d.info.Open(); err == nil {
		log.WithFields(log.Fields{
			"name":         d.name,
			"vendor_id":    d.info.VendorID,
			"product_id":   d.info.ProductID,
			"serial":       d.info.Serial,
			"manufacturer": d.info.Manufacturer,
			"product":      d.info.Product,
		}).Debug("stream deck opened")
	}
	return
}

func (d *StreamDeck) Close() (err error) {
	if d.dev != nil {
		err = d.dev.Close()
	}
	return
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

func (d *StreamDeck) Events() <-chan device.Event {
	var (
		c = make(chan device.Event)
		b = make([]byte, d.keyStateOffset+d.keys)
	)
	go func(c chan<- device.Event) {
		defer close(c)

		// Trigger button presses
		/*
			go func(c chan<- device.Event) {
				for now := range time.Tick(time.Second / 10) {
					for i, state := range d.buttonState {
						if state == 0 {
							continue
						}
						var (
							first    = d.buttonPress[i]
							firstAgo = now.Sub(first)
							last     = d.buttonTrigger[i]
							lastAgo  = now.Sub(last)
						)
						for _, s := range ButtonTriggerSchedule {
							if firstAgo < s.After {
								continue
							}
							if lastAgo < s.Trigger {
								continue
							}
							d.buttonTrigger[i] = now
							c <- ButtonEvent{
								Button:   d.button[i],
								Pressed:  true,
								Duration: firstAgo,
							}
						}
					}
				}
			}(c)
		*/

		// Read events from device
		for {
			copy(d.keyState, b[d.keyStateOffset:])
			if _, err := d.dev.Read(b); err != nil {
				close(c)
				return
			}

			for i := d.keyStateOffset; i < len(b); i++ {
				j := uint8(i - d.keyStateOffset)
				if d.translateKey != nil {
					j = d.translateKey(j, uint8(d.cols))
				}
				if b[i] != d.keyState[j] {
					var (
						duration time.Duration
						pressed  = b[i] == 1
					)
					if pressed {
						// Press action immediately triggers a press
						d.keyPress[j] = time.Now()
						d.keyTrigger[j] = d.keyPress[j]
					} else {
						duration = time.Since(d.keyPress[j])
					}

					log.WithFields(log.Fields{
						"index":    j,
						"pressed":  pressed,
						"duration": duration,
					}).Debug("stream deck button")

					var t device.EventType
					if pressed {
						t = device.KeyPressed
					} else {
						t = device.KeyReleased
					}
					c <- device.Event{
						Type:     t,
						Pos:      image.Pt(int(j)%d.cols, int(j)/d.cols),
						Duration: duration,
					}
				}
			}
		}
	}(c)
	return c
}

func (d *StreamDeck) SetBrightness(v uint8) error {
	if v > 100 {
		v = 100
	}
	return d.sendFeatureReport(append(d.commandSetBrightness, v))
}

/*
func (d *StreamDeck) SetColor(c color.Color) error {
	i := image.NewUniform(c)
	for _, b := range d.button {
		if err := b.SetImage(i); err != nil {
			return err
		}
	}
	return nil
}

func (d *StreamDeck) SetImage(i image.Image) error {
	var (
		r = i.Bounds()
		w = r.Dx() / d.cols
		h = r.Dy() / d.rows
		o = image.NewRGBA(image.Rect(0, 0, w, h))
		b int
	)
	for y := 0; y < d.rows; y++ {
		for x := 0; x < d.cols; x++ {
			draw.Draw(o, o.Bounds(), i, image.Pt(x*w, y*h), draw.Src)
			if err := d.key[b].SetImage(o); err != nil {
				return err
			}
			b++
		}
	}
	return nil
}
*/

func (d *StreamDeck) getFeatureReport(p []byte) ([]byte, error) {
	b := make([]byte, d.featureReportSize)
	copy(b, p)
	if _, err := d.dev.GetFeatureReport(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (d *StreamDeck) sendFeatureReport(p []byte) error {
	b := make([]byte, d.featureReportSize)
	copy(b, p)
	_, err := d.dev.SendFeatureReport(b)
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
	device.RegisterUSB(elgatoVendorID, streamDeckProductID, deviceType{
		productID:            streamDeckProductID,
		name:                 "Stream Deck",
		cols:                 5,
		rows:                 3,
		keys:                 15,
		pixels:               72,
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
	}.Driver())
}
