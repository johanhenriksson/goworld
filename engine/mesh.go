package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]*render.VertexBuffer

// Mesh base
type Mesh struct {
	*Object
	Pass render.DrawPass

	material *render.Material
	vao      *render.VertexArray
}

// NewMesh creates a new mesh object
func NewMesh(parent *Object, material *render.Material) *Mesh {
	m := &Mesh{
		Object:   parent,
		Pass:     render.GeometryPass,
		material: material,
		vao:      render.CreateVertexArray(render.Triangles),
	}
	for _, buffer := range m.material.Buffers {
		m.vao.AddBuffer(buffer)
		m.material.SetupBufferPointers(buffer)
	}
	parent.Attach(m)
	return m
}

// Buffer mesh data to the GPU
func (m *Mesh) Buffer(name string, data render.VertexData) error {
	return m.vao.Buffer(name, data)
}

func (m *Mesh) AddIndex(datatype render.GLType) {
	m.vao.AddIndexBuffer(datatype)
}

// Update the mesh.
func (m *Mesh) Update(dt float32) {}

// Draw the mesh.
func (m *Mesh) Draw(args render.DrawArgs) {
	if m.Pass == render.GeometryPass && args.Pass != render.GeometryPass && args.Pass != render.LightPass {
		return
	}
	if m.Pass == render.ForwardPass && args.Pass != render.ForwardPass {
		return
	}

	m.material.Use()

	// set up uniforms
	m.material.Mat4("model", &args.Transform)
	m.material.Mat4("view", &args.View)
	m.material.Mat4("projection", &args.Projection)
	m.material.Mat4("mvp", &args.MVP)
	m.material.Vec3("eye", &args.Position)

	m.vao.Draw()
}
