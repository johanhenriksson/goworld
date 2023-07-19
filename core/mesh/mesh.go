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

type Mesh interface {
	object.Component

	Vertices() vertex.Mesh
	Mode() DrawMode
	CastShadows() bool
	Material() *material.Def
	MaterialID() uint64

	Texture(string) texture.Ref

	BoundingSphere() shape.Sphere
}

// mesh base
type Static struct {
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
func New(mode DrawMode, mat *material.Def) *Static {
	return NewPrimitiveMesh(vertex.Triangles, mode, mat)
}

// NewLines creates a new line mesh component
func NewLines(mat *material.Def) *Static {
	return NewPrimitiveMesh(vertex.Lines, Lines, mat)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive vertex.Primitive, mode DrawMode, mat *material.Def) *Static {
	m := object.NewComponent(&Static{
		mode:     mode,
		mat:      mat,
		matId:    material.Hash(mat),
		textures: make(map[string]texture.Ref),
	})
	return m
}

func (m *Static) Name() string {
	return "Mesh"
}

func (m *Static) Vertices() vertex.Mesh {
	return m.data
}

func (m *Static) SetVertices(data vertex.Mesh) {
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

func (m *Static) Texture(slot string) texture.Ref {
	return m.textures[slot]
}

func (m *Static) SetTexture(slot string, ref texture.Ref) {
	m.textures[slot] = ref
}

func (m *Static) CastShadows() bool {
	return m.mode != Lines
}

func (m *Static) Mode() DrawMode {
	return m.mode
}

func (m *Static) Material() *material.Def {
	return m.mat
}

func (m *Static) MaterialID() uint64 {
	return m.matId
}

func (m *Static) BoundingSphere() shape.Sphere {
	return shape.Sphere{
		Center: m.Transform().WorldPosition().Add(m.center),
		Radius: m.radius,
	}
}
