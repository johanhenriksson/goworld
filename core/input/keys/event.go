package keys

import "fmt"

type Event interface {
	Code() Code
	Action() Action
	Character() rune
	Modifier() Modifier

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

func (e event) Code() Code         { return e.code }
func (e event) Character() rune    { return e.char }
func (e event) Modifier() Modifier { return e.mods }
func (e event) Action() Action     { return e.action }
func (e event) Handled() bool      { return e.handled }

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
