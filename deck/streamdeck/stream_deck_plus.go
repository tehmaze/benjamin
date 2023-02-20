package streamdeck

import (
	"encoding/hex"
	"image"
	"log"

	"github.com/tehmaze/benjamin/deck"
)

const (
	streamDeckPlusProductID = 0x0084
	streamDeckPlusCols      = 4
	streamDeckPlusRows      = 2
	streamDeckPlusPixels    = 120
	streamDeckPlusLCDWidth  = 200
	streamDeckPlusLCDHeight = 100
)

func handleKeyOrEncoderInput(d *StreamDeck, b []byte, c chan<- deck.Event) {
	log.Println("handle key input:\n", hex.Dump(b))
	switch b[1] {
	case 0x00: // Key
		d.handleKeyInput(b, c)
	case 0x02: // LCD
		d.handleLCDInput(b, c)
	case 0x03: // Encoder
		d.handleEncoderInput(b, c)
	}
}

func init() {
	deck.RegisterUSB(elgatoVendorID, streamDeckPlusProductID, deviceType{
		productID:            streamDeckPlusProductID,
		name:                 "Stream Deck +",
		cols:                 streamDeckPlusCols,
		rows:                 streamDeckPlusRows,
		keys:                 streamDeckPlusCols * streamDeckPlusRows,
		pixels:               streamDeckPlusPixels,
		displays:             4,
		displaySize:          image.Pt(streamDeckPlusLCDWidth, streamDeckPlusLCDHeight),
		encoders:             4,
		margin:               24,
		dpi:                  124,
		padding:              16,
		featureReportSize:    32,
		firmwareOffset:       6,
		keyStateOffset:       4,
		imagePageSize:        1024,
		imagePageHeaderSize:  8,
		imagePageHeader:      streamDeckRev2PageHeader,
		handleInput:          handleKeyOrEncoderInput,
		toImageFormat:        toJPEGVerbatim,
		commandFirmware:      streamDeckRev2Firmware,
		commandReset:         streamDeckRev2Reset,
		commandSetBrightness: streamDeckRev2SetBrightness,
	}.Driver)
}
