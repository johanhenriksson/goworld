package editor

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/physics"
)

type ObjectEditor struct {
	object.Object

	Body *physics.RigidBody

	Custom T

	target object.Component
}

var _ Selectable = &ObjectEditor{}

type DefaultEditor struct {
	object.Object
	GUI gui.Fragment
}

func (d *DefaultEditor) Actions() []Action     { return nil }
func (d *DefaultEditor) Bounds() physics.Shape { return nil }

func NewObjectEditor(target object.Component, editor T) *ObjectEditor {
	if editor == nil {
		// instantiate default object inspector
		editor = object.New("DefaultEditor", &DefaultEditor{
			GUI: gui.NewFragment(gui.FragmentArgs{
				Slot:     "sidebar:content",
				Position: gui.FragmentLast,
				Render: func() node.T {
					return Inspector(target, nil)
				},
			}),
		})
	}
	object.Disable(editor)

	var body *physics.RigidBody
	if editor.Bounds() != nil {
		body = physics.NewRigidBody("Collider:"+target.Name(), 0)
		object.Attach(body, editor.Bounds())
	}

	return object.New("ObjectEditor", &ObjectEditor{
		Object: object.Ghost(target),
		target: target,

		Custom: editor,

		// the bounds collider must be held outside the editor object, so that it is not
		// disabled along with the editor on deselection. this prevents re-selection

		Body: body,
	})
}

func (e *ObjectEditor) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)
	e.Custom.Update(scene, dt)
}

func (e *ObjectEditor) Select(ev mouse.Event) {
	object.Enable(e.Custom)
}

func (e *ObjectEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.Custom)
	return true
}

func (e *ObjectEditor) Target() object.Component {
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
