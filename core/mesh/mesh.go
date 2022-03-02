package mesh

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	object.Component

	Mesh() vertex.Mesh
	SetMesh(vertex.Mesh)
	Material() material.T
	Mode() DrawMode
	CastShadows() bool
	Vao() vertex.Array
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

func (m *mesh) Vao() vertex.Array {
	// todo: remove this
	return m.vao
}

func (m mesh) Material() material.T {
	return m.mat
}

func (m mesh) CastShadows() bool {
	return m.mode != Lines
}

func (m mesh) Mode() DrawMode {
	return m.mode
}
