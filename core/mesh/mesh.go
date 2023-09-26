package mesh

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	object.Register[*Static](Deserialize)
}

type Mesh interface {
	object.Component

	Primitive() vertex.Primitive
	CastShadows() bool
	Material() *material.Def
	MaterialID() material.ID

	Texture(texture.Slot) texture.Ref

	// Bounding sphere used for view frustum culling
	BoundingSphere() shape.Sphere

	// returns the VertexData property
	// this is kinda ugly - but other components might need to subscribe to changes ?
	Mesh() *object.Property[vertex.Mesh]
}

// mesh base
type Static struct {
	object.Component

	primitive vertex.Primitive
	shadows   bool
	mat       *material.Def
	matId     material.ID

	textures map[texture.Slot]texture.Ref

	// bounding radius
	center vec3.T
	radius float32

	VertexData object.Property[vertex.Mesh]
}

// New creates a new mesh component
func New(mat *material.Def) *Static {
	return NewPrimitiveMesh(vertex.Triangles, mat)
}

// NewLines creates a new line mesh component
func NewLines() *Static {
	return NewPrimitiveMesh(vertex.Lines, nil)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(primitive vertex.Primitive, mat *material.Def) *Static {
	m := object.NewComponent(&Static{
		mat:       mat,
		matId:     material.Hash(mat),
		textures:  make(map[texture.Slot]texture.Ref),
		primitive: primitive,
		shadows:   true,

		VertexData: object.NewProperty[vertex.Mesh](nil),
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

func (m *Static) Primitive() vertex.Primitive         { return m.primitive }
func (m *Static) Mesh() *object.Property[vertex.Mesh] { return &m.VertexData }

func (m *Static) Texture(slot texture.Slot) texture.Ref {
	return m.textures[slot]
}

func (m *Static) SetTexture(slot texture.Slot, ref texture.Ref) {
	m.textures[slot] = ref
}

func (m *Static) CastShadows() bool {
	return m.primitive == vertex.Triangles && m.shadows && !m.mat.Transparent
}

func (m *Static) SetShadows(shadows bool) {
	m.shadows = shadows
}

func (m *Static) Material() *material.Def {
	return m.mat
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

type MeshState struct {
	object.ComponentState
	Primitive vertex.Primitive
	Material  material.Def
}

func (m *Static) Serialize(enc object.Encoder) error {
	return enc.Encode(MeshState{
		// send help
		ComponentState: object.NewComponentState(m.Component),
		Primitive:      m.primitive,
		Material:       *m.Material(),
	})
}

func Deserialize(dec object.Decoder) (object.Component, error) {
	var state MeshState
	if err := dec.Decode(&state); err != nil {
		return nil, err
	}

	obj := NewPrimitiveMesh(state.Primitive, &state.Material)
	obj.Component = state.ComponentState.New()
	return obj, nil
}
