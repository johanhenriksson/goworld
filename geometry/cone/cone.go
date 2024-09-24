package cone

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Cone struct {
	Object
	Mesh     *Mesh
	Collider *physics.Mesh
}

func New(pool Pool, args Args) *Cone {
	return NewObject(pool, "Cone", &Cone{
		Mesh:     NewMesh(pool, args),
		Collider: physics.NewMesh(pool),
	})
}

// A Cone is a forward rendered colored cone mesh
type Mesh struct {
	*mesh.Static
	Radius   Property[float32]
	Height   Property[float32]
	Segments Property[int]
	Color    Property[color.T]
}

type Args struct {
	Mat      *material.Def
	Radius   float32
	Height   float32
	Segments int
	Color    color.T
}

func NewMesh(pool Pool, args Args) *Mesh {
	if args.Mat == nil {
		args.Mat = material.ColoredForward()
	}
	cone := NewComponent(pool, &Mesh{
		Static:   mesh.New(pool, args.Mat),
		Radius:   NewProperty(args.Radius),
		Height:   NewProperty(args.Height),
		Segments: NewProperty(args.Segments),
		Color:    NewProperty(args.Color),
	})
	cone.Radius.OnChange.Subscribe(func(radius float32) { cone.generate() })
	cone.Height.OnChange.Subscribe(func(height float32) { cone.generate() })
	cone.Segments.OnChange.Subscribe(func(segments int) { cone.generate() })
	cone.Color.OnChange.Subscribe(func(color color.T) { cone.generate() })
	cone.generate()
	return cone
}

func (c *Mesh) generate() {
	radius, height := c.Radius.Get(), c.Height.Get()
	color := c.Color.Get()
	segments := c.Segments.Get()

	data := make([]vertex.Vertex, 6*segments)

	// cone
	top := vec3.New(0, height, 0)
	sangle := 2 * math.Pi / float32(segments)
	for i := 0; i < segments; i++ {
		a1 := sangle * (float32(i) + 0.5)
		a2 := sangle * (float32(i) + 1.5)
		v1 := vec3.New(math.Cos(a1), 0, -math.Sin(a1)).Scaled(radius)
		v2 := vec3.New(math.Cos(a2), 0, -math.Sin(a2)).Scaled(radius)
		v1t, v2t := top.Sub(v1), top.Sub(v2)
		n := vec3.Cross(v1t, v2t).Normalized()

		o := 3 * i
		data[o+0] = vertex.New(v2, n, vec2.Zero, color)
		data[o+1] = vertex.New(top, n, vec2.Zero, color)
		data[o+2] = vertex.New(v1, n, vec2.Zero, color)
	}

	// bottom
	base := vec3.Zero
	n := vec3.New(0, -1, 0)
	for i := 0; i < segments; i++ {
		a1 := sangle * (float32(i) + 0.5)
		a2 := sangle * (float32(i) + 1.5)
		v1 := vec3.New(math.Cos(a1), 0, -math.Sin(a1)).Scaled(radius)
		v2 := vec3.New(math.Cos(a2), 0, -math.Sin(a2)).Scaled(radius)
		o := 3 * (i + segments)
		data[o+0] = vertex.New(v1, n, vec2.Zero, color)
		data[o+1] = vertex.New(base, n, vec2.Zero, color)
		data[o+2] = vertex.New(v2, n, vec2.Zero, color)
	}

	key := Key("cone", c)
	mesh := vertex.NewTriangles(key, data, []uint16{})
	c.VertexData.Set(mesh)
}
