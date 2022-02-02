package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
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
	mat := assets.GetMaterialShared("default")
	cube := &Cube{
		T:    mesh.New(mat, mesh.Deferred),
		Size: size,
	}
	cube.generate()
	return cube
}

func (c *Cube) generate() {
	s := c.Size / 2
	data := []vertex.T{
		// XP
		{P: vec3.New(s, -s, s), N: vec3.UnitX, T: vec2.New(1, 1)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitX, T: vec2.New(1, 0)},
		{P: vec3.New(s, s, -s), N: vec3.UnitX, T: vec2.New(0, 0)},
		{P: vec3.New(s, -s, s), N: vec3.UnitX, T: vec2.New(1, 1)},
		{P: vec3.New(s, s, -s), N: vec3.UnitX, T: vec2.New(0, 0)},
		{P: vec3.New(s, s, s), N: vec3.UnitX, T: vec2.New(0, 1)},

		// XN
		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, T: vec2.New(0, 1)},
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, T: vec2.New(1, 0)},
		{P: vec3.New(-s, -s, -s), N: vec3.UnitXN, T: vec2.New(0, 0)},
		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, T: vec2.New(0, 1)},
		{P: vec3.New(-s, s, s), N: vec3.UnitXN, T: vec2.New(1, 1)},
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, T: vec2.New(1, 0)},

		// YP
		{P: vec3.New(-s, s, -s), N: vec3.UnitY, T: vec2.New(0, 0)},
		{P: vec3.New(-s, s, s), N: vec3.UnitY, T: vec2.New(0, 1)},
		{P: vec3.New(s, s, -s), N: vec3.UnitY, T: vec2.New(1, 0)},
		{P: vec3.New(s, s, -s), N: vec3.UnitY, T: vec2.New(1, 0)},
		{P: vec3.New(-s, s, s), N: vec3.UnitY, T: vec2.New(0, 1)},
		{P: vec3.New(s, s, s), N: vec3.UnitY, T: vec2.New(1, 1)},

		// YN
		{P: vec3.New(-s, -s, -s), N: vec3.UnitYN, T: vec2.New(0, 0)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, T: vec2.New(1, 0)},
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, T: vec2.New(0, 1)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, T: vec2.New(1, 0)},
		{P: vec3.New(s, -s, s), N: vec3.UnitYN, T: vec2.New(1, 1)},
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, T: vec2.New(0, 1)},

		// ZP
		{P: vec3.New(-s, -s, s), N: vec3.UnitZ, T: vec2.New(1, 0)},
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, T: vec2.New(0, 0)},
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, T: vec2.New(1, 1)},
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, T: vec2.New(0, 0)},
		{P: vec3.New(s, s, s), N: vec3.UnitZ, T: vec2.New(0, 1)},
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, T: vec2.New(1, 1)},

		// ZN
		{P: vec3.New(-s, -s, -s), N: vec3.UnitZN, T: vec2.New(0, 0)},
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, T: vec2.New(0, 1)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, T: vec2.New(1, 0)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, T: vec2.New(1, 0)},
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, T: vec2.New(0, 1)},
		{P: vec3.New(s, s, -s), N: vec3.UnitZN, T: vec2.New(1, 1)},
	}
	c.Buffer(data)
}
