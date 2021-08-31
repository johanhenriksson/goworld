package box

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T struct {
	*engine.Mesh
	Args
}

type Args struct {
	Size  vec3.T
	Color render.Color
}

func New(args Args) *T {
	b := &T{
		Mesh: engine.NewLineMesh(),
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
	parent := object.New("Box")
	return Attach(parent, args)
}

func Builder(out **T, args Args) *object.Builder {
	b := object.Build("Box")
	*out = New(args)
	return b.Attach(*out)
}

func (b *T) compute() {
	var x, y, z float32
	w, h, d := b.Size.X, b.Size.Y, b.Size.Z
	c := b.Color.Vec4()
	vertices := []vertex.C{
		// bottom square
		{P: vec3.New(x, y, z), C: c},
		{P: vec3.New(x+w, y, z), C: c},
		{P: vec3.New(x, y, z), C: c},
		{P: vec3.New(x, y, z+d), C: c},
		{P: vec3.New(x+w, y, z), C: c},
		{P: vec3.New(x+w, y, z+d), C: c},
		{P: vec3.New(x, y, z+w), C: c},
		{P: vec3.New(x+w, y, z+d), C: c},

		// top square
		{P: vec3.New(x, y+h, z), C: c},
		{P: vec3.New(x+w, y+h, z), C: c},
		{P: vec3.New(x, y+h, z), C: c},
		{P: vec3.New(x, y+h, z+d), C: c},
		{P: vec3.New(x+w, y+h, z), C: c},
		{P: vec3.New(x+w, y+h, z+d), C: c},
		{P: vec3.New(x, y+h, z+w), C: c},
		{P: vec3.New(x+w, y+h, z+d), C: c},

		// connecting lines
		{P: vec3.New(x, y, z), C: c},
		{P: vec3.New(x, y+h, z), C: c},
		{P: vec3.New(x+w, y, z), C: c},
		{P: vec3.New(x+w, y+h, z), C: c},
		{P: vec3.New(x, y, z+d), C: c},
		{P: vec3.New(x, y+h, z+d), C: c},
		{P: vec3.New(x+w, y, z+d), C: c},
		{P: vec3.New(x+w, y+h, z+d), C: c},
	}
	b.Buffer(vertices)
}
