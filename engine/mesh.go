package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/render"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]*render.VertexBuffer

// Mesh base
type Mesh struct {
	*object.Link
	Pass     render.Pass
	Material *render.Material

	vao *render.VertexArray
}

// NewMesh creates a new mesh object
func NewMesh(material *render.Material) *Mesh {
	return NewPrimitiveMesh(render.Triangles, render.Geometry, material)
}

// NewLineMesh creates a new mesh for drawing lines
func NewLineMesh() *Mesh {
	material := assets.GetMaterialShared("lines")
	return NewPrimitiveMesh(render.Lines, render.Line, material)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive render.GLPrimitive, pass render.Pass, material *render.Material) *Mesh {
	m := &Mesh{
		Link:     object.NewLink(nil),
		Pass:     pass,
		Material: material,
		vao:      render.CreateVertexArray(primitive),
	}
	return m
}

func (m *Mesh) SetIndexType(t render.GLType) {
	// get rid of this later
	m.vao.SetIndexType(t)
}

func (m *Mesh) DrawDeferred(args DrawArgs) {
	if m.Pass != render.Geometry {
		return
	}

	m.Material.Use()
	shader := m.Material.Shader

	// set up uniforms
	shader.Mat4("model", &args.Transform)
	shader.Mat4("view", &args.View)
	shader.Mat4("projection", &args.Projection)
	shader.Mat4("mvp", &args.MVP)
	shader.Vec3("eye", &args.Position)

	m.vao.Draw()
}

func (m *Mesh) DrawForward(args DrawArgs) {
	if m.Pass != render.Forward {
		return
	}

	m.Material.Use()
	shader := m.Material.Shader

	// set up uniforms
	shader.Mat4("model", &args.Transform)
	shader.Mat4("view", &args.View)
	shader.Mat4("projection", &args.Projection)
	shader.Mat4("mvp", &args.MVP)

	m.vao.Draw()
}

func (m *Mesh) DrawLines(args DrawArgs) {
	if m.Pass != render.Line {
		return
	}

	m.Material.Use()
	m.Material.Mat4("mvp", &args.MVP)

	m.vao.Draw()
}

func (m Mesh) Buffer(data interface{}) {
	pointers := m.Material.VertexPointers(data)

	// compatibility hack
	if len(pointers) == 0 {
		panic(fmt.Errorf("error buffering mesh %s - no pointers", m.String()))
	} else {
		m.vao.BufferTo(pointers, data)
	}
}
