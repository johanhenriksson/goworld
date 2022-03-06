package cache

import (
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type MeshBackend interface {
	Instantiate(mesh vertex.Mesh, mat material.T) GpuMesh
	Update(bmesh GpuMesh, mesh vertex.Mesh)
	Delete(bmesh GpuMesh)
	Destroy()
}

// glmeshes is a mesh buffer backend for OpenGL
type glmeshes struct{}

func (m *glmeshes) Instantiate(mesh vertex.Mesh, mat material.T) GpuMesh {
	vao := gl_vertex_array.New(mesh.Primitive())

	ptrs := mesh.Pointers()
	ptrs.Bind(mat)
	vao.SetPointers(ptrs)

	return vao
}

func (m *glmeshes) Update(bmesh GpuMesh, mesh vertex.Mesh) {
	vao := bmesh.(vertex.Array)

	vao.SetIndexSize(mesh.IndexSize())
	vao.SetElements(mesh.Elements())
	vao.Buffer("vertex", mesh.VertexData())
	vao.Buffer("index", mesh.IndexData())
}

func (m *glmeshes) Delete(bmesh GpuMesh) {
	vao := bmesh.(vertex.Array)
	vao.Delete()
}

func (m *glmeshes) Destroy() {
}
