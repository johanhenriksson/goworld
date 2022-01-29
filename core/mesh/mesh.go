package mesh

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/material"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]*render.VertexBuffer

type T interface {
	object.Component

	DrawForward(render.Args)
	DrawDeferred(render.Args)
	DrawLines(render.Args)

	SetIndexType(t render.GLType)
	Buffer(data interface{})
}

// mesh base
type mesh struct {
	object.Component

	mat material.T
	vao *render.VertexArray
}

// New creates a new mesh component
func New(mat material.T) T {
	return NewPrimitiveMesh(render.Triangles, mat)
}

// NewLines creates a new line mesh component
func NewLines() T {
	material := assets.GetMaterialShared("lines")
	return NewPrimitiveMesh(render.Lines, material)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive render.GLPrimitive, mat material.T) *mesh {
	m := &mesh{
		Component: object.NewComponent(),
		mat:       mat,
		vao:       render.CreateVertexArray(primitive),
	}
	return m
}

func (m *mesh) SetIndexType(t render.GLType) {
	// get rid of this later
	m.vao.SetIndexType(t)
}

func (m mesh) Name() string {
	return "Mesh"
}

func (m *mesh) DrawDeferred(args render.Args) {
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
	m.mat.Use()

	// set up uniforms
	m.mat.Mat4("model", args.Transform)
	m.mat.Mat4("view", args.View)
	m.mat.Mat4("projection", args.Projection)
	m.mat.Mat4("mvp", args.MVP)

	m.vao.Draw()
}

func (m *mesh) DrawLines(args render.Args) {
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
