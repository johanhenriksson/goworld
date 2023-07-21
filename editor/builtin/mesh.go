package builtin

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	var kind mesh.Mesh
	editor.Register(kind, NewMeshEditor)
}

type MeshEditor struct {
	object.Object
	target mesh.Mesh
	shape  *physics.Mesh
	mesh   vertex.Mesh
}

func NewMeshEditor(ctx *editor.Context, mesh mesh.Mesh) *MeshEditor {
	editor := object.New("MeshEditor", &MeshEditor{
		target: mesh,
		shape:  physics.NewMesh(),
	})

	// propagate mesh updates to the editor collider shape
	mesh.Mesh().OnChange.Subscribe(editor, func(m vertex.Mesh) {
		editor.shape.Mesh.Set(m)
	})

	return editor
}

func (e *MeshEditor) Bounds() physics.Shape {
	return e.shape
}

func (e *MeshEditor) Actions() []editor.Action {
	return nil
}
