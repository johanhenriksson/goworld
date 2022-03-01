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
	"github.com/johanhenriksson/goworld/render/vertex_array"
	"github.com/johanhenriksson/goworld/render/vertex_buffer"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]vertex_buffer.T

type T interface {
	object.Component
	deferred.Drawable
	deferred.ShadowDrawable
	engine.ForwardDrawable
	engine.LineDrawable

	SetPointers(vertex.Pointers)
	BufferRaw(name string, elements int, data []byte)
	Buffer(name string, data interface{})
}

// mesh base
type mesh struct {
	object.Component

	mat  material.T
	vao  vertex_array.T
	mode DrawMode
}

// New creates a new mesh component
func New(mat material.T, mode DrawMode) T {
	return NewPrimitiveMesh(render.Triangles, mat, mode)
}

// NewLines creates a new line mesh component
func NewLines() T {
	material := assets.GetMaterialShared("lines")
	return NewPrimitiveMesh(render.Lines, material, Lines)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive render.Primitive, mat material.T, mode DrawMode) *mesh {
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

func (m mesh) Material() material.T {
	return m.mat
}

func (m *mesh) DrawDeferred(args render.Args) error {
	if m.mode != Deferred {
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

func (m *mesh) Buffer(name string, data interface{}) {
	m.vao.Buffer(name, data)

	if name == "vertex" {
		pointers := vertex.ParsePointers(data)
		pointers.Bind(m.mat)
		m.vao.SetPointers(pointers)
	}
}

//
// these functions do not really belong here
//

func (m *mesh) BufferRaw(name string, elements int, data []byte) {
	m.vao.BufferRaw(name, elements, data)
}

func (m *mesh) SetPointers(ptrs vertex.Pointers) {
	ptrs.Bind(m.mat)
	m.vao.SetPointers(ptrs)
}
