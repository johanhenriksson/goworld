package builtin

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	editor.RegisterEditor(&mesh.Static{}, NewMeshEditor)
}

type MeshEditor struct {
	object.Object
	mesh   vertex.Mesh
	target *mesh.Static

	Shape *physics.Mesh
	Body  *physics.RigidBody
	Wire  *lines.Wireframe
	GUI   gui.Fragment
}

func NewMeshEditor(ctx *editor.Context, mesh *mesh.Static) *MeshEditor {
	editor := object.NewObject(ctx.Objects, "MeshEditor", &MeshEditor{
		Object: object.Ghost(ctx.Objects, mesh.Name(), mesh.Transform()),
		target: mesh,

		Shape: physics.NewMesh(ctx.Objects),
		Body:  physics.NewRigidBody(ctx.Objects, 0),
		Wire:  lines.NewWireframe(ctx.Objects, mesh.VertexData.Get(), color.White),

		GUI: editor.PropertyEditorFragment(ctx.Objects, gui.FragmentLast, func() node.T {
			return editor.Inspector(
				mesh,
			)
		}),
	})

	// grab reference to mesh shape & subscribe to changes
	editor.Shape.Mesh.Set(mesh.VertexData.Get())
	mesh.VertexData.OnChange.Subscribe(func(m vertex.Mesh) {
		editor.Shape.Mesh.Set(m)
		editor.Wire.Source.Set(m)
	})

	return editor
}

func (e *MeshEditor) Target() object.Component { return e.target }

func (e *MeshEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
	object.Enable(e.Wire)
}

func (e *MeshEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.GUI)
	object.Disable(e.Wire)
	return true
}

func (e *MeshEditor) Actions() []editor.Action {
	return nil
}
