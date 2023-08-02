package cube

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Object struct {
	object.Object
	*Mesh
}

func NewObject(args Args) *Object {
	return object.New("Cube", &Object{
		Mesh: New(args),
	})
}

// Mesh is a vertex colored cube mesh
type Mesh struct {
	*mesh.Static
	Args
}

type Args struct {
	Mat  *material.Def
	Size float32
}

// New creates a vertex colored cube mesh with a given size
func New(args Args) *Mesh {
	if args.Mat == nil {
		args.Mat = material.ColoredForward()
	}
	cube := object.NewComponent(&Mesh{
		Static: mesh.New(args.Mat),
		Args:   args,
	})
	cube.SetTexture(texture.Diffuse, texture.Checker)
	cube.generate()
	return cube
}

func (c *Mesh) generate() {
	s := c.Size / 2

	topLeft := vec2.New(0, 0)
	topRight := vec2.New(1, 0)
	bottomLeft := vec2.New(0, 1)
	bottomRight := vec2.New(1, 1)

	vertices := []vertex.T{
		// X+
		{P: vec3.New(s, -s, s), N: vec3.UnitX, T: bottomRight}, // 0
		{P: vec3.New(s, -s, -s), N: vec3.UnitX, T: bottomLeft}, // 1
		{P: vec3.New(s, s, -s), N: vec3.UnitX, T: topLeft},     // 2
		{P: vec3.New(s, s, s), N: vec3.UnitX, T: topRight},     // 3

		// X-
		{P: vec3.New(-s, -s, -s), N: vec3.UnitXN, T: bottomRight}, // 4
		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, T: bottomLeft},   // 5
		{P: vec3.New(-s, s, s), N: vec3.UnitXN, T: topLeft},       // 6
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, T: topRight},     // 7

		// Y+
		{P: vec3.New(s, s, -s), N: vec3.UnitY, T: bottomRight}, // 8
		{P: vec3.New(-s, s, -s), N: vec3.UnitY, T: bottomLeft}, // 9
		{P: vec3.New(-s, s, s), N: vec3.UnitY, T: topLeft},     // 10
		{P: vec3.New(s, s, s), N: vec3.UnitY, T: topRight},     // 11

		// Y-
		{P: vec3.New(-s, -s, -s), N: vec3.UnitYN, T: bottomRight}, // 12
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, T: bottomLeft},   // 13
		{P: vec3.New(s, -s, s), N: vec3.UnitYN, T: topLeft},       // 14
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, T: topRight},     // 15

		// Z+
		{P: vec3.New(-s, -s, s), N: vec3.UnitZ, T: bottomRight}, // 16
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, T: bottomLeft},   // 17
		{P: vec3.New(s, s, s), N: vec3.UnitZ, T: topLeft},       // 18
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, T: topRight},     // 19

		// Z-
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, T: bottomRight}, // 20
		{P: vec3.New(-s, -s, -s), N: vec3.UnitZN, T: bottomLeft}, // 21
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, T: topLeft},     // 22
		{P: vec3.New(s, s, -s), N: vec3.UnitZN, T: topRight},     // 23
	}

	indices := []uint16{
		0, 1, 2,
		0, 2, 3,

		4, 5, 6,
		4, 6, 7,

		8, 9, 10,
		8, 10, 11,

		12, 13, 14,
		12, 14, 15,

		16, 17, 18,
		16, 18, 19,

		20, 21, 22,
		20, 22, 23,
	}

	key := object.Key("cube", c)
	mesh := vertex.NewTriangles(key, vertices, indices)
	c.VertexData.Set(mesh)
}
