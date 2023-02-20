package loupedeck

const (
	loupedeckVendorID = 0x2ec2
)

// Knobs and buttons.
const (
	KnobTopLeft     = 0x01
	KnobCenterLeft  = 0x02
	KnobBottomLeft  = 0x03
	KnobTopRight    = 0x04
	KnobCenterRight = 0x05
	KnobBottomRight = 0x06
	Button0         = 0x07
	Button1         = 0x08
	Button2         = 0x09
	Button3         = 0x0a
	Button4         = 0x0b
	Button5         = 0x0c
	Button6         = 0x0d
	Button7         = 0x0e
)

// Command bytes
const (
	cmdButtonPress   = 0x00
	cmdKnobRotate    = 0x01
	cmdSetColor      = 0x02
	cmdSerial        = 0x03
	cmdReset         = 0x06
	cmdVersion       = 0x07
	cmdSetBrightness = 0x09
	cmdMCU           = 0x0d
	cmdDraw          = 0x0f
	cmdFrameBuffer   = 0x10
	cmdTouch         = 0x4d
	cmdTouchEnd      = 0x6d
)

// Response word
const (
	resConfirm            = 0x0302
	resTick               = 0x0400
	resSetBrightness      = 0x0409
	resConfirmFrameBuffer = 0x0410
	resButton             = 0x0500
	resEncoder            = 0x0501
	resReset              = 0x0506
	resDraw               = 0x050f
	resTouch              = 0x094d
	resTouchEnd           = 0x096d
	resVersion            = 0x180d
)

const (
	hapticShort  = 0x01
	hapticMedium = 0x0a
	hapticLong   = 0x0f
)

const (
	maxBrightness = 10
)
