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
	Pass     render.Pass
	Material *render.Material

	name string
	vao  *render.VertexArray
}

// NewMesh creates a new mesh object
func NewMesh(name string, material *render.Material) *Mesh {
	return NewPrimitiveMesh(name, render.Triangles, render.Geometry, material)
}

// NewLineMesh creates a new mesh for drawing lines
func NewLineMesh(name string) *Mesh {
	material := assets.GetMaterialShared("lines")
	return NewPrimitiveMesh(name, render.Lines, render.Line, material)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(name string, primitive render.GLPrimitive, pass render.Pass, material *render.Material) *Mesh {
	m := &Mesh{
		Transform: Identity(),
		Pass:      pass,
		name:      name,
		Material:  material,
		vao:       render.CreateVertexArray(primitive),
	}
	return m
}

// Returns the name of the mesh
func (m *Mesh) Name() string {
	return m.name
}

// Buffer mesh data to GPU memory
// func (m *Mesh) Buffer(name string, data render.VertexData) error {
// 	m.vao.Buffer(name, data)
// 	for _, buffer := range m.material.Buffers {
// 		m.material.SetupBufferPointers(buffer)
// 	}
// 	return nil
// }

func (m *Mesh) SetIndexType(t render.GLType) {
	// get rid of this later
	m.vao.SetIndexType(t)
}

func (m *Mesh) Collect(pass DrawPass, args DrawArgs) {
	if m.Pass == pass.Type() && pass.Visible(m, args) {
		pass.Queue(m, args.Apply(m.Transform))
	}
}

func (m *Mesh) DrawDeferred(args DrawArgs) {
	m.Material.Use()
	shader := m.Material.Shader // UsePass(render.Geometry)

	// set up uniforms
	shader.Mat4("model", &args.Transform)
	shader.Mat4("view", &args.View)
	shader.Mat4("projection", &args.Projection)
	shader.Mat4("mvp", &args.MVP)
	shader.Vec3("eye", &args.Position)

	m.vao.Draw()
}

func (m *Mesh) DrawForward(args DrawArgs) {
	m.Material.Use()
	shader := m.Material.Shader // UsePass(render.Geometry)

	// set up uniforms
	shader.Mat4("model", &args.Transform)
	shader.Mat4("view", &args.View)
	shader.Mat4("projection", &args.Projection)
	shader.Mat4("mvp", &args.MVP)

	m.vao.Draw()
}

func (m *Mesh) DrawLines(args DrawArgs) {
	m.Material.Use()
	m.Material.Mat4("mvp", &args.MVP)

	m.vao.Draw()
}

func (m Mesh) Buffer(data interface{}) {
	pointers := m.Material.VertexPointers(data)

	// compatibility hack
	if len(pointers) == 0 {
		fmt.Println("error buffering mesh", m.Name, "- no pointers")
	} else {
		m.vao.BufferTo(pointers, data)
	}
}
