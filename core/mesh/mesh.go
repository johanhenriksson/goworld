package mesh

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex_array"
	"github.com/johanhenriksson/goworld/render/vertex_buffer"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]vertex_buffer.T

type T interface {
	object.Component

	DrawForward(render.Args)
	DrawDeferred(render.Args)
	DrawLines(render.Args)

	SetIndexType(t types.Type)
	Buffer(data interface{})
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

func (m *mesh) SetIndexType(t types.Type) {
	// get rid of this later
	m.vao.SetIndexType(t)
}

func (m mesh) Name() string {
	return "Mesh"
}

func (m *mesh) DrawDeferred(args render.Args) {
	if m.mode != Deferred {
		return
	}

	m.mat.Use()

	// set up uniforms
	m.mat.Mat4("model", args.Transform)
	m.mat.Mat4("view", args.View)
	m.mat.Mat4("projection", args.Projection)
	m.mat.Mat4("mvp", args.MVP)
	m.mat.Vec3("eye", args.Position)

	m.vao.Draw()
}

func (m *mesh) DrawForward(args render.Args) {
	if m.mode != Forward {
		return
	}

	m.mat.Use()

	// set up uniforms
	m.mat.Mat4("model", args.Transform)
	m.mat.Mat4("view", args.View)
	m.mat.Mat4("projection", args.Projection)
	m.mat.Mat4("mvp", args.MVP)

	m.vao.Draw()
}

func (m *mesh) DrawLines(args render.Args) {
	if m.mode != Lines {
		return
	}

	m.mat.Use()
	m.mat.Mat4("mvp", args.MVP)

	m.vao.Draw()
}

func (m *mesh) Buffer(data interface{}) {
	pointers := m.mat.VertexPointers(data)

	// compatibility hack
	// ... but for what? this never seems to happen
	// more like a sanity check
	if len(pointers) == 0 {
		panic(fmt.Errorf("error buffering mesh %s - no pointers", m.Name()))
	}

	m.vao.BufferTo(pointers, data)
}
