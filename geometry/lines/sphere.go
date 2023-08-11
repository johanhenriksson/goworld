package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Sphere struct {
	*mesh.Static
	SphereArgs
}

type SphereArgs struct {
	Radius float32
	XColor color.T
	YColor color.T
	ZColor color.T
}

func NewSphere(args SphereArgs) *Sphere {
	b := object.NewComponent(&Sphere{
		Static:     mesh.NewLines(),
		SphereArgs: args,
	})
	b.compute()
	return b
}

func (b *Sphere) compute() {
	segments := 32
	angle := 2 * math.Pi / float32(segments)
	vertices := make([]vertex.C, 0, 2*3*segments)

	// x ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a0), 0, math.Sin(a0)).Scaled(b.Radius),
			C: b.XColor.Vec4(),
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a1), 0, math.Sin(a1)).Scaled(b.Radius),
			C: b.XColor.Vec4(),
		})
	}

	// y ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(0, math.Sin(a0), math.Cos(a0)).Scaled(b.Radius),
			C: b.YColor.Vec4(),
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(0, math.Sin(a1), math.Cos(a1)).Scaled(b.Radius),
			C: b.YColor.Vec4(),
		})
	}

	// z ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a0), math.Sin(a0), 0).Scaled(b.Radius),
			C: b.ZColor.Vec4(),
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a1), math.Sin(a1), 0).Scaled(b.Radius),
			C: b.ZColor.Vec4(),
		})
	}

	key := object.Key("sphere", b)
	mesh := vertex.NewLines(key, vertices, []uint16{})
	b.VertexData.Set(mesh)
}
