package cone

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Cone struct {
	object.Object
	*Mesh
	Collider *physics.Mesh
}

func NewObject(args Args) *Cone {
	return object.New("Cone", &Cone{
		Mesh:     New(args),
		Collider: physics.NewMesh(),
	})
}

// A Cone is a forward rendered colored cone mesh
type Mesh struct {
	*mesh.Static
	Args
}

type Args struct {
	Mat      *material.Def
	Radius   float32
	Height   float32
	Segments int
	Color    color.T
}

func New(args Args) *Mesh {
	if args.Mat == nil {
		args.Mat = material.ColoredForward()
	}
	cone := object.NewComponent(&Mesh{
		Static: mesh.New(args.Mat),
		Args:   args,
	})
	cone.generate()
	return cone
}

func (c *Mesh) generate() {
	data := make([]vertex.C, 6*c.Segments)

	// cone
	top := vec3.New(0, c.Height, 0)
	sangle := 2 * math.Pi / float32(c.Segments)
	for i := 0; i < c.Segments; i++ {
		a1 := sangle * (float32(i) + 0.5)
		a2 := sangle * (float32(i) + 1.5)
		v1 := vec3.New(math.Cos(a1), 0, -math.Sin(a1)).Scaled(c.Radius)
		v2 := vec3.New(math.Cos(a2), 0, -math.Sin(a2)).Scaled(c.Radius)
		v1t, v2t := top.Sub(v1), top.Sub(v2)
		n := vec3.Cross(v1t, v2t).Normalized()

		o := 3 * i
		data[o+0] = vertex.C{P: v2, N: n, C: c.Color.Vec4()}
		data[o+1] = vertex.C{P: top, N: n, C: c.Color.Vec4()}
		data[o+2] = vertex.C{P: v1, N: n, C: c.Color.Vec4()}
	}

	// bottom
	base := vec3.Zero
	n := vec3.New(0, -1, 0)
	for i := 0; i < c.Segments; i++ {
		a1 := sangle * (float32(i) + 0.5)
		a2 := sangle * (float32(i) + 1.5)
		v1 := vec3.New(math.Cos(a1), 0, -math.Sin(a1)).Scaled(c.Radius)
		v2 := vec3.New(math.Cos(a2), 0, -math.Sin(a2)).Scaled(c.Radius)
		o := 3 * (i + c.Segments)
		data[o+0] = vertex.C{P: v1, N: n, C: c.Color.Vec4()}
		data[o+1] = vertex.C{P: base, N: n, C: c.Color.Vec4()}
		data[o+2] = vertex.C{P: v2, N: n, C: c.Color.Vec4()}
	}

	key := object.Key("cone", c)
	mesh := vertex.NewTriangles(key, data, []uint16{})
	c.VertexData.Set(mesh)
}
