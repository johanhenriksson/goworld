package plane

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Plane struct {
	object.Object
	*Mesh
}

func NewObject(args Args) *Plane {
	return object.New("Plane", &Plane{
		Mesh: New(args),
	})
}

// Plane is a single segment, two-sided 3D plane
type Mesh struct {
	*mesh.Static
	Args
}

type Args struct {
	Size float32
	Mat  *material.Def
}

func New(args Args) *Mesh {
	if args.Mat == nil {
		args.Mat = material.StandardForward()
	}
	plane := object.NewComponent(&Mesh{
		Static: mesh.New(args.Mat),
		Args:   args,
	})
	plane.generate()
	return plane
}

func (p *Mesh) generate() {
	s := p.Size / 2
	y := float32(0.001)

	vertices := []vertex.T{
		{P: vec3.New(-s, y, -s), N: vec3.UnitY, T: vec2.New(0, 1)}, // o1
		{P: vec3.New(s, y, -s), N: vec3.UnitY, T: vec2.New(1, 1)},  // x1
		{P: vec3.New(-s, y, s), N: vec3.UnitY, T: vec2.New(0, 0)},  // z1
		{P: vec3.New(s, y, s), N: vec3.UnitY, T: vec2.New(1, 0)},   // d1

		{P: vec3.New(-s, -y, -s), N: vec3.UnitYN, T: vec2.New(0, 0)}, // o2
		{P: vec3.New(s, -y, -s), N: vec3.UnitYN, T: vec2.New(0, 0)},  // x2
		{P: vec3.New(-s, -y, s), N: vec3.UnitYN, T: vec2.New(0, 0)},  // z2
		{P: vec3.New(s, -y, s), N: vec3.UnitYN, T: vec2.New(0, 0)},   // d2
	}

	indices := []uint16{
		0, 2, 1, 1, 2, 3, // top
		5, 6, 4, 7, 6, 5, // bottom
	}

	key := object.Key("plane", p)
	mesh := vertex.NewTriangles(key, vertices, indices)
	p.VertexData.Set(mesh)
}
