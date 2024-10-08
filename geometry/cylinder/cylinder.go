package cylinder

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

func init() {
	Register[*Mesh](Type{
		Name: "Cyllinder",
		Path: []string{"Geometry"},
		Create: func(ctx Pool) (Component, error) {
			return New(ctx, Args{
				Mat:      material.StandardDeferred(),
				Radius:   0.5,
				Height:   1.8,
				Segments: 32,
				Color:    color.White,
			}), nil
		},
	})
}

type Cylinder struct {
	Object
	Mesh     *Mesh
	Collider *physics.Mesh
}

func New(pool Pool, args Args) *Cylinder {
	return NewObject(pool, "Cylinder", &Cylinder{
		Mesh:     NewMesh(pool, args),
		Collider: physics.NewMesh(pool),
	})
}

// A Cylinder is a forward rendered colored cyllinder mesh
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
		args.Mat = material.StandardForward()
	}
	cyllinder := NewComponent(pool, &Mesh{
		Static:   mesh.New(pool, args.Mat),
		Radius:   NewProperty(args.Radius),
		Height:   NewProperty(args.Height),
		Segments: NewProperty(args.Segments),
		Color:    NewProperty(args.Color),
	})
	cyllinder.Radius.OnChange.Subscribe(func(radius float32) { cyllinder.generate() })
	cyllinder.Height.OnChange.Subscribe(func(height float32) { cyllinder.generate() })
	cyllinder.Segments.OnChange.Subscribe(func(segments int) { cyllinder.generate() })
	cyllinder.Color.OnChange.Subscribe(func(color color.T) { cyllinder.generate() })
	// this should not run on the main thread
	cyllinder.generate()
	return cyllinder
}

func (c *Mesh) generate() {
	// vertex order: clockwise
	radius := c.Radius.Get()
	height := c.Height.Get()
	segments := c.Segments.Get()
	color := c.Color.Get()

	data := make([]vertex.Vertex, 2*2*3*segments)
	hh := height / 2
	sangle := 2 * math.Pi / float32(segments)

	// top
	top := vec3.New(0, hh, 0)
	bottom := vec3.New(0, -hh, 0)
	for i := 0; i < segments; i++ {
		o := 12 * i // segment vertex offset

		right := sangle * (float32(i) + 0.5)
		left := sangle * (float32(i) + 1.5)
		topRight := vec3.New(math.Cos(right), 0, -math.Sin(right)).Scaled(radius)
		topRight.Y = hh
		topLeft := vec3.New(math.Cos(left), 0, -math.Sin(left)).Scaled(radius)
		topLeft.Y = hh
		bottomRight := vec3.New(math.Cos(right), 0, -math.Sin(right)).Scaled(radius)
		bottomRight.Y = -hh
		bottomLeft := vec3.New(math.Cos(left), 0, -math.Sin(left)).Scaled(radius)
		bottomLeft.Y = -hh

		// top face
		data[o+0] = vertex.New(topLeft, vec3.Up, vec2.Zero, color)
		data[o+1] = vertex.New(top, vec3.Up, vec2.Zero, color)
		data[o+2] = vertex.New(topRight, vec3.Up, vec2.Zero, color)

		// bottom face
		data[o+3] = vertex.New(bottomRight, vec3.Down, vec2.Zero, color)
		data[o+4] = vertex.New(bottom, vec3.Down, vec2.Zero, color)
		data[o+5] = vertex.New(bottomLeft, vec3.Down, vec2.Zero, color)

		// calculate segment normal
		nv1 := topRight.Sub(bottomLeft)
		nv2 := bottomRight.Sub(bottomLeft)
		n := vec3.Cross(nv1, nv2)

		// side face 1
		data[o+6] = vertex.New(topRight, n, vec2.Zero, color)
		data[o+7] = vertex.New(bottomLeft, n, vec2.Zero, color)
		data[o+8] = vertex.New(topLeft, n, vec2.Zero, color)

		// side face 2
		data[o+9] = vertex.New(bottomRight, n, vec2.Zero, color)
		data[o+10] = vertex.New(bottomLeft, n, vec2.Zero, color)
		data[o+11] = vertex.New(topRight, n, vec2.Zero, color)
	}

	key := Key("cylinder", c)
	mesh := vertex.NewTriangles(key, data, []uint32{})
	c.VertexData.Set(mesh)
}
