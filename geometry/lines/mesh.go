package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	*mesh.Static
	mesh  vertex.Mesh
	color color.T
}

func NewMesh(msh vertex.Mesh, clr color.T) *Mesh {
	b := object.NewComponent(&Mesh{
		Static: mesh.NewLines(),
		mesh:   msh,
		color:  clr,
	})
	b.compute()
	return b
}

func (b *Mesh) compute() {
	clr := b.color.Vec4()

	indices := make([]uint16, 0, b.mesh.IndexCount()*2)
	vertices := make([]vertex.C, 0, b.mesh.VertexCount())

	b.mesh.Triangles(func(t vertex.Triangle) {
		offset := uint16(len(vertices))
		vertices = append(vertices, vertex.C{P: t.A, C: clr})
		vertices = append(vertices, vertex.C{P: t.B, C: clr})
		vertices = append(vertices, vertex.C{P: t.C, C: clr})

		// A-B
		indices = append(indices, offset+0)
		indices = append(indices, offset+1)

		// B-C
		indices = append(indices, offset+1)
		indices = append(indices, offset+2)

		// C-A
		indices = append(indices, offset+2)
		indices = append(indices, offset+0)
	})

	key := object.Key("wireframe", b)
	mesh := vertex.NewLines(key, vertices, indices)
	b.VertexData.Set(mesh)
}
