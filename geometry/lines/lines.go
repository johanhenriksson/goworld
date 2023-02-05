package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
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
}

func New(args Args) *T {
	b := object.New(&T{
		T:    mesh.NewLines(args.Mat),
		Args: args,
	})
	b.compute()
	return b
}

func (li *T) compute() {
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

	key := object.Key("lines", li)
	mesh := vertex.NewLines(key, vertices, []uint16{})
	li.SetMesh(mesh)
}
