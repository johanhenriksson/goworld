package builtin

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/physics"
)

func init() {
	editor.Register(&physics.Box{}, NewBoxEditor)
}

type BoxEditor struct {
	object.Object
	target *physics.Box

	Shape physics.Shape
	Body  *physics.RigidBody
}

func NewBoxEditor(ctx *editor.Context, box *physics.Box) *BoxEditor {
	body := physics.NewRigidBody("Collider", 0)
	body.Shape = physics.NewBox(box.Size())

	return object.New("BoxEditor", &BoxEditor{
		target: box,
		Body:   body,
		Shape:  body.Shape,
	})
}

func (e *BoxEditor) Actions() []editor.Action {
	return nil
}
