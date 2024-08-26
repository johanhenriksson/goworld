package plane

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	Register[*Mesh](Type{
		Name: "Plane",
		Path: []string{"Geometry"},
		Create: func(ctx Pool) (Component, error) {
			return New(ctx, Args{
				Mat:  material.StandardDeferred(),
				Size: vec2.New(1, 1),
			}), nil
		},
	})
}

type Plane struct {
	Object
	Mesh *Mesh
	Size Property[float32]
}

func New(pool Pool, args Args) *Plane {
	return NewObject(pool, "Plane", &Plane{
		Mesh: NewMesh(pool, args),
	})
}

// Plane is a single segment, two-sided 3D plane
type Mesh struct {
	*mesh.Static
	Size Property[vec2.T]

	data vertex.MutableMesh[vertex.T, uint16]
}

type Args struct {
	Size vec2.T
	Mat  *material.Def
}

func NewMesh(pool Pool, args Args) *Mesh {
	if args.Mat == nil {
		args.Mat = material.StandardForward()
	}
	p := NewComponent(pool, &Mesh{
		Static: mesh.New(pool, args.Mat),
		Size:   NewProperty[vec2.T](args.Size),
	})
	p.data = vertex.NewTriangles[vertex.T, uint16](Key("plane", p), nil, nil)
	p.Size.OnChange.Subscribe(func(f vec2.T) { p.refresh() })
	p.refresh()
	return p
}

func (p *Mesh) refresh() {
	s := p.Size.Get().Scaled(0.5)
	y := float32(0.001)

	uv := p.Size.Get().Scaled(1.0 / 8)
	vertices := []vertex.T{
		{P: vec3.New(-s.X, y, -s.Y), N: vec3.UnitY, T: vec2.New(0, uv.Y)},   // o1
		{P: vec3.New(s.X, y, -s.Y), N: vec3.UnitY, T: vec2.New(uv.X, uv.Y)}, // x1
		{P: vec3.New(-s.X, y, s.Y), N: vec3.UnitY, T: vec2.New(0, 0)},       // z1
		{P: vec3.New(s.X, y, s.Y), N: vec3.UnitY, T: vec2.New(uv.X, 0)},     // d1

		{P: vec3.New(-s.X, -y, -s.Y), N: vec3.UnitYN, T: vec2.New(0, uv.Y)},   // o2
		{P: vec3.New(s.X, -y, -s.Y), N: vec3.UnitYN, T: vec2.New(uv.X, uv.Y)}, // x2
		{P: vec3.New(-s.X, -y, s.Y), N: vec3.UnitYN, T: vec2.New(0, 0)},       // z2
		{P: vec3.New(s.X, -y, s.Y), N: vec3.UnitYN, T: vec2.New(uv.X, 0)},     // d2
	}

	indices := []uint16{
		0, 2, 1, 1, 2, 3, // top
		5, 6, 4, 7, 6, 5, // bottom
	}

	p.data.Update(vertices, indices)
	p.VertexData.Set(p.data)
}
