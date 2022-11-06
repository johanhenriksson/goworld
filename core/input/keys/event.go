package keys

import "fmt"

type Event interface {
	Code() Code
	Action() Action
	Character() rune
	Modifier(Modifier) bool

	Handled() bool
	StopPropagation()
}

type event struct {
	handled bool
	code    Code
	char    rune
	action  Action
	mods    Modifier
}

func (e event) Code() Code      { return e.code }
func (e event) Character() rune { return e.char }
func (e event) Action() Action  { return e.action }
func (e event) Handled() bool   { return e.handled }

func (e event) Modifier(mod Modifier) bool {
	return e.mods&mod == mod
}

func (e *event) StopPropagation() {
	e.handled = true
}

func (e event) String() string {
	switch e.action {
	case Press:
		return fmt.Sprintf("KeyEvent: %s %d %d", e.action, e.code, e.mods)
	case Release:
		return fmt.Sprintf("KeyEvent: %s %d %d", e.action, e.code, e.mods)
	case Repeat:
		return fmt.Sprintf("KeyEvent: %s %d %d", e.action, e.code, e.mods)
	case Char:
		return fmt.Sprintf("KeyEvent: %s %c", e.action, e.char)
	}
	return fmt.Sprintf("KeyEvent: Invalid Action %x", e.action)
}

func NewCharEvent(char rune, mods Modifier) Event {
	return &event{
		action: Char,
		char:   char,
		mods:   mods,
	}
}

func NewPressEvent(code Code, action Action, mods Modifier) Event {
	return &event{
		code:   code,
		action: action,
		mods:   mods,
	}
}
