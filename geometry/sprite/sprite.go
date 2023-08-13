package sprite

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Sprite struct {
	object.Object
	*Mesh
}

func NewObject(args Args) *Sprite {
	return object.New("Sprite", &Sprite{
		Mesh: New(args),
	})
}

// Sprite is a single segment, two-sided 3D plane
type Mesh struct {
	*mesh.Static
	Size   object.Property[vec2.T]
	Sprite object.Property[texture.Ref]

	mesh vertex.MutableMesh[vertex.T, uint16]
}

var _ mesh.Mesh = &Mesh{}

type Args struct {
	Size    vec2.T
	Texture texture.Ref
}

func New(args Args) *Mesh {
	sprite := object.NewComponent(&Mesh{
		Static: mesh.New(Material()),
		Size:   object.NewProperty(args.Size),
		Sprite: object.NewProperty(args.Texture),
	})

	sprite.mesh = vertex.NewTriangles[vertex.T, uint16]("sprite", nil, nil)
	sprite.generate()

	sprite.SetTexture(texture.Diffuse, args.Texture)
	sprite.Sprite.OnChange.Subscribe(func(tex texture.Ref) {
		sprite.SetTexture(texture.Diffuse, tex)
	})

	return sprite
}

func (p *Mesh) generate() {
	w, h := p.Size.Get().X, p.Size.Get().Y
	vertices := []vertex.T{
		{P: vec3.New(-0.5*w, -0.5*h, 0), T: vec2.New(0, 1)},
		{P: vec3.New(0.5*w, 0.5*h, 0), T: vec2.New(1, 0)},
		{P: vec3.New(-0.5*w, 0.5*h, 0), T: vec2.New(0, 0)},
		{P: vec3.New(0.5*w, -0.5*h, 0), T: vec2.New(1, 1)},
	}
	indices := []uint16{
		0, 1, 2,
		0, 3, 1,
	}
	p.mesh.Update(vertices, indices)
	p.VertexData.Set(p.mesh)
}
