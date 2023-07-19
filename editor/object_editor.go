package editor

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/physics"
)

type ObjectEditor struct {
	object.Object
	target object.Object
	GUI    gui.Fragment
}

func NewObjectEditor(target object.Object) *ObjectEditor {
	return object.New("ObjectEditor", &ObjectEditor{
		target: target,
		GUI: InspectorGUI(
			target,
			propedit.Transform("transform", target.Transform()),
		),
	})
}

func (e *ObjectEditor) Actions() []Action {
	return []Action{
		{
			Name: "Move",
			Key:  keys.G,
			Callback: func(mgr ToolManager) {
				mgr.MoveTool(e.target)
			},
		},
	}
}

func (e *ObjectEditor) Bounds() physics.Shape {
	return nil
}
