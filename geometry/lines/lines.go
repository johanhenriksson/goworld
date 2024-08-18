package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Lines struct {
	*mesh.Static
	Args
}

type Args struct {
	Lines []Line

	lineMesh vertex.MutableMesh[vertex.C, uint16]
}

func New(pool object.Pool, args Args) *Lines {
	b := object.NewComponent(pool, &Lines{
		Static: mesh.NewLines(pool),
		Args:   args,
	})
	b.lineMesh = vertex.NewLines(object.Key("lines", b), []vertex.C{}, []uint16{})
	b.VertexData.Set(b.lineMesh)
	b.Refresh()
	return b
}

func (li *Lines) Add(from, to vec3.T, clr color.T) {
	li.Lines = append(li.Lines, Line{
		Start: from,
		End:   to,
		Color: clr,
	})
}

func (li *Lines) Clear() {
	li.Lines = li.Lines[:0]
}

func (li *Lines) Count() int {
	return len(li.Lines)
}

func (li *Lines) Refresh() {
	count := len(li.Lines)
	vertices := make([]vertex.C, 2*count)
	for i := 0; i < count; i++ {
		line := li.Lines[i]
		a := &vertices[2*i+0]
		b := &vertices[2*i+1]
		a.P = line.Start
		a.C = line.Color
		b.P = line.End
		b.C = line.Color
	}
	li.lineMesh.Update(vertices, []uint16{})
}
