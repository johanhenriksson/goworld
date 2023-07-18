package builtin

import (
	"log"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/physics"
)

func init() {
	editor.Register(&physics.Box{}, NewBoxEditor)
}

type BoxEditor struct {
	object.Object
	target *physics.Box
	shape  physics.Shape

	GUI gui.Fragment
}

func NewBoxEditor(ctx *editor.Context, box *physics.Box) *BoxEditor {
	return object.New("BoxEditor", &BoxEditor{
		target: box,
		shape:  physics.NewBox(box.Size()),

		GUI: editor.InspectorGUI(box, nil),
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
