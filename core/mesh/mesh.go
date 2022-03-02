package mesh

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/deferred"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	object.Component
	deferred.Drawable
	deferred.ShadowDrawable
	engine.ForwardDrawable
	engine.LineDrawable

	Mesh() vertex.Mesh
	SetMesh(vertex.Mesh)
	Material() material.T
}

// mesh base
type mesh struct {
	object.Component

	mat  material.T
	data vertex.Mesh
	vao  vertex.Array
	mode DrawMode
}

// New creates a new mesh component
func New(mat material.T, mode DrawMode) T {
	return NewPrimitiveMesh(vertex.Triangles, mat, mode)
}

// NewLines creates a new line mesh component
func NewLines() T {
	material := assets.GetMaterialShared("lines")
	return NewPrimitiveMesh(vertex.Lines, material, Lines)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive vertex.Primitive, mat material.T, mode DrawMode) *mesh {
	m := &mesh{
		Component: object.NewComponent(),
		mode:      mode,
		mat:       mat,
		vao:       gl_vertex_array.New(primitive),
	}
	return m
}

func (m mesh) Name() string {
	return "Mesh"
}

func (m mesh) Mesh() vertex.Mesh {
	return m.data
}

func (m *mesh) SetMesh(data vertex.Mesh) {
	ptrs := data.Pointers()
	ptrs.Bind(m.mat)
	m.vao.SetPointers(ptrs)
	m.vao.SetIndexSize(data.IndexSize())
	m.vao.SetElements(data.Elements())
	m.vao.Buffer("vertex", data.VertexData())
	m.vao.Buffer("index", data.IndexData())

	m.data = data
}

func (m mesh) Material() material.T {
	return m.mat
}

func (m *mesh) DrawDeferred(args render.Args) error {
	if m.mode != Deferred {
		return nil
	}

	// the actual draw call belongs to the renderer
	// this should be extracted

	if err := m.mat.Use(); err != nil {
		return fmt.Errorf("failed to assign material %s in mesh %s: %w", m.mat.Name(), m.Name(), err)
	}

	// set up uniforms
	m.mat.Mat4("model", args.Transform)
	m.mat.Mat4("view", args.View)
	m.mat.Mat4("projection", args.Projection)
	m.mat.Mat4("mvp", args.MVP)
	m.mat.Vec3("eye", args.Position)

	return m.vao.Draw()
}

func (m *mesh) DrawForward(args render.Args) error {
	if m.mode != Forward {
		return nil
	}

	if err := m.mat.Use(); err != nil {
		return fmt.Errorf("failed to assign material %s in mesh %s: %w", m.mat.Name(), m.Name(), err)
	}

	// set up uniforms
	m.mat.Mat4("model", args.Transform)
	m.mat.Mat4("view", args.View)
	m.mat.Mat4("projection", args.Projection)
	m.mat.Mat4("mvp", args.MVP)

	return m.vao.Draw()
}

func (m *mesh) DrawLines(args render.Args) error {
	if m.mode != Lines {
		return nil
	}

	if err := m.mat.Use(); err != nil {
		return fmt.Errorf("failed to assign material %s in mesh %s: %w", m.mat.Name(), m.Name(), err)
	}

	m.mat.Mat4("mvp", args.MVP)

	return m.vao.Draw()
}

func (m *mesh) DrawShadow(args render.Args) error {
	if m.mode == Lines {
		// lines cant cast shadows
		return nil
	}

	return m.vao.Draw()
}
