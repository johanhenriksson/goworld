package editor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/physics"
)

// EditorGhost connects a scene object with an editor implementation.
// Its primary purpose is to toggle the editor implementation, and provide
// editor object picking by inserting a collider in the editor physics world.
// EditorGhost mirrors the transformation of its target object.
type EditorGhost struct {
	object.Object
	target object.Component

	Editor T
	Body   *physics.RigidBody
}

func NewEditorGhost(target object.Component, editor T) *EditorGhost {
	if editor == nil {
		if obj, isObject := target.(object.Object); isObject {
			// use default object editor
			editor = NewObjectEditor(obj)
		} else {
			// use default component editor
			editor = NewComponentEditor(target)
		}
	}
	object.Disable(editor)

	// collider for object picking
	var body *physics.RigidBody
	if editor.Bounds() != nil {
		body = physics.NewRigidBody("Collider:"+target.Name(), 0)
		object.Attach(body, editor.Bounds())
	}

	return object.New("Ghost", &EditorGhost{
		Object: object.Ghost(target),
		target: target,

		Editor: editor,

		// the bounds rigidbody must be held outside the editor object, so that it is not
		// disabled along with the editor on deselection. this prevents re-selection
		Body: body,
	})
}

func (e *EditorGhost) Name() string {
	return fmt.Sprintf("EditorGhost[%s]", e.target.Name())
}

func (e *EditorGhost) Update(scene object.Component, dt float32) {
	e.Object.Update(scene, dt)
	e.Editor.Update(scene, dt)
}

func (e *EditorGhost) Select(ev mouse.Event) {
	object.Enable(e.Editor)
}

func (e *EditorGhost) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.Editor)
	return true
}

func (e *EditorGhost) Target() object.Component {
	return e.target
}

func (e *EditorGhost) Actions() []Action {
	return e.Editor.Actions()
}
