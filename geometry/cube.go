package geometry

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Cube mesh, textured
type Cube struct {
	mesh.T
	Size float32
}

// NewCube creates a new textured cube mesh with a given size
func NewCube(size float32) *Cube {
	cube := &Cube{
		T:    mesh.New(mesh.Deferred),
		Size: size,
	}
	cube.generate()
	return cube
}

func (c *Cube) generate() {
	s := c.Size / 2
	vertices := []vertex.T{
		// XP
		{P: vec3.New(s, -s, s), N: vec3.UnitX, T: vec2.New(1, 1)},  // 0
		{P: vec3.New(s, -s, -s), N: vec3.UnitX, T: vec2.New(1, 0)}, // 1
		{P: vec3.New(s, s, -s), N: vec3.UnitX, T: vec2.New(0, 0)},  // 2
		{P: vec3.New(s, -s, s), N: vec3.UnitX, T: vec2.New(1, 1)},  // 1
		{P: vec3.New(s, s, -s), N: vec3.UnitX, T: vec2.New(0, 0)},  // 2
		{P: vec3.New(s, s, s), N: vec3.UnitX, T: vec2.New(0, 1)},   // 3

		// XN
		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, T: vec2.New(0, 1)},  // 4
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, T: vec2.New(1, 0)},  // 5
		{P: vec3.New(-s, -s, -s), N: vec3.UnitXN, T: vec2.New(0, 0)}, // 6
		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, T: vec2.New(0, 1)},  // 4
		{P: vec3.New(-s, s, s), N: vec3.UnitXN, T: vec2.New(1, 1)},   // 7
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, T: vec2.New(1, 0)},  // 5

		// YP
		{P: vec3.New(-s, s, -s), N: vec3.UnitY, T: vec2.New(0, 0)}, // 8
		{P: vec3.New(-s, s, s), N: vec3.UnitY, T: vec2.New(0, 1)},  // 9
		{P: vec3.New(s, s, -s), N: vec3.UnitY, T: vec2.New(1, 0)},  // 10
		{P: vec3.New(s, s, -s), N: vec3.UnitY, T: vec2.New(1, 0)},  // 10
		{P: vec3.New(-s, s, s), N: vec3.UnitY, T: vec2.New(0, 1)},  // 9
		{P: vec3.New(s, s, s), N: vec3.UnitY, T: vec2.New(1, 1)},   // 11

		// YN
		{P: vec3.New(-s, -s, -s), N: vec3.UnitYN, T: vec2.New(0, 0)}, // 12
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, T: vec2.New(1, 0)},  // 13
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, T: vec2.New(0, 1)},  // 14
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, T: vec2.New(1, 0)},  // 13
		{P: vec3.New(s, -s, s), N: vec3.UnitYN, T: vec2.New(1, 1)},   // 15
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, T: vec2.New(0, 1)},  // 14

		// ZP
		{P: vec3.New(-s, -s, s), N: vec3.UnitZ, T: vec2.New(1, 0)}, // 16
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, T: vec2.New(0, 0)},  // 17
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, T: vec2.New(1, 1)},  // 18
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, T: vec2.New(0, 0)},  // 17
		{P: vec3.New(s, s, s), N: vec3.UnitZ, T: vec2.New(0, 1)},   // 19
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, T: vec2.New(1, 1)},  // 18

		// ZN
		{P: vec3.New(-s, -s, -s), N: vec3.UnitZN, T: vec2.New(0, 0)}, // 20
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, T: vec2.New(0, 1)},  // 21
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, T: vec2.New(1, 0)},  // 22
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, T: vec2.New(1, 0)},  // 22
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, T: vec2.New(0, 1)},  // 21
		{P: vec3.New(s, s, -s), N: vec3.UnitZN, T: vec2.New(1, 1)},   // 23
	}

	indices := []uint8{
		0, 1, 2,
		1, 2, 3,

		4, 5, 6,
		4, 7, 5,

		8, 9, 10,
		10, 9, 11,

		12, 13, 14,
		13, 15, 14,

		16, 17, 18,
		17, 19, 18,

		20, 21, 22,
		22, 21, 23,
	}

	mesh := vertex.NewTriangles("cube", vertices, indices)
	c.SetMesh(mesh)
}
