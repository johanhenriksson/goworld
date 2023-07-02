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
	object.G
	target *physics.Box

	Body *physics.RigidBody
}

func NewBoxEditor(ctx *editor.Context, box *physics.Box) *BoxEditor {
	body := physics.NewRigidBody(0)
	body.Shape = box

	return object.Group("BoxEditor", &BoxEditor{
		target: box,
		Body:   body,
	})
}

func (e *BoxEditor) Actions() []editor.Action {
	return nil
}
