package builtin

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	editor.Register(&physics.Mesh{}, NewPhysicsMeshEditor)
}

type PhysicsMeshEditor struct {
	object.Object
	target *physics.Mesh
	shape  *physics.Mesh
	mesh   vertex.Mesh
}

func NewPhysicsMeshEditor(ctx *editor.Context, mesh *physics.Mesh) *PhysicsMeshEditor {
	editor := object.New("PhysicsMeshEditor", &PhysicsMeshEditor{
		target: mesh,
		shape:  physics.NewMesh(),
	})

	// grab reference to mesh shape & subscribe to changes
	editor.shape.Mesh.Set(mesh.Mesh.Get())
	mesh.Mesh.OnChange.Subscribe(editor, func(m vertex.Mesh) {
		editor.shape.Mesh.Set(m)
	})

	return editor
}

func (e *PhysicsMeshEditor) Bounds() physics.Shape {
	return e.shape
}

func (e *PhysicsMeshEditor) Actions() []editor.Action {
	return nil
}
