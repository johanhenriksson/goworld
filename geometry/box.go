package geometry

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Box struct {
	*engine.Mesh
	Size  vec3.T
	Color render.Color
}

func NewBox(size vec3.T, color render.Color) *Box {
	b := &Box{
		Mesh:  engine.NewLineMesh("Box"),
		Size:  size,
		Color: color,
	}
	b.compute()
	return b
}

func (b *Box) compute() {
	x, y, z := b.Position.X, b.Position.Y, b.Position.Z
	w, h, d := b.Size.X, b.Size.Y, b.Size.Z
	vertices := ColorVertices{
		// bottom square
		ColorVertex{Position: vec3.New(x, y, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y, z+d), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y, z+d), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y, z+w), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y, z+d), Color: b.Color},

		// top square
		ColorVertex{Position: vec3.New(x, y+h, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y+h, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y+h, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y+h, z+d), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y+h, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y+h, z+d), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y+h, z+w), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y+h, z+d), Color: b.Color},

		// connecting lines
		ColorVertex{Position: vec3.New(x, y, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y+h, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y+h, z), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y, z+d), Color: b.Color},
		ColorVertex{Position: vec3.New(x, y+h, z+d), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y, z+d), Color: b.Color},
		ColorVertex{Position: vec3.New(x+w, y+h, z+d), Color: b.Color},
	}
	b.Buffer("geometry", vertices)
}
