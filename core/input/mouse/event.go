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

	Handled() bool
	StopPropagation()
}

type event struct {
	action   Action
	button   Button
	position vec2.T
	delta    vec2.T
	scroll   vec2.T
	mods     keys.Modifier
	handled  bool
}

func (e event) Action() Action          { return e.action }
func (e event) Button() Button          { return e.button }
func (e event) Position() vec2.T        { return e.position }
func (e event) Delta() vec2.T           { return e.delta }
func (e event) Scroll() vec2.T          { return e.scroll }
func (e event) Modifier() keys.Modifier { return e.mods }
func (e event) Handled() bool           { return e.handled }

func (e *event) StopPropagation() {
	e.handled = true
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
	}
	return "MouseEvent: Invalid"
}
