package streamdeck

import "github.com/tehmaze/benjamin/deck"

const (
	streamDeckXLV2ProductID = 0x008f
	streamDeckXLV2Cols      = 8
	streamDeckXLV2Rows      = 4
	streamDeckXLV2Pixels    = 96
)

func init() {
	deck.RegisterUSB(elgatoVendorID, streamDeckXLV2ProductID, deviceType{
		productID:            streamDeckXLV2ProductID,
		name:                 "Stream Deck XL V2",
		cols:                 streamDeckXLV2Cols,
		rows:                 streamDeckXLV2Rows,
		keys:                 streamDeckXLV2Cols * streamDeckXLV2Rows,
		pixels:               streamDeckXLV2Pixels,
		margin:               82,
		dpi:                  166,
		padding:              16,
		featureReportSize:    32,
		firmwareOffset:       6,
		keyStateOffset:       3,
		imagePageSize:        1024,
		imagePageHeaderSize:  8,
		imagePageHeader:      streamDeckRev2PageHeader,
		toImageFormat:        toJPEG(streamDeckXLV2Pixels),
		commandFirmware:      streamDeckRev2Firmware,
		commandReset:         streamDeckRev2Reset,
		commandSetBrightness: streamDeckRev2SetBrightness,
	}.Driver)
}
