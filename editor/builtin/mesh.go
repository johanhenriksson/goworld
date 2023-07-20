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
	return object.New("MeshEditor", &MeshEditor{
		target: mesh,
		shape:  physics.NewMesh(),
	})
}

func (e *MeshEditor) Bounds() physics.Shape {
	return e.shape
}

func (e *MeshEditor) Actions() []editor.Action {
	return nil
}

func (e *MeshEditor) Update(scene object.Component, dt float32) {
	// todo: replace with an event
	mesh := e.target.Vertices()
	if mesh != e.mesh {
		e.shape.SetMeshData(mesh)
		e.mesh = mesh
	}
	e.target.Update(scene, dt)
}
