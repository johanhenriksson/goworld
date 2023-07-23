package builtin

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/physics"
)

func init() {
	editor.Register(&physics.Box{}, NewBoxEditor)
}

type BoxEditor struct {
	object.Object
	target *physics.Box

	Shape *physics.Box
	Body  *physics.RigidBody

	GUI gui.Fragment
}

func NewBoxEditor(ctx *editor.Context, box *physics.Box) *BoxEditor {
	editor := object.New("PhysicsBoxEditor", &BoxEditor{
		Object: object.Ghost("Ghost:"+box.Name(), box.Transform()),
		target: box,

		Shape: physics.NewBox(box.Extents.Get()),
		Body:  physics.NewRigidBody(0),

		GUI: editor.SidebarFragment(gui.FragmentLast, func() node.T {
			return editor.Inspector(
				box,
				propedit.Vec3Field("extents", "Extents", propedit.Vec3Props{
					Value:    box.Extents.Get(),
					OnChange: box.Extents.Set,
				}),
			)
		}),
	})

	box.OnChange().Subscribe(func(s physics.Shape) {
		// manually apply the local scaling factor
		editor.Shape.Extents.Set(box.Extents.Get()) // .Mul(box.Transform().Scale()))
	})

	return editor
}

func (e *BoxEditor) Target() object.Component { return e.target }

func (e *BoxEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
}

func (e *BoxEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.GUI)
	return true
}

func (e *BoxEditor) Actions() []editor.Action {
	return nil
}
