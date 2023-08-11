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
	Radius object.Property[float32]
	Color  object.Property[color.T]

	data vertex.MutableMesh[vertex.C, uint16]
}

type SphereArgs struct {
	Radius float32
	Color  color.T
}

func NewSphere(args SphereArgs) *Sphere {
	b := object.NewComponent(&Sphere{
		Static: mesh.NewLines(),
		Radius: object.NewProperty(args.Radius),
		Color:  object.NewProperty(args.Color),
	})
	b.Radius.OnChange.Subscribe(func(float32) { b.refresh() })
	b.Color.OnChange.Subscribe(func(color.T) { b.refresh() })
	b.data = vertex.NewLines[vertex.C, uint16](object.Key("sphere", b), nil, nil)
	b.refresh()
	return b
}

func (b *Sphere) refresh() {
	r := b.Radius.Get()
	color := b.Color.Get().Vec4()
	segments := 32
	angle := 2 * math.Pi / float32(segments)
	vertices := make([]vertex.C, 0, 2*3*segments)

	// x ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a0), 0, math.Sin(a0)).Scaled(r),
			C: color,
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a1), 0, math.Sin(a1)).Scaled(r),
			C: color,
		})
	}

	// y ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(0, math.Sin(a0), math.Cos(a0)).Scaled(r),
			C: color,
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(0, math.Sin(a1), math.Cos(a1)).Scaled(r),
			C: color,
		})
	}

	// z ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a0), math.Sin(a0), 0).Scaled(r),
			C: color,
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a1), math.Sin(a1), 0).Scaled(r),
			C: color,
		})
	}

	b.data.Update(vertices, []uint16{})
	b.VertexData.Set(b.data)
}
