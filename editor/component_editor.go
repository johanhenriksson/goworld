package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/physics"
)

type ComponentEditor struct {
	object.Object
	target object.Component
	GUI    gui.Fragment
}

func NewComponentEditor(target object.Component) *ComponentEditor {
	return object.New("ComponentEditor", &ComponentEditor{
		target: target,
		GUI: InspectorGUI(
			target,
		),
	})
}

func (d *ComponentEditor) Actions() []Action     { return nil }
func (d *ComponentEditor) Bounds() physics.Shape { return nil }
