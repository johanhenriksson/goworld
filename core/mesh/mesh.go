package mesh

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	object.Register[*Static](object.Type{
		Name: "Mesh",
		Create: func(pool object.Pool) (object.Component, error) {
			return New(pool, nil), nil
		},
	})
}

type Mesh interface {
	object.Component

	Primitive() vertex.Primitive
	CastShadows() bool
	Material() *material.Def
	MaterialID() material.ID

	Texture(texture.Slot) assets.Texture

	// Bounding sphere used for view frustum culling
	BoundingSphere() shape.Sphere

	// returns the VertexData property
	// this is kinda ugly - but other components might need to subscribe to changes ?
	Mesh() *object.Property[vertex.Mesh]
}

// mesh base
type Static struct {
	object.Component

	matId material.ID

	// bounding radius
	center vec3.T
	radius float32

	Prim        object.Property[vertex.Primitive]
	CastsShadow object.Property[bool]
	Mat         object.Property[*material.Def]
	Textures    object.Dict[texture.Slot, assets.Texture]
	VertexData  object.Property[vertex.Mesh]
}

var _ Mesh = (*Static)(nil)

// todo: this needs an argument struct

// New creates a new mesh component
func New(pool object.Pool, mat *material.Def) *Static {
	return NewPrimitiveMesh(pool, vertex.Triangles, mat)
}

// NewLines creates a new line mesh component
func NewLines(pool object.Pool) *Static {
	return NewPrimitiveMesh(pool, vertex.Lines, nil)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(pool object.Pool, primitive vertex.Primitive, mat *material.Def) *Static {
	m := object.NewComponent(pool, &Static{
		Mat:         object.NewProperty(mat),
		CastsShadow: object.NewProperty(true),
		Prim:        object.NewProperty(primitive),
		Textures:    object.NewDict[texture.Slot, assets.Texture](),
		VertexData:  object.NewProperty[vertex.Mesh](nil),

		matId: material.Hash(mat),
	})
	m.VertexData.OnChange.Subscribe(func(data vertex.Mesh) {
		// refresh bounding sphere
		min := data.Min()
		max := data.Max()
		m.center = max.Sub(min).Scaled(0.5)
		m.radius = m.center.Length()
	})
	return m
}

func (m *Static) Name() string {
	return "Mesh"
}

func (m *Static) Primitive() vertex.Primitive         { return m.Prim.Get() }
func (m *Static) Mesh() *object.Property[vertex.Mesh] { return &m.VertexData }

func (m *Static) Texture(slot texture.Slot) assets.Texture {
	t, exists := m.Textures.Get(slot)
	if !exists {
		panic(fmt.Errorf("texture slot %s does not exist", slot))
	}
	return t
}

func (m *Static) SetTexture(slot texture.Slot, ref assets.Texture) {
	m.Textures.Set(slot, ref)
}

func (m *Static) CastShadows() bool {
	return m.Prim.Get() == vertex.Triangles &&
		m.CastsShadow.Get() &&
		!m.Mat.Get().Transparent
}

func (m *Static) Material() *material.Def {
	return m.Mat.Get()
}

func (m *Static) MaterialID() material.ID {
	return m.matId
}

func (m *Static) BoundingSphere() shape.Sphere {
	return shape.Sphere{
		Center: m.Transform().WorldPosition().Add(m.center),
		Radius: m.radius,
	}
}
