package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]*render.VertexBuffer

// Mesh base
type Mesh struct {
	*Transform
	Pass DrawPass

	material *render.Material
	vao      *render.VertexArray
}

// NewMesh creates a new mesh object
func NewMesh(material *render.Material) *Mesh {
	return NewPrimitiveMesh(render.Triangles, material)
}

// NewLineMesh creates a new mesh for drawing lines
func NewLineMesh() *Mesh {
	material := assets.GetMaterialCached("lines")
	return NewPrimitiveMesh(render.Lines, material)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive render.GLPrimitive, material *render.Material) *Mesh {
	m := &Mesh{
		Transform: Identity(),
		Pass:      DrawGeometry,
		material:  material,
		vao:       render.CreateVertexArray(primitive),
	}
	for _, buffer := range m.material.Buffers {
		m.vao.AddBuffer(buffer)
		m.material.SetupBufferPointers(buffer)
	}
	return m
}

// Buffer mesh data to GPU memory
func (m *Mesh) Buffer(name string, data render.VertexData) error {
	return m.vao.Buffer(name, data)
}

// AddIndex adds an index buffer of the given type.
func (m *Mesh) AddIndex(datatype render.GLType) {
	m.vao.AddIndexBuffer(datatype)
}

// Draw the mesh.
func (m *Mesh) Draw(args DrawArgs) {
	if m.Pass == DrawGeometry && args.Pass != DrawGeometry && args.Pass != DrawShadow {
		return
	}
	if m.Pass == DrawForward && args.Pass != DrawForward {
		return
	}
	if m.Pass == DrawLines {
		fmt.Println("draw line mesh!")
	}

	m.material.Use()
	args = args.Apply(m.Transform)

	// set up uniforms
	m.material.Mat4("model", &args.Transform)
	m.material.Mat4("view", &args.View)
	m.material.Mat4("projection", &args.Projection)
	m.material.Mat4("mvp", &args.MVP)
	m.material.Vec3("eye", &args.Position)

	m.vao.Draw()
}
