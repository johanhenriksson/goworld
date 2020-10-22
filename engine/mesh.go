package engine

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]*render.VertexBuffer

// Mesh base
type Mesh struct {
	*Transform
	Pass render.Pass

	name     string
	material *render.Material
	vao      *render.VertexArray
}

// NewMesh creates a new mesh object
func NewMesh(name string, material *render.Material) *Mesh {
	return NewPrimitiveMesh(name, render.Triangles, render.Geometry, material)
}

// NewLineMesh creates a new mesh for drawing lines
func NewLineMesh(name string) *Mesh {
	material := assets.GetMaterialShared("lines")
	return NewPrimitiveMesh(name, render.Lines, render.LinePass, material)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(name string, primitive render.GLPrimitive, pass render.Pass, material *render.Material) *Mesh {
	m := &Mesh{
		Transform: Identity(),
		Pass:      pass,
		name:      name,
		material:  material,
		vao:       render.CreateVertexArray(primitive),
	}
	for _, buffer := range m.material.Buffers {
		m.vao.AddBuffer(buffer)
		m.material.SetupBufferPointers(buffer)
	}
	return m
}

// Returns the name of the mesh
func (m *Mesh) Name() string {
	return m.name
}

// Buffer mesh data to GPU memory
func (m *Mesh) Buffer(name string, data render.VertexData) error {
	return m.vao.Buffer(name, data)
}

// AddIndex adds an index buffer of the given type.
func (m *Mesh) AddIndex(datatype render.GLType) {
	m.vao.AddIndexBuffer(datatype)
}

func (m *Mesh) Collect(pass DrawPass, args DrawArgs) {
	if m.Pass == pass.Type() && pass.Visible(m, args) {
		pass.Queue(m, args.Apply(m.Transform))
	}
}

func (m *Mesh) DrawDeferred(args DrawArgs) {
	m.material.Use()
	shader := m.material.Shader // UsePass(render.Geometry)

	// set up uniforms
	shader.Mat4("model", &args.Transform)
	shader.Mat4("view", &args.View)
	shader.Mat4("projection", &args.Projection)
	shader.Mat4("mvp", &args.MVP)
	shader.Vec3("eye", &args.Position)

	m.vao.Draw()
}

func (m *Mesh) DrawForward(args DrawArgs) {
	m.material.Use()
	shader := m.material.Shader // UsePass(render.Geometry)

	// set up uniforms
	shader.Mat4("model", &args.Transform)
	shader.Mat4("view", &args.View)
	shader.Mat4("projection", &args.Projection)
	shader.Mat4("mvp", &args.MVP)

	m.vao.Draw()
}

func (m *Mesh) DrawLines(args DrawArgs) {
	m.material.Use()
	m.material.Mat4("mvp", &args.MVP)

	m.vao.Draw()
}
