package editor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/physics"
)

// ObjectEditor connects a scene object with an editor implementation.
// Its primary purpose is to toggle the editor implementation, and provide
// editor object picking by inserting a collider in the editor physics world.
// ObjectEditors mirrors the transformation of its target object.
type ObjectEditor struct {
	object.Object

	Body *physics.RigidBody

	Editor T

	target object.Component
}

func NewObjectEditor(target object.Component, editor T) *ObjectEditor {
	if editor == nil {
		editor = NewDefaultEditor(target)
	}
	object.Disable(editor)

	// collider for object picking
	var body *physics.RigidBody
	if editor.Bounds() != nil {
		body = physics.NewRigidBody("Collider:"+target.Name(), 0)
		object.Attach(body, editor.Bounds())
	}

	return object.New("ObjectEditor", &ObjectEditor{
		Object: object.Ghost(target),
		target: target,

		Editor: editor,

		// the bounds rigidbody must be held outside the editor object, so that it is not
		// disabled along with the editor on deselection. this prevents re-selection

		Body: body,
	})
}

func (e *ObjectEditor) Name() string {
	_, isObject := e.target.(object.Object)
	return fmt.Sprintf("ObjectEditor[%s,%t]", e.target.Name(), isObject)
}

func (e *ObjectEditor) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)
	e.Editor.Update(scene, dt)
}

func (e *ObjectEditor) Select(ev mouse.Event) {
	object.Enable(e.Editor)
}

func (e *ObjectEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.Editor)
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
	actions = append(actions, e.Editor.Actions()...)
	return actions
}
