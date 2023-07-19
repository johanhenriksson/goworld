package builtin

import (
	"log"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
)

func init() {
	editor.Register(&physics.Box{}, NewBoxEditor)
}

type BoxEditor struct {
	object.Object
	target *physics.Box
	shape  *physics.Box

	GUI gui.Fragment
}

func NewBoxEditor(ctx *editor.Context, box *physics.Box) *BoxEditor {
	shape := physics.NewBox(box.Extents.Get())

	return object.New("BoxEditor", &BoxEditor{
		target: box,
		shape:  shape,

		GUI: editor.InspectorGUI(
			box,
			propedit.Vec3Field("extents", "Extents", propedit.Vec3Props{
				Value: box.Extents.Get(),
				OnChange: func(t vec3.T) {
					box.Extents.Set(t)
					shape.Extents.Set(t)
				},
			}),
		),
	})
}

func (e *BoxEditor) Bounds() physics.Shape {
	return e.shape
}

func (e *BoxEditor) OnEnable() {
	log.Println("ENABLE PHYSICS BOX EDITOR FOR", e.target.Parent().Name())
}

func (e *BoxEditor) OnDisable() {
	log.Println("DISABLE PHYSICS BOX EDITOR FOR", e.target.Parent().Name())
}

func (e *BoxEditor) Actions() []editor.Action {
	return nil
}
