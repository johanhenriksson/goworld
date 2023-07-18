package builtin

import (
	"log"

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
	return object.New("MeshEditor", &PhysicsMeshEditor{
		target: mesh,
		shape:  physics.NewMesh(),
	})
}

func (e *PhysicsMeshEditor) Bounds() physics.Shape {
	return e.shape
}

func (e *PhysicsMeshEditor) OnEnable() {
	log.Println("ENABLE PHYSICS MESH EDITOR FOR", e.target.Parent().Name())
}

func (e *PhysicsMeshEditor) OnDisable() {
	log.Println("DISABLE PHYSICS MESH EDITOR FOR", e.target.Parent().Name())
}

func (e *PhysicsMeshEditor) Actions() []editor.Action {
	return nil
}

func (e *PhysicsMeshEditor) Update(scene object.Component, dt float32) {
	mesh := e.target.MeshData()
	if mesh != e.mesh {
		e.shape.SetMeshData(mesh)
		e.mesh = mesh
	}
	e.target.Update(scene, dt)
}
