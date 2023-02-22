package streamdeck

import (
	"math"
)

const (
	gen1ImagePageHeaderSize = 16
	gen1ImagePageSize       = 7819
)

func gen1Reset(d *Device) error {
	return d.sendFeatureReport([]byte{
		0x0b,
		0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
}

func gen1SetBrightness(d *Device, v float64) error {
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	perc := byte(math.Ceil(v * 100))
	return d.sendFeatureReport([]byte{
		0x05,
		0x55, 0xaa, 0xd1, 0x01, perc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
}

func gen1ImagePageHeader(pageIndex, keyIndex, dataSize int, isLast bool) []byte {
	var last byte
	if isLast {
		last = 0x01
	}
	return []byte{
		0x02, 0x01,
		byte(pageIndex + 1), 0x00,
		last,
		byte(keyIndex + 1),
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

func gen1(device *Device) model {
	return &baseModel{
		Device:              device,
		reset:               gen1Reset,
		setBrightness:       gen1SetBrightness,
		imagePageHeader:     gen1ImagePageHeader,
		imagePageHeaderSize: gen1ImagePageHeaderSize,
		imagePageSize:       gen1ImagePageSize,
	}
}
