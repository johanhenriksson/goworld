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
	editor.Register(&mesh.Static{}, NewMeshEditor)
}

type MeshEditor struct {
	object.Object
	mesh   vertex.Mesh
	target *mesh.Static

	Shape *physics.Mesh
	Body  *physics.RigidBody
	Mesh  *lines.Mesh
	GUI   gui.Fragment
}

func NewMeshEditor(ctx *editor.Context, mesh *mesh.Static) *MeshEditor {
	editor := object.New("MeshEditor", &MeshEditor{
		Object: object.Ghost(mesh.Name(), mesh.Transform()),
		target: mesh,

		Shape: physics.NewMesh(),
		Body:  physics.NewRigidBody(0),
		Mesh:  lines.NewMesh(mesh.VertexData.Get(), color.White),

		GUI: editor.SidebarFragment(gui.FragmentLast, func() node.T {
			return editor.Inspector(
				mesh,
			)
		}),
	})

	// grab reference to mesh shape & subscribe to changes
	editor.Shape.Mesh.Set(mesh.VertexData.Get())
	mesh.VertexData.OnChange.Subscribe(func(m vertex.Mesh) {
		editor.Shape.Mesh.Set(m)
		editor.Mesh.Source.Set(m)
	})

	return editor
}

func (e *MeshEditor) Target() object.Component { return e.target }

func (e *MeshEditor) Select(ev mouse.Event) {
	object.Enable(e.GUI)
	object.Enable(e.Mesh)
}

func (e *MeshEditor) Deselect(ev mouse.Event) bool {
	// todo: check with editor if we can deselect?
	object.Disable(e.GUI)
	object.Disable(e.Mesh)
	return true
}

func (e *MeshEditor) Actions() []editor.Action {
	return nil
}
