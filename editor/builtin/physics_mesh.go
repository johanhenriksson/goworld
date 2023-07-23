package builtin

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	editor.Register(&physics.Mesh{}, NewPhysicsMeshEditor)
}

type PhysicsMeshEditor struct {
	object.Object
	mesh vertex.Mesh

	target *physics.Mesh
	Shape  *physics.Mesh
	Body   *physics.RigidBody
	GUI    gui.Fragment
}

func NewPhysicsMeshEditor(ctx *editor.Context, mesh *physics.Mesh) *PhysicsMeshEditor {
	editor := object.New("PhysicsMeshEditor", &PhysicsMeshEditor{
		Object: object.Ghost(mesh.Name(), mesh.Transform()),
		target: mesh,

		Shape: physics.NewMesh(),
		Body:  physics.NewRigidBody(0),

		GUI: editor.SidebarFragment(gui.FragmentLast, func() node.T {
			return editor.Inspector(
				mesh,
			)
		}),
	})

	// grab reference to mesh shape & subscribe to changes
	editor.Shape.Mesh.Set(mesh.Mesh.Get())
	mesh.Mesh.OnChange.Subscribe(func(m vertex.Mesh) {
		editor.Shape.Mesh.Set(m)
	})

	return editor
}

func (e *PhysicsMeshEditor) Target() object.Component { return e.target }

func (e *PhysicsMeshEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
}

func (e *PhysicsMeshEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.GUI)
	return true
}

func (e *PhysicsMeshEditor) Actions() []editor.Action {
	return nil
}
