package geometry

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Box struct {
	*engine.Mesh
	Size  vec3.T
	Color render.Color
}

func NewBox(size vec3.T, color render.Color) *Box {
	b := &Box{
		Mesh:  engine.NewLineMesh(),
		Size:  size,
		Color: color,
	}
	b.compute()
	return b
}

func (b *Box) compute() {
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
