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

	data   vertex.MutableMesh[vertex.C, uint16]
	xcolor color.T
	ycolor color.T
	zcolor color.T
}

type SphereArgs struct {
	Radius float32
	Color  color.T
}

func NewSphere(pool object.Pool, args SphereArgs) *Sphere {
	b := object.NewComponent(pool, &Sphere{
		Static: mesh.New(pool, nil),
		Radius: object.NewProperty(args.Radius),
		Color:  object.NewProperty(args.Color),
	})
	b.Radius.OnChange.Subscribe(func(float32) { b.refresh() })
	b.Color.OnChange.Subscribe(func(c color.T) {
		b.SetAxisColors(c, c, c)
		b.refresh()
	})
	b.data = vertex.NewLines[vertex.C, uint16](object.Key("sphere", b), nil, nil)
	b.SetAxisColors(args.Color, args.Color, args.Color)
	return b
}

func (b *Sphere) SetAxisColors(x color.T, y color.T, z color.T) {
	b.xcolor = x
	b.ycolor = y
	b.zcolor = z
	b.refresh()
}

func (b *Sphere) refresh() {
	segments := 32
	radius := b.Radius.Get()
	angle := 2 * math.Pi / float32(segments)
	vertices := make([]vertex.C, 0, 2*3*segments)

	// x ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(0, math.Sin(a0), math.Cos(a0)).Scaled(radius),
			C: b.xcolor,
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(0, math.Sin(a1), math.Cos(a1)).Scaled(radius),
			C: b.xcolor,
		})
	}

	// y ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a0), 0, math.Sin(a0)).Scaled(radius),
			C: b.ycolor,
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a1), 0, math.Sin(a1)).Scaled(radius),
			C: b.ycolor,
		})
	}

	// z ring
	for i := 0; i < segments; i++ {
		a0 := float32(i) * angle
		a1 := float32(i+1) * angle
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a0), math.Sin(a0), 0).Scaled(radius),
			C: b.zcolor,
		})
		vertices = append(vertices, vertex.C{
			P: vec3.New(math.Cos(a1), math.Sin(a1), 0).Scaled(radius),
			C: b.zcolor,
		})
	}

	b.data.Update(vertices, nil)
	b.VertexData.Set(b.data)
}
