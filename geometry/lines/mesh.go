package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	*mesh.Static
	Color  object.Property[color.T]
	Source object.Property[vertex.Mesh]

	data vertex.MutableMesh[vertex.C, uint32]
}

func NewMesh(msh vertex.Mesh, clr color.T) *Mesh {
	b := object.NewComponent(&Mesh{
		Static: mesh.NewLines(),
		Color:  object.NewProperty(clr),
		Source: object.NewProperty(msh),
	})
	b.Color.OnChange.Subscribe(func(color.T) { b.refresh() })
	b.Source.OnChange.Subscribe(func(vertex.Mesh) { b.refresh() })
	b.data = vertex.NewLines[vertex.C, uint32](object.Key("wiremesh", b), nil, nil)
	b.refresh()
	return b
}

func (b *Mesh) refresh() {
	clr := b.Color.Get().Vec4()
	msh := b.Source.Get()

	indices := make([]uint32, 0, msh.IndexCount()*2)
	vertices := make([]vertex.C, 0, msh.VertexCount())

	msh.Triangles(func(t vertex.Triangle) {
		offset := uint32(len(vertices))
		vertices = append(vertices, vertex.C{P: t.A, C: clr})
		vertices = append(vertices, vertex.C{P: t.B, C: clr})
		vertices = append(vertices, vertex.C{P: t.C, C: clr})

		indices = append(indices, offset+0, offset+1) // A-B
		indices = append(indices, offset+1, offset+2) // B-C
		indices = append(indices, offset+2, offset+0) // C-A
	})

	b.data.Update(vertices, indices)
	b.VertexData.Set(b.data)
}
