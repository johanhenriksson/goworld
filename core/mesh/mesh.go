package mesh

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	object.T

	Mesh() vertex.Mesh
	SetMesh(vertex.Mesh)
	Mode() DrawMode
	CastShadows() bool
	Material() *material.Def
	MaterialID() uint64
}

// mesh base
type mesh struct {
	object.T

	data  vertex.Mesh
	mode  DrawMode
	mat   *material.Def
	matId uint64
}

// New creates a new mesh component
func New(mode DrawMode, mat *material.Def) T {
	return NewPrimitiveMesh(vertex.Triangles, mode, mat)
}

// NewLines creates a new line mesh component
func NewLines(mat *material.Def) T {
	return NewPrimitiveMesh(vertex.Lines, Lines, mat)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive vertex.Primitive, mode DrawMode, mat *material.Def) *mesh {
	m := object.New(&mesh{
		mode:  mode,
		mat:   mat,
		matId: material.Hash(mat),
	})
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

func (m mesh) Material() *material.Def {
	return m.mat
}

func (m mesh) MaterialID() uint64 {
	return m.matId
}
