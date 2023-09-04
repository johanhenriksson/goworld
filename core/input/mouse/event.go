package mouse

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Event interface {
	Action() Action
	Button() Button
	Position() vec2.T
	Delta() vec2.T
	Scroll() vec2.T
	Modifier() keys.Modifier
	Project(vec2.T) Event
	Locked() bool

	Handled() bool
	Consume()
}

type event struct {
	action   Action
	button   Button
	position vec2.T
	delta    vec2.T
	scroll   vec2.T
	mods     keys.Modifier
	handled  bool
	locked   bool
}

func (e event) Action() Action          { return e.action }
func (e event) Button() Button          { return e.button }
func (e event) Position() vec2.T        { return e.position }
func (e event) Delta() vec2.T           { return e.delta }
func (e event) Scroll() vec2.T          { return e.scroll }
func (e event) Modifier() keys.Modifier { return e.mods }
func (e event) Handled() bool           { return e.handled }
func (e event) Locked() bool            { return e.locked }

func (e *event) Consume() {
	e.handled = true
}

func (e *event) Project(relativePos vec2.T) Event {
	projected := *e
	projected.position = projected.position.Sub(relativePos)
	return &projected
}

func (e event) String() string {
	switch e.action {
	case Move:
		return fmt.Sprintf("MouseEvent: Moved to %.0f,%.0f (delta %.0f,%.0f)",
			e.position.X, e.position.Y,
			e.delta.X, e.delta.Y)
	case Press:
		return fmt.Sprintf("MouseEvent: Press %s at %.0f,%.0f", e.button, e.position.X, e.position.Y)
	case Release:
		return fmt.Sprintf("MouseEvent: Release %s at %.0f,%.0f", e.button, e.position.X, e.position.Y)
	case Scroll:
		return fmt.Sprintf("MouseEvent: Scroll %.0f,%.0f", e.scroll.X, e.scroll.Y)
	}
	return "MouseEvent: Invalid"
}

func NewButtonEvent(button Button, action Action, pos vec2.T, mod keys.Modifier, locked bool) Event {
	return &event{
		action:   action,
		button:   button,
		mods:     mod,
		position: pos,
		locked:   locked,
	}
}

func NewMouseEnterEvent() Event {
	return &event{
		action: Enter,
	}
}

func NewMouseLeaveEvent() Event {
	return &event{
		action: Leave,
	}
}

func NopEvent() Event {
	return &event{action: -1}
}
