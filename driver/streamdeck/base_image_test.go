package streamdeck

import (
	"image"
	"image/color"
	"testing"
)

func TestFlipH(t *testing.T) {
	i := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 2; x++ {
			i.Set(x, y, color.Black)
		}
		for x := 2; x < 4; x++ {
			i.Set(x, y, color.White)
		}
	}

	/*
	 * +--+--+--+--+      +--+--+--+--+
	 * |00|00|ff|ff|      |ff|ff|00|00|
	 * +--+--+--+--+      +--+--+--+--+
	 * |00|00|ff|ff|      |ff|ff|00|00|
	 * +--+--+--+--+  ->  +--+--+--+--+
	 * |00|00|ff|ff|      |ff|ff|00|00|
	 * +--+--+--+--+      +--+--+--+--+
	 * |00|00|ff|ff|      |ff|ff|00|00|
	 * +--+--+--+--+      +--+--+--+--+
	 */
	flipH(i)

	if c := i.NRGBAAt(0, 0); c.R != 0xff || c.G != 0xff || c.B != 0xff {
		t.Errorf("expected (0,0) to be white (0xff), got %#02x%02x%02x", c.R, c.G, c.B)
	}
	if c := i.NRGBAAt(2, 0); c.R != 0x00 || c.G != 0x00 || c.B != 0x00 {
		t.Errorf("expected (2,0) to be black (0x00), got %#02x%02x%02x", c.R, c.G, c.B)
	}
}

func TestFlipV(t *testing.T) {
	i := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	for x := 0; x < 4; x++ {
		for y := 0; y < 2; y++ {
			i.Set(x, y, color.Black)
		}
		for y := 2; y < 4; y++ {
			i.Set(x, y, color.White)
		}
	}

	/*
	 * +--+--+--+--+      +--+--+--+--+
	 * |00|00|00|00|      |ff|ff|ff|ff|
	 * +--+--+--+--+      +--+--+--+--+
	 * |00|00|00|00|      |ff|ff|ff|ff|
	 * +--+--+--+--+  ->  +--+--+--+--+
	 * |ff|ff|ff|ff|      |00|00|00|00|
	 * +--+--+--+--+      +--+--+--+--+
	 * |ff|ff|ff|ff|      |00|00|00|00|
	 * +--+--+--+--+      +--+--+--+--+
	 */
	flipV(i)

	if c := i.NRGBAAt(0, 0); c.R != 0xff || c.G != 0xff || c.B != 0xff {
		t.Errorf("expected (0,0) to be white (0xffffff), got %#02x%02x%02x", c.R, c.G, c.B)
	}
	if c := i.NRGBAAt(0, 2); c.R != 0x00 || c.G != 0x00 || c.B != 0x00 {
		t.Errorf("expected (0,2) to be black (0x000000), got %#02x%02x%02x", c.R, c.G, c.B)
	}
}

func TestRotate180(t *testing.T) {
	i := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		if y < 2 {
			for x := 0; x < 2; x++ {
				i.Set(x, y, color.Black)
			}
			for x := 2; x < 4; x++ {
				i.Set(x, y, color.White)
			}
		} else {
			for x := 0; x < 4; x++ {
				i.Set(x, y, color.White)
			}
		}
	}

	/*
	 * +--+--+--+--+      +--+--+--+--+
	 * |00|00|ff|ff|      |ff|ff|ff|ff|
	 * +--+--+--+--+      +--+--+--+--+
	 * |00|00|ff|ff|      |ff|ff|ff|ff|
	 * +--+--+--+--+  ->  +--+--+--+--+
	 * |ff|ff|ff|ff|      |ff|ff|00|00|
	 * +--+--+--+--+      +--+--+--+--+
	 * |ff|ff|ff|ff|      |ff|ff|00|00|
	 * +--+--+--+--+      +--+--+--+--+
	 */
	rotate180(i)

	if c := i.NRGBAAt(0, 0); c.R != 0xff || c.G != 0xff || c.B != 0xff {
		t.Errorf("expected (0,0) to be white (0xffffff), got %#02x%02x%02x", c.R, c.G, c.B)
	}
	if c := i.NRGBAAt(2, 2); c.R != 0x00 || c.G != 0x00 || c.B != 0x00 {
		t.Errorf("expected (0,2) to be black (0x000000), got %#02x%02x%02x", c.R, c.G, c.B)
	}
}
