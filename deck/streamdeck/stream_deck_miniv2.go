package streamdeck

import "github.com/tehmaze/benjamin/deck"

const (
	streamDeckMiniV2ProductID = 0x0090
	streamDeckMiniV2Cols      = 3
	streamDeckMiniV2Rows      = 2
	streamDeckMiniV2Pixels    = 80
)

func init() {
	deck.RegisterUSB(elgatoVendorID, streamDeckMiniV2ProductID, deviceType{
		productID:            streamDeckMiniV2ProductID,
		name:                 "Stream Deck Mini V2",
		cols:                 streamDeckMiniV2Cols,
		rows:                 streamDeckMiniV2Rows,
		keys:                 streamDeckMiniV2Cols * streamDeckMiniV2Rows,
		pixels:               streamDeckMiniV2Pixels,
		dpi:                  138,
		padding:              16,
		featureReportSize:    17,
		firmwareOffset:       5,
		keyStateOffset:       0,
		translateKey:         translateRightToLeft,
		imagePageSize:        1024,
		imagePageHeaderSize:  16,
		imagePageHeader:      streamDeckRev1PageHeader,
		toImageFormat:        toBMP,
		commandFirmware:      streamDeckRev1Firmware,
		commandReset:         streamDeckRev1Reset,
		commandSetBrightness: streamDeckRev1SetBrightness,
	}.Driver)
}
