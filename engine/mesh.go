package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]*render.VertexBuffer

// Mesh base
type Mesh struct {
	*ComponentBase

	material *render.Material
	vao      *render.VertexArray
	vbos     MeshBufferMap
}

// NewMesh creates a new mesh object
func NewMesh(material *render.Material) *Mesh {
	m := &Mesh{
		material: material,
		vao:      render.CreateVertexArray(),
		vbos:     MeshBufferMap{},
	}
	for _, buffer := range m.material.Buffers {
		m.addBuffer(buffer)
	}
	return m
}

// addBuffer adds a named buffer to the mesh VAO.
func (m *Mesh) addBuffer(name string) *render.VertexBuffer {
	// create new vbo
	vbo := render.CreateVertexBuffer()

	// set up vertex array pointers for this buffer
	m.vao.Bind()
	vbo.Bind()
	m.material.SetupBufferPointers(name)

	// store reference & return vbo object
	m.vbos[name] = vbo
	return vbo
}

// Buffer mesh data to the GPU
func (m *Mesh) Buffer(name string, data render.VertexData) error {
	vbo, exists := m.vbos[name]
	if !exists {
		return fmt.Errorf("Unknown VBO: %s", name)
	}
	m.vao.Bind()
	m.vao.Length = int32(data.Elements())
	return vbo.Buffer(data)
}

// Update the mesh.
func (m *Mesh) Update(dt float32) {}

// Draw the mesh.
func (m *Mesh) Draw(args render.DrawArgs) {
	if args.Pass != render.GeometryPass && args.Pass != render.LightPass {
		return
	}

	m.material.Use()

	// set up uniforms
	m.material.Mat4f("model", &args.Transform)
	m.material.Mat4f("view", &args.View)
	m.material.Mat4f("projection", &args.Projection)

	m.vao.DrawElements()
}
