package engine

import (
	"fmt"
	"github.com/johanhenriksson/goworld/render"
)

type IndexedMesh struct {
	*Mesh
	idx *render.VertexBuffer
}

func NewIndexedMesh(mat *render.Material) *IndexedMesh {
	m := &IndexedMesh{
		Mesh: NewMesh(mat),
		idx:  render.CreateIndexBuffer(),
	}
	m.vao.Bind()
	m.idx.Bind()
	return m
}

func (m *IndexedMesh) BufferIndices(indices render.VertexData) {
	m.vao.Bind()
	m.vao.Length = int32(indices.Elements())
	m.idx.Buffer(indices)
}

// Buffer mesh data to the GPU
func (m *IndexedMesh) Buffer(name string, data render.VertexData) error {
	vbo, exists := m.vbos[name]
	if !exists {
		return fmt.Errorf("Unknown VBO: %s", name)
	}
	m.vao.Bind()
	return vbo.Buffer(data)
}

// Draw the mesh
func (m *IndexedMesh) Draw(args render.DrawArgs) {
	m.material.Use()
	m.material.Mat4f("model", &args.Transform)
	m.material.Mat4f("view", &args.View)
	m.material.Mat4f("projection", &args.Projection)

	m.vao.DrawIndexed()
}
