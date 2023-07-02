package editor

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
)

type ObjectEditor struct {
	object.G

	Bounds collider.T
	Custom T

	target object.T
}

type DefaultEditor struct {
	object.G
	GUI gui.Fragment
}

func (d *DefaultEditor) Actions() []Action { return nil }

func NewObjectEditor(target object.T, bounds collider.Box, editor T) *ObjectEditor {
	var boundsCol collider.T
	if editor != nil {
		boundsCol = collider.NewBox(bounds)
	} else {
		// instantiate default object inspector
		editor = object.Group("DefaultEditor", &DefaultEditor{
			GUI: gui.NewFragment(gui.FragmentArgs{
				Slot:     "sidebar:content",
				Position: gui.FragmentLast,
				Render: func() node.T {
					return Inspector(target, nil)
				},
			}),
		})
	}
	editor.SetActive(false)

	return object.Group("ObjectEditor", &ObjectEditor{
		G:      object.Ghost(target),
		target: target,

		Custom: editor,

		// the bounds collider must be held outside the editor object, so that it is not
		// disabled along with the editor on deselection. this prevents re-selection
		Bounds: boundsCol,
	})
}

var _ Selectable = &ObjectEditor{}

func (e *ObjectEditor) Update(scene object.T, dt float32) {
	e.G.Update(scene, dt)
	e.Custom.Update(scene, dt)
}

func (e *ObjectEditor) Select(ev mouse.Event, collider collider.T) {
	e.Custom.SetActive(true)
}

func (e *ObjectEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	e.Custom.SetActive(false)
	return true
}

func (e *ObjectEditor) Target() object.T {
	return e.target
}

func (e *ObjectEditor) Actions() []Action {
	actions := []Action{
		{
			Name: "Move",
			Key:  keys.G,
			Callback: func(mgr ToolManager) {
				mgr.MoveTool(e.target)
			},
		},
	}
	actions = append(actions, e.Custom.Actions()...)
	return actions
}
