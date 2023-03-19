package geometry

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// ColorCube is a vertex colored cube mesh
type ColorCube struct {
	mesh.T
	Size  float32
	Color color.T
}

// NewColorCube creates a vertex colored cube mesh with a given size
func NewColorCube(mat *material.Def, color color.T, size float32) *ColorCube {
	cube := object.New(&ColorCube{
		T:     mesh.New(mesh.Forward, mat),
		Size:  size,
		Color: color,
	})
	cube.generate()
	return cube
}

func (c *ColorCube) generate() {
	s := c.Size / 2
	co := c.Color.Vec4()
	vertices := []vertex.C{
		{P: vec3.New(s, -s, s), N: vec3.UnitX, C: co},  // 0
		{P: vec3.New(s, -s, -s), N: vec3.UnitX, C: co}, // 1
		{P: vec3.New(s, s, -s), N: vec3.UnitX, C: co},  // 2
		{P: vec3.New(s, -s, s), N: vec3.UnitX, C: co},  // 0
		{P: vec3.New(s, s, -s), N: vec3.UnitX, C: co},  // 2
		{P: vec3.New(s, s, s), N: vec3.UnitX, C: co},   // 3

		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, C: co},  // 4
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, C: co},  // 5
		{P: vec3.New(-s, -s, -s), N: vec3.UnitXN, C: co}, // 6
		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, C: co},  // 4
		{P: vec3.New(-s, s, s), N: vec3.UnitXN, C: co},   // 7
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, C: co},  // 5

		{P: vec3.New(-s, s, -s), N: vec3.UnitY, C: co}, // 8
		{P: vec3.New(-s, s, s), N: vec3.UnitY, C: co},  // 9
		{P: vec3.New(s, s, -s), N: vec3.UnitY, C: co},  // 10
		{P: vec3.New(s, s, -s), N: vec3.UnitY, C: co},  // 10
		{P: vec3.New(-s, s, s), N: vec3.UnitY, C: co},  // 9
		{P: vec3.New(s, s, s), N: vec3.UnitY, C: co},   // 11

		{P: vec3.New(-s, -s, -s), N: vec3.UnitYN, C: co}, // 12
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, C: co},  // 13
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, C: co},  // 14
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, C: co},  // 13
		{P: vec3.New(s, -s, s), N: vec3.UnitYN, C: co},   // 15
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, C: co},  // 14

		{P: vec3.New(-s, -s, s), N: vec3.UnitZ, C: co}, // 16
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, C: co},  // 17
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, C: co},  // 18
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, C: co},  // 17
		{P: vec3.New(s, s, s), N: vec3.UnitZ, C: co},   // 19
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, C: co},  // 18

		{P: vec3.New(-s, -s, -s), N: vec3.UnitZN, C: co}, // 20
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, C: co},  // 21
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, C: co},  // 22
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, C: co},  // 22
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, C: co},  // 21
		{P: vec3.New(s, s, -s), N: vec3.UnitZN, C: co},   // 23
	}

	indices := []uint16{}

	key := object.Key("colorcube", c)
	mesh := vertex.NewTriangles(key, vertices, indices)
	c.SetMesh(mesh)
}
