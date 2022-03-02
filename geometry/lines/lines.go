package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T struct {
	mesh.T
	Args
}

type Args struct {
	Lines []Line
}

func New(args Args) *T {
	b := &T{
		T:    mesh.NewLines(),
		Args: args,
	}
	b.compute()
	return b
}

func Attach(parent object.T, args Args) *T {
	box := New(args)
	parent.Attach(box)
	return box
}

func NewObject(args Args) *T {
	parent := object.New("Lines")
	return Attach(parent, args)
}

func Builder(out **T, args Args) *object.Builder {
	b := object.Build("Lines")
	*out = New(args)
	return b.Attach(*out)
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

	mesh := vertex.NewLines("lines", vertices, []uint16{})
	li.SetMesh(mesh)
}
