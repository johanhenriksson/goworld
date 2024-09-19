package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type BoxObject struct {
	Object
	Box *Box
}

func NewBoxObject(pool Pool, args BoxArgs) *BoxObject {
	return NewObject(pool, "Box", &BoxObject{
		Box: NewBox(pool, args),
	})
}

type Box struct {
	*mesh.Static
	Extents Property[vec3.T]
	Color   Property[color.T]

	data vertex.MutableMesh[vertex.C, uint16]
}

// Args are kinda like props
// If they change, we should recomupte the mesh

type BoxArgs struct {
	Extents vec3.T
	Color   color.T
}

func NewBox(pool Pool, args BoxArgs) *Box {
	b := NewComponent(pool, &Box{
		Static:  mesh.New(pool, nil),
		Extents: NewProperty(args.Extents),
		Color:   NewProperty(args.Color),
	})
	b.data = vertex.NewLines[vertex.C, uint16](Key("box", b), nil, nil)
	b.Extents.OnChange.Subscribe(func(vec3.T) { b.refresh() })
	b.Color.OnChange.Subscribe(func(color.T) { b.refresh() })
	b.refresh()
	return b
}

func (b *Box) refresh() {
	halfsize := b.Extents.Get().Scaled(0.5)
	w, h, d := halfsize.X, halfsize.Y, halfsize.Z
	c := b.Color.Get()

	vertices := []vertex.C{
		// bottom square
		{P: vec3.New(-w, -h, -d), C: c}, // 0
		{P: vec3.New(+w, -h, -d), C: c}, // 1
		{P: vec3.New(-w, -h, +d), C: c}, // 2
		{P: vec3.New(+w, -h, +d), C: c}, // 3

		// top square
		{P: vec3.New(-w, +h, -d), C: c}, // 4
		{P: vec3.New(+w, +h, -d), C: c}, // 5
		{P: vec3.New(-w, +h, +d), C: c}, // 6
		{P: vec3.New(+w, +h, +d), C: c}, // 7
	}
	indices := []uint16{
		// bottom
		0, 1,
		0, 2,
		1, 3,
		2, 3,

		// top
		4, 5,
		4, 6,
		5, 7,
		6, 7,

		// sides
		0, 4,
		1, 5,
		2, 6,
		3, 7,
	}

	b.data.Update(vertices, indices)
	b.VertexData.Set(b.data)
}
