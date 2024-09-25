package sprite

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	object.Register[*Mesh](object.Type{})
}

// Sprite is a single segment, one-sided 3D plane
type Mesh struct {
	*mesh.Static
	Size   object.Property[vec2.T]
	Sprite object.Property[assets.Texture]

	mesh vertex.MutableMesh[vertex.Vertex, uint32]
}

var _ mesh.Mesh = &Mesh{}

type Args struct {
	Size    vec2.T
	Texture assets.Texture
}

func New(pool object.Pool, args Args) *Mesh {
	sprite := object.NewComponent(pool, &Mesh{
		Static: mesh.New(pool, Material()),
		Size:   object.NewProperty(args.Size),
		Sprite: object.NewProperty(args.Texture),
	})

	sprite.mesh = vertex.NewTriangles[vertex.Vertex, uint32](fmt.Sprintf("sprite_%.2f_%.2f", args.Size.X, args.Size.Y), nil, nil)
	sprite.generate()

	sprite.SetTexture(texture.Diffuse, args.Texture)
	sprite.Sprite.OnChange.Subscribe(func(tex assets.Texture) {
		sprite.SetTexture(texture.Diffuse, tex)
	})

	return sprite
}

func (p *Mesh) generate() {
	w, h := p.Size.Get().X, p.Size.Get().Y
	vertices := []vertex.Vertex{
		vertex.New(vec3.New(-0.5*w, -0.5*h, 0), vec3.Zero, vec2.New(0, 1), color.White),
		vertex.New(vec3.New(0.5*w, 0.5*h, 0), vec3.Zero, vec2.New(1, 0), color.White),
		vertex.New(vec3.New(-0.5*w, 0.5*h, 0), vec3.Zero, vec2.New(0, 0), color.White),
		vertex.New(vec3.New(0.5*w, -0.5*h, 0), vec3.Zero, vec2.New(1, 1), color.White),
	}
	indices := []uint32{
		0, 1, 2,
		0, 3, 1,
	}
	p.mesh.Update(vertices, indices)
	p.VertexData.Set(p.mesh)
}
