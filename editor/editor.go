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

// EditorUpdater is an optional interface that can be implemented by
// components to receive editor updates.
type EditorUpdater interface {
	object.Component
	EditorUpdate(scene object.Component, dt float32)
}
