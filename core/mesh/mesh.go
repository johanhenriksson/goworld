package mesh

import (
	"log"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Component interface {
	object.Component

	Mesh() vertex.Mesh
	SetMesh(vertex.Mesh)
	Mode() DrawMode
	CastShadows() bool
	Material() *material.Def
	MaterialID() uint64

	Texture(string) texture.Ref
	SetTexture(string, texture.Ref)

	BoundingSphere() shape.Sphere
}

// mesh base
type mesh struct {
	object.Component

	data  vertex.Mesh
	mode  DrawMode
	mat   *material.Def
	matId uint64

	textures map[string]texture.Ref

	// bounding radius
	center vec3.T
	radius float32
}

// New creates a new mesh component
func New(mode DrawMode, mat *material.Def) Component {
	return NewPrimitiveMesh(vertex.Triangles, mode, mat)
}

// NewLines creates a new line mesh component
func NewLines(mat *material.Def) Component {
	return NewPrimitiveMesh(vertex.Lines, Lines, mat)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive vertex.Primitive, mode DrawMode, mat *material.Def) *mesh {
	m := object.NewComponent(&mesh{
		mode:     mode,
		mat:      mat,
		matId:    material.Hash(mat),
		textures: make(map[string]texture.Ref),
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

	// refresh AABB
	min := data.Min()
	max := data.Max()
	m.radius = math.Max(min.Length(), max.Length())
	m.center = max.Sub(min).Scaled(0.5)

	log.Println("mesh", m, ": trigger mesh update event")

	// raise a mesh update event
	for _, handler := range object.GetAll[UpdateHandler](m) {
		handler.OnMeshUpdate(data)
	}
}

func (m *mesh) Texture(slot string) texture.Ref {
	return m.textures[slot]
}

func (m *mesh) SetTexture(slot string, ref texture.Ref) {
	m.textures[slot] = ref
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

func (m mesh) BoundingSphere() shape.Sphere {
	return shape.Sphere{
		Center: m.Transform().WorldPosition().Add(m.center),
		Radius: m.radius,
	}
}
