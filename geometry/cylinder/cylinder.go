package cylinder

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Cylinder struct {
	object.Object
	Mesh *Mesh
}

func NewObject(args Args) *Cylinder {
	return object.New("Cyllinder", &Cylinder{
		Mesh: New(args),
	})
}

// A Cylinder is a forward rendered colored cyllinder mesh
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
	cyllinder := object.NewComponent(&Mesh{
		Static: mesh.New(mesh.Forward, args.Mat),
		Args:   args,
	})
	// this should not run on the main thread
	cyllinder.generate()
	return cyllinder
}

func (c *Mesh) generate() {
	// vertex order: clockwise

	data := make([]vertex.C, 2*2*3*c.Segments)
	hh := c.Height / 2
	sangle := 2 * math.Pi / float32(c.Segments)
	color := c.Color.Vec4()

	// top
	top := vec3.New(0, hh, 0)
	bottom := vec3.New(0, -hh, 0)
	for i := 0; i < c.Segments; i++ {
		o := 12 * i // segment vertex offset

		right := sangle * (float32(i) + 0.5)
		left := sangle * (float32(i) + 1.5)
		topRight := vec3.New(math.Cos(right), 0, -math.Sin(right)).Scaled(c.Radius)
		topRight.Y = hh
		topLeft := vec3.New(math.Cos(left), 0, -math.Sin(left)).Scaled(c.Radius)
		topLeft.Y = hh
		bottomRight := vec3.New(math.Cos(right), 0, -math.Sin(right)).Scaled(c.Radius)
		bottomRight.Y = -hh
		bottomLeft := vec3.New(math.Cos(left), 0, -math.Sin(left)).Scaled(c.Radius)
		bottomLeft.Y = -hh

		// top face
		data[o+0] = vertex.C{P: topLeft, N: vec3.Up, C: color}
		data[o+1] = vertex.C{P: top, N: vec3.Up, C: color}
		data[o+2] = vertex.C{P: topRight, N: vec3.Up, C: color}

		// bottom face
		data[o+3] = vertex.C{P: bottomRight, N: vec3.Down, C: color}
		data[o+4] = vertex.C{P: bottom, N: vec3.Down, C: color}
		data[o+5] = vertex.C{P: bottomLeft, N: vec3.Down, C: color}

		// calculate segment normal
		nv1 := topRight.Sub(bottomLeft)
		nv2 := bottomRight.Sub(bottomLeft)
		n := vec3.Cross(nv1, nv2)

		// side face 1
		data[o+6] = vertex.C{P: topLeft, N: n, C: color}
		data[o+7] = vertex.C{P: bottomLeft, N: n, C: color}
		data[o+8] = vertex.C{P: topRight, N: n, C: color}

		// side face 2
		data[o+9] = vertex.C{P: topRight, N: n, C: color}
		data[o+10] = vertex.C{P: bottomLeft, N: n, C: color}
		data[o+11] = vertex.C{P: bottomRight, N: n, C: color}
	}

	key := object.Key("cyllinder", c)
	mesh := vertex.NewTriangles(key, data, []uint16{})
	c.VertexData.Set(mesh)
}
