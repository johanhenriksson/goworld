package mesh

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/object"
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

	//
	// used for querying:
	//

	CastShadows() bool

	//
	// used for rendering:
	//

	Material() *material.Def
	MaterialID() material.ID

	// Texture returns the texture attached to the given slot
	Texture(texture.Slot) assets.Texture

	// Mesh returns the mesh data
	Mesh() assets.Mesh
}

// mesh base
type Static struct {
	object.Component

	matId material.ID

	// bounding radius
	center vec3.T
	radius float32

	CastsShadow object.Property[bool]
	Mat         object.Property[*material.Def]
	Textures    object.Dict[texture.Slot, assets.Texture]
	VertexData  object.Property[assets.Mesh]
}

var _ Mesh = (*Static)(nil)

// NewMesh creates a new mesh
func New(pool object.Pool, mat *material.Def) *Static {
	return object.NewComponent(pool, &Static{
		Mat:         object.NewProperty(mat),
		CastsShadow: object.NewProperty(true),
		Textures:    object.NewDict[texture.Slot, assets.Texture](),
		VertexData:  object.NewProperty[assets.Mesh](nil),

		matId: material.Hash(mat),
	})
}

func (m *Static) Name() string {
	return "Mesh"
}

func (m *Static) Mesh() assets.Mesh {
	return m.VertexData.Get()
}

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
	if ref := m.VertexData.Get(); ref != nil {
		mat := m.Mat.Get()
		if mat == nil {
			return false
		}
		return m.CastsShadow.Get() &&
			mat.Primitive == vertex.Triangles &&
			!mat.Transparent
	}
	return false
}

func (m *Static) Material() *material.Def {
	return m.Mat.Get()
}

func (m *Static) MaterialID() material.ID {
	return m.matId
}
