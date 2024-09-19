package cube

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets/fs"
	. "github.com/johanhenriksson/goworld/core/object"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	Register[*Mesh](Type{
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
	Mesh *mesh.Static
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
}

type Args struct {
	Mat  *material.Def
	Size float32
}

// NewMesh creates a vertex colored cube mesh with a given size
func NewMesh(pool Pool, args Args) *mesh.Static {
	if args.Mat == nil {
		args.Mat = material.StandardForward()
	}
	m := mesh.New(pool, args.Mat)
	m.VertexData.Set(newCube(args.Size))
	return m
}

type cube struct {
	key     string
	version int
	mesh    vertex.MutableMesh[vertex.T, uint16]
	size    float32
}

func newCube(size float32) *cube {
	return &cube{
		key:     fmt.Sprintf("cube(%f)", size),
		version: 1,
		size:    size,
	}
}

func (c *cube) generate() vertex.Mesh {
	s := c.size
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

	return vertex.NewTriangles(c.key, vertices, indices)
}

func (c *cube) Key() string  { return c.key }
func (c *cube) Version() int { return c.version }

func (c *cube) LoadMesh(fs fs.Filesystem) vertex.Mesh {
	// what is responsible for caching the result of this?
	// ideally, if the mesh already exists, the key/version will be the same and cause a cache hit before this is called
	// if store the result here, it will be cached once in every reference!

	// sidenote: we do have access to the file system
	// so we could potentially load a mesh from a file here, avoiding re-generation of heavier assets

	return c.generate()
}
