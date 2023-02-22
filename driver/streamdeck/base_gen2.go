package streamdeck

import "math"

const (
	gen2ImagePageHeaderSize = 8
	gen2ImagePageSize       = 1024
)

func gen2Reset(d *Device) error {
	return d.sendFeatureReport([]byte{
		0x03, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
}

func gen2SetBrightness(d *Device, v float64) error {
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	perc := byte(math.Ceil(v * 100))
	return d.sendFeatureReport([]byte{
		0x03, 0x08, perc, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
}

func gen2ImagePageHeader(pageIndex, keyIndex, dataSize int, isLast bool) []byte {
	var last byte
	if isLast {
		last = 0x01
	}
	return []byte{
		0x02, 0x07,
		byte(keyIndex),
		last,
		byte(dataSize),
		byte(dataSize >> 8),
		byte(pageIndex),
		byte(pageIndex >> 8),
	}
}

func gen2(device *Device) model {
	return &baseModel{
		Device:              device,
		reset:               gen2Reset,
		setBrightness:       gen2SetBrightness,
		imagePageHeader:     gen2ImagePageHeader,
		imagePageHeaderSize: gen2ImagePageHeaderSize,
		imagePageSize:       gen2ImagePageSize,
	}
}
