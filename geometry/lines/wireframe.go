package lines

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Wireframe struct {
	*mesh.Static
	Color  object.Property[color.T]
	Source object.Property[assets.Mesh]

	data   vertex.MutableMesh[vertex.Vertex, uint32]
	offset float32
}

func NewWireframe(pool object.Pool, msh assets.Mesh, clr color.T) *Wireframe {
	w := object.NewComponent(pool, &Wireframe{
		Static: mesh.New(pool, nil),
		Color:  object.NewProperty(clr),
		Source: object.NewProperty(msh),
	})
	w.Color.OnChange.Subscribe(func(color.T) { w.refresh() })
	w.Source.OnChange.Subscribe(func(assets.Mesh) { w.refresh() })
	w.data = vertex.NewLines[vertex.Vertex, uint32](object.Key("wireframe", w), nil, nil)
	w.refresh()
	return w
}

func (w *Wireframe) refresh() {
	clr := w.Color.Get()
	ref := w.Source.Get()
	if ref == nil {
		w.data.Update(nil, nil)
		w.VertexData.Set(w.data)
		return
	}

	// this is only slightly illegal since we actually need access to the mesh data
	// ideally wireframe should be done using another method that doesn't require
	// the mesh data to be manually converted to lines
	msh := ref.LoadMesh(assets.FS)

	indices := make([]uint32, 0, msh.IndexCount()*2)
	vertices := make([]vertex.Vertex, 0, msh.VertexCount())

	for t := range msh.Triangles() {
		index := uint32(len(vertices))
		offset := t.Normal().Scaled(w.offset)
		vertices = append(vertices, vertex.Vertex{P: t.A.Add(offset), C: clr}) // +0
		vertices = append(vertices, vertex.Vertex{P: t.B.Add(offset), C: clr}) // +1
		vertices = append(vertices, vertex.Vertex{P: t.C.Add(offset), C: clr}) // +2

		indices = append(indices, index+0, index+1) // A-B
		indices = append(indices, index+1, index+2) // B-C
		indices = append(indices, index+2, index+0) // C-A
	}

	w.data.Update(vertices, indices)
	w.VertexData.Set(w.data)
}
