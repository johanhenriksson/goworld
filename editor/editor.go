package editor

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
)

type T interface {
	object.Object
	Target() object.Component

	Select(ev mouse.Event)
	Deselect(ev mouse.Event) bool

	Actions() []Action
}
