package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Box struct {
	*mesh.Static
	BoxArgs
}

// Args are kinda like props
// If they change, we should recomupte the mesh

type BoxArgs struct {
	Extents vec3.T
	Color   color.T
}

func NewBox(args BoxArgs) *Box {
	b := object.NewComponent(&Box{
		Static:  mesh.NewLines(),
		BoxArgs: args,
	})
	b.compute()
	return b
}

func (b *Box) compute() {
	var x, y, z float32
	w, h, d := b.Extents.X/2, b.Extents.Y/2, b.Extents.Z/2
	c := b.Color.Vec4()

	key := object.Key("box", b)
	mesh := vertex.NewLines(key, []vertex.C{
		// bottom square
		{P: vec3.New(x-w, y-h, z-d), C: c}, // 0
		{P: vec3.New(x+w, y-h, z-d), C: c}, // 1
		{P: vec3.New(x-w, y-h, z+d), C: c}, // 2
		{P: vec3.New(x+w, y-h, z+d), C: c}, // 3

		// top square
		{P: vec3.New(x-w, y+h, z-d), C: c}, // 4
		{P: vec3.New(x+w, y+h, z-d), C: c}, // 5
		{P: vec3.New(x-w, y+h, z+d), C: c}, // 6
		{P: vec3.New(x+w, y+h, z+d), C: c}, // 7
	}, []uint16{
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
	})
	b.VertexData.Set(mesh)
}
