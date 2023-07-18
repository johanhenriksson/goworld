package editor

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/physics"
)

type DefaultEditor struct {
	object.Object
	GUI gui.Fragment
}

func (d *DefaultEditor) Actions() []Action     { return nil }
func (d *DefaultEditor) Bounds() physics.Shape { return nil }

func NewDefaultEditor(target object.Component) *DefaultEditor {
	return object.New("DefaultEditor", &DefaultEditor{
		GUI: InspectorGUI(target, nil),
	})
}
