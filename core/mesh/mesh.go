package mesh

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	object.Component

	Mesh() vertex.Mesh
	SetMesh(vertex.Mesh)
	Mode() DrawMode
	CastShadows() bool
}

// mesh base
type mesh struct {
	object.Component

	data vertex.Mesh
	mode DrawMode
}

// New creates a new mesh component
func New(mode DrawMode) T {
	return NewPrimitiveMesh(vertex.Triangles, mode)
}

// NewLines creates a new line mesh component
func NewLines() T {
	return NewPrimitiveMesh(vertex.Lines, Lines)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive vertex.Primitive, mode DrawMode) *mesh {
	m := &mesh{
		Component: object.NewComponent(),
		mode:      mode,
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
	m.data = data
}

func (m mesh) CastShadows() bool {
	return m.mode != Lines
}

func (m mesh) Mode() DrawMode {
	return m.mode
}
