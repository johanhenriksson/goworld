package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T struct {
	mesh.T
	Args
}

type Args struct {
	Mat   *material.Def
	Lines []Line

	lineMesh vertex.MutableMesh[vertex.C, uint16]
}

func New(args Args) *T {
	b := object.New(&T{
		T:    mesh.NewLines(args.Mat),
		Args: args,
	})
	b.lineMesh = vertex.NewLines(object.Key("lines", b), []vertex.C{}, []uint16{})
	b.SetMesh(b.lineMesh)
	b.Refresh()
	return b
}

func (li *T) Add(from, to vec3.T, clr color.T) {
	li.Lines = append(li.Lines, Line{
		Start: from,
		End:   to,
		Color: clr,
	})
}

func (li *T) Clear() {
	li.Lines = li.Lines[:0]
}

func (li *T) Count() int {
	return len(li.Lines)
}

func (li *T) Refresh() {
	count := len(li.Lines)
	vertices := make([]vertex.C, 2*count)
	for i := 0; i < count; i++ {
		line := li.Lines[i]
		a := &vertices[2*i+0]
		b := &vertices[2*i+1]
		a.P = line.Start
		a.C = line.Color.Vec4()
		b.P = line.End
		b.C = line.Color.Vec4()
	}
	li.lineMesh.Update(vertices, []uint16{})
}
