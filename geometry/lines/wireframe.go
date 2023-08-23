package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Wireframe struct {
	*mesh.Static
	Color  object.Property[color.T]
	Source object.Property[vertex.Mesh]

	data   vertex.MutableMesh[vertex.C, uint32]
	offset float32
}

func NewWireframe(msh vertex.Mesh, clr color.T) *Wireframe {
	w := object.NewComponent(&Wireframe{
		Static: mesh.NewLines(),
		Color:  object.NewProperty(clr),
		Source: object.NewProperty(msh),
	})
	w.Color.OnChange.Subscribe(func(color.T) { w.refresh() })
	w.Source.OnChange.Subscribe(func(vertex.Mesh) { w.refresh() })
	w.data = vertex.NewLines[vertex.C, uint32](object.Key("wireframe", w), nil, nil)
	w.refresh()
	return w
}

func (w *Wireframe) refresh() {
	clr := w.Color.Get().Vec4()
	msh := w.Source.Get()
	if msh == nil {
		return
	}

	indices := make([]uint32, 0, msh.IndexCount()*2)
	vertices := make([]vertex.C, 0, msh.VertexCount())

	msh.Triangles(func(t vertex.Triangle) {
		index := uint32(len(vertices))
		offset := t.Normal().Scaled(w.offset)
		vertices = append(vertices, vertex.C{P: t.A.Add(offset), C: clr}) // +0
		vertices = append(vertices, vertex.C{P: t.B.Add(offset), C: clr}) // +1
		vertices = append(vertices, vertex.C{P: t.C.Add(offset), C: clr}) // +2

		indices = append(indices, index+0, index+1) // A-B
		indices = append(indices, index+1, index+2) // B-C
		indices = append(indices, index+2, index+0) // C-A
	})

	w.data.Update(vertices, indices)
	w.VertexData.Set(w.data)
}
