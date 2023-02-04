package plane

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Plane is a colored, one segment, one-sided 3D plane
type T struct {
	mesh.T
	Args
}

type Args struct {
	Size  float32
	Color color.T
	Mat   *material.Def
}

func New(args Args) *T {
	plane := object.New(&T{
		T:    mesh.New(mesh.Forward, args.Mat),
		Args: args,
	})
	plane.generate()
	return plane
}

func (p *T) generate() {
	s := p.Size / 2
	y := float32(0.001)
	c := p.Color.Vec4()

	vertices := []vertex.C{
		{P: vec3.New(-s, y, -s), N: vec3.UnitY, C: c}, // o1
		{P: vec3.New(s, y, -s), N: vec3.UnitY, C: c},  // x1
		{P: vec3.New(-s, y, s), N: vec3.UnitY, C: c},  // z1
		{P: vec3.New(s, y, s), N: vec3.UnitY, C: c},   // d1

		{P: vec3.New(-s, -y, -s), N: vec3.UnitYN, C: c}, // o2
		{P: vec3.New(s, -y, -s), N: vec3.UnitYN, C: c},  // x2
		{P: vec3.New(-s, -y, s), N: vec3.UnitYN, C: c},  // z2
		{P: vec3.New(s, -y, s), N: vec3.UnitYN, C: c},   // d2
	}

	indices := []uint16{
		0, 2, 1, 1, 2, 3,
		5, 6, 4, 7, 6, 5,
	}

	key := object.Key("plane", p)
	mesh := vertex.NewTriangles(key, vertices, indices)
	p.SetMesh(mesh)
}
