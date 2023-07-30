package box

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	*mesh.Static
	Args
}

// Args are kinda like props
// If they change, we should recomupte the mesh

type Args struct {
	Size  vec3.T
	Color color.T
}

func New(args Args) *Mesh {
	b := object.NewComponent(&Mesh{
		Static: mesh.NewLines(),
		Args:   args,
	})
	b.compute()
	return b
}

func (b *Mesh) compute() {
	var x, y, z float32
	w, h, d := b.Size.X, b.Size.Y, b.Size.Z
	c := b.Color.Vec4()

	key := object.Key("box", b)
	mesh := vertex.NewLines(key, []vertex.C{
		// bottom square
		{P: vec3.New(x, y, z), C: c},     // 0
		{P: vec3.New(x+w, y, z), C: c},   // 1
		{P: vec3.New(x, y, z+d), C: c},   // 2
		{P: vec3.New(x+w, y, z+d), C: c}, // 3

		// top square
		{P: vec3.New(x, y+h, z), C: c},     // 4
		{P: vec3.New(x+w, y+h, z), C: c},   // 5
		{P: vec3.New(x, y+h, z+d), C: c},   // 6
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
