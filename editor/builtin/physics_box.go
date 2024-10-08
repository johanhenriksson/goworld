package builtin

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

func init() {
	editor.RegisterEditor(&physics.Box{}, NewBoxEditor)
}

type BoxEditor struct {
	Object
	target *physics.Box

	Shape *physics.Box
	Body  *physics.RigidBody
	Mesh  *lines.Box

	GUI gui.Fragment
}

func NewBoxEditor(ctx *editor.Context, box *physics.Box) *BoxEditor {
	editor := NewObject(ctx.Objects, "PhysicsBoxEditor", &BoxEditor{
		Object: Ghost(ctx.Objects, box.Name(), box.Transform()),
		target: box,

		Shape: physics.NewBox(ctx.Objects, box.Extents.Get()),
		Body:  physics.NewRigidBody(ctx.Objects, 0),
		Mesh: lines.NewBox(ctx.Objects, lines.BoxArgs{
			Extents: box.Extents.Get(),
			Color:   color.Green,
		}),

		GUI: editor.PropertyEditorFragment(ctx.Objects, gui.FragmentLast, func() node.T {
			return editor.Inspector(
				box,
				propedit.Vec3Field("extents", "Extents", propedit.Vec3Props{
					Value:    box.Extents.Get(),
					OnChange: box.Extents.Set,
				}),
			)
		}),
	})

	// keep properties in sync
	box.Extents.OnChange.Subscribe(func(vec3.T) {
		editor.Shape.Extents.Set(box.Extents.Get())
		editor.Mesh.Extents.Set(box.Extents.Get())
	})

	return editor
}

func (e *BoxEditor) Target() Component { return e.target }

func (e *BoxEditor) Select(ev mouse.Event) {
	Enable(e.GUI)
	Enable(e.Mesh)
}

func (e *BoxEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	Disable(e.GUI)
	Disable(e.Mesh)
	return true
}

func (e *BoxEditor) Actions() []editor.Action {
	return nil
}
