package builtin

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/editor/propedit"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

func init() {
	editor.Register(&physics.Sphere{}, NewSphereEditor)
}

type SphereEditor struct {
	object.Object
	target *physics.Sphere

	Shape *physics.Sphere
	Body  *physics.RigidBody
	Mesh  *lines.Sphere

	GUI gui.Fragment
}

func NewSphereEditor(ctx *editor.Context, sphere *physics.Sphere) *SphereEditor {
	editor := object.New("PhysicsSphereEditor", &SphereEditor{
		Object: object.Ghost(sphere.Name(), sphere.Transform()),
		target: sphere,

		Shape: physics.NewSphere(sphere.Radius.Get()),
		Body:  physics.NewRigidBody(0),
		Mesh: lines.NewSphere(lines.SphereArgs{
			Radius: sphere.Radius.Get(),
			Color:  color.Green,
		}),

		GUI: editor.SidebarFragment(gui.FragmentLast, func() node.T {
			return editor.Inspector(
				sphere,
				propedit.FloatField("radius", "Radius", propedit.FloatProps{
					Value:    sphere.Radius.Get(),
					OnChange: sphere.Radius.Set,
				}),
			)
		}),
	})

	// keep properties in sync
	sphere.Radius.OnChange.Subscribe(func(r float32) {
		editor.Shape.Radius.Set(r)
		editor.Mesh.Radius.Set(r)
	})

	return editor
}

func (e *SphereEditor) Target() object.Component { return e.target }

func (e *SphereEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
	object.Enable(e.Mesh)
}

func (e *SphereEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.GUI)
	object.Disable(e.Mesh)
	return true
}

func (e *SphereEditor) Actions() []editor.Action {
	return nil
}
