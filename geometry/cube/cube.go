package cube

import (
	. "github.com/johanhenriksson/goworld/core/object"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	Register[*Mesh](TypeInfo{
		Name: "Cube",
		Path: []string{"Geometry"},
		Create: func(ctx Pool) (Component, error) {
			return New(ctx, Args{
				Size: 1,
			}), nil
		},
	})
}

type CubeObject struct {
	Object
	Mesh *Mesh
}

func New(pool Pool, args Args) *CubeObject {
	return NewObject(pool, "Cube", &CubeObject{
		Mesh: NewMesh(pool, args),
	})
}

// Mesh is a vertex colored cube mesh
type Mesh struct {
	// these are actually dynamic meshes, but since they generate so quickly
	// it might not make sense to generate in the background
	*mesh.Static

	Size Property[float32]
	Mat  Property[*material.Def]
}

type Args struct {
	Mat  *material.Def
	Size float32
}

// NewMesh creates a vertex colored cube mesh with a given size
func NewMesh(pool Pool, args Args) *Mesh {
	if args.Mat == nil {
		args.Mat = material.StandardForward()
	}
	c := NewComponent(pool, &Mesh{
		Static: mesh.New(pool, args.Mat),
		Size:   NewProperty(args.Size),
		Mat:    NewProperty(args.Mat),
	})
	c.Size.OnChange.Subscribe(func(size float32) {
		c.generate()
	})
	// todo: subscribe to material changes
	c.generate()
	return c
}

func (c *Mesh) generate() {
	s := c.Size.Get() / 2
	if s < 0 {
		// return an empty mesh?
		s = 0
	}

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

	key := Key("cube", c)
	mesh := vertex.NewTriangles(key, vertices, indices)
	c.VertexData.Set(mesh)
}
