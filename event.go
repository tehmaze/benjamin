package benjamin

import (
	"fmt"
	"image"
	"time"
)

// Event from user interaction.
type Event struct {
	// Type of event.
	Type EventType

	// Peripheral the event was generated with.
	Peripheral Peripheral

	// Data associated with the event.
	Data EventData
}

func (e Event) String() string {
	return fmt.Sprintf("type=%s data=%s", e.Type, e.Data)
}

type EventType int

const (
	TypeError EventType = iota
	TypeDisplayPress
	TypeDisplayLongPress
	TypeDisplaySwipe
	TypeEncoderChange
	TypeEncoderPress
	TypeEncoderRelease
	TypeKeyPress
	TypeKeyRelease
	TypeMax
)

var eventTypeName = map[EventType]string{
	TypeError:            "Error",
	TypeDisplayPress:     "DisplayPress",
	TypeDisplayLongPress: "DisplayLongPress",
	TypeDisplaySwipe:     "DisplaySwipe",
	TypeEncoderChange:    "EncoderChange",
	TypeEncoderPress:     "EncoderPress",
	TypeEncoderRelease:   "EncoderRelease",
	TypeKeyPress:         "KeyPress",
	TypeKeyRelease:       "KeyRelease",
}

func (t EventType) String() string {
	if s, ok := eventTypeName[t]; ok {
		return s
	}
	return "invalid"
}

type EventData interface {
	// Device the event was on.
	Device() Device

	// Time of event.
	Time() time.Time

	// Stringer interface
	fmt.Stringer
}

type EventHandler interface {
	Handle(Event)
}

type EventHandlerFunc func(Event)

func (f EventHandlerFunc) Handle(event Event) {
	f(event)
}

type BaseEvent struct {
	On Device
	At time.Time
}

func makeBaseEvent(device Device) BaseEvent {
	return BaseEvent{
		On: device,
		At: time.Now(),
	}
}

func (b BaseEvent) Device() Device  { return b.On }
func (b BaseEvent) Time() time.Time { return b.At }

type Error struct {
	BaseEvent
	Error error
}

func NewError(device Device, err error) Event {
	return Event{
		Type: TypeError,
		Data: Error{
			BaseEvent: makeBaseEvent(device),
			Error:     err,
		},
	}
}

func (event Error) String() string {
	return fmt.Sprintf("error:%s", event.Error)
}

type DisplayPress struct {
	BaseEvent
	Display
	Position image.Point
}

func NewDisplayPress(device Device, display Display, at image.Point) Event {
	return Event{
		Type:       TypeDisplayPress,
		Peripheral: display,
		Data: DisplayPress{
			BaseEvent: makeBaseEvent(device),
			Display:   display,
			Position:  at,
		},
	}
}

func (event DisplayPress) String() string {
	return fmt.Sprintf("display %d press: position=%s", event.Display.Index(), event.Position)
}

type DisplayLongPress struct {
	BaseEvent
	Display
	Position image.Point
}

func NewDisplayLongPress(device Device, display Display, at image.Point) Event {
	return Event{
		Type:       TypeDisplayLongPress,
		Peripheral: display,
		Data: DisplayLongPress{
			BaseEvent: makeBaseEvent(device),
			Display:   display,
			Position:  at,
		},
	}
}

func (event DisplayLongPress) String() string {
	return fmt.Sprintf("display %d long press: position=%s", event.Display.Index(), event.Position)
}

type DisplaySwipe struct {
	BaseEvent
	Display
	From, To image.Point
}

func NewDisplaySwipe(device Device, display Display, from, to image.Point) Event {
	return Event{
		Type:       TypeDisplaySwipe,
		Peripheral: display,
		Data: DisplaySwipe{
			BaseEvent: makeBaseEvent(device),
			Display:   display,
			From:      from,
			To:        to,
		},
	}
}

func (event DisplaySwipe) String() string {
	return fmt.Sprintf("display %d swipe: position=%s->%s", event.Display.Index(), event.From, event.To)
}

type EncoderChange struct {
	BaseEvent
	Encoder
	Change int
	Bits   int
}

func NewEncoderChange(device Device, encoder Encoder, change, bits int) Event {
	return Event{
		Type:       TypeEncoderChange,
		Peripheral: encoder,
		Data: EncoderChange{
			BaseEvent: makeBaseEvent(device),
			Encoder:   encoder,
			Change:    change,
			Bits:      bits,
		},
	}
}

func (event EncoderChange) String() string {
	return fmt.Sprintf("encoder %d change: %d", event.Encoder.Index(), event.Change)
}

type EncoderPress struct {
	BaseEvent
	Encoder
}

func NewEncoderPress(device Device, encoder Encoder) Event {
	return Event{
		Type:       TypeEncoderPress,
		Peripheral: encoder,
		Data: EncoderPress{
			BaseEvent: makeBaseEvent(device),
			Encoder:   encoder,
		},
	}
}

func (event EncoderPress) String() string {
	return fmt.Sprintf("encoder %d press", event.Encoder.Index())
}

type EncoderRelease struct {
	BaseEvent
	Encoder
	After time.Duration
}

func NewEncoderRelease(device Device, encoder Encoder, after time.Duration) Event {
	return Event{
		Type:       TypeEncoderRelease,
		Peripheral: encoder,
		Data: EncoderRelease{
			BaseEvent: makeBaseEvent(device),
			Encoder:   encoder,
			After:     after,
		},
	}
}

func (event EncoderRelease) String() string {
	return fmt.Sprintf("encoder %d release: after=%s", event.Encoder.Index(), event.After)
}

type KeyPress struct {
	BaseEvent
	Key
}

func NewKeyPress(device Device, key Key) Event {
	return Event{
		Type:       TypeKeyPress,
		Peripheral: key,
		Data: KeyPress{
			BaseEvent: makeBaseEvent(device),
			Key:       key,
		},
	}
}

func (event KeyPress) String() string {
	return fmt.Sprintf("key %s press", event.Key.Position())
}

type KeyRelease struct {
	BaseEvent
	After time.Duration
	Key
}

func NewKeyRelease(device Device, key Key, after time.Duration) Event {
	return Event{
		Type:       TypeKeyRelease,
		Peripheral: key,
		Data: KeyRelease{
			BaseEvent: makeBaseEvent(device),
			After:     after,
			Key:       key,
		},
	}
}

func (event KeyRelease) String() string {
	return fmt.Sprintf("key %s release: after=%s", event.Key.Position(), event.After)
}
