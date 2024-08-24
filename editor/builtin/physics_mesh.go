package builtin

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	editor.RegisterEditor(&physics.Mesh{}, NewPhysicsMeshEditor)
}

type PhysicsMeshEditor struct {
	Object
	mesh vertex.Mesh

	target *physics.Mesh
	Shape  *physics.Mesh
	Body   *physics.RigidBody
	Wire   *lines.Wireframe
	GUI    gui.Fragment
}

func NewPhysicsMeshEditor(ctx *editor.Context, mesh *physics.Mesh) *PhysicsMeshEditor {
	editor := NewObject(ctx.Objects, "PhysicsMeshEditor", &PhysicsMeshEditor{
		Object: Ghost(ctx.Objects, mesh.Name(), mesh.Transform()),
		target: mesh,

		Shape: physics.NewMesh(ctx.Objects),
		Body:  physics.NewRigidBody(ctx.Objects, 0),
		Wire:  lines.NewWireframe(ctx.Objects, mesh.Mesh.Get(), color.Green),

		GUI: editor.PropertyEditorFragment(ctx.Objects, gui.FragmentLast, func() node.T {
			return editor.Inspector(
				mesh,
			)
		}),
	})

	// grab reference to mesh shape & subscribe to changes
	editor.Shape.Mesh.Set(mesh.Mesh.Get())
	mesh.Mesh.OnChange.Subscribe(func(m vertex.Mesh) {
		editor.Shape.Mesh.Set(m)
		editor.Wire.Source.Set(m)
	})

	return editor
}

func (e *PhysicsMeshEditor) Target() Component { return e.target }

func (e *PhysicsMeshEditor) Select(ev mouse.Event) {
	Enable(e.GUI)
	Enable(e.Wire)
}

func (e *PhysicsMeshEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	Disable(e.GUI)
	Disable(e.Wire)
	return true
}

func (e *PhysicsMeshEditor) Actions() []editor.Action {
	return nil
}
