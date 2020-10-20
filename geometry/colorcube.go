package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// ColorCube is a vertex colored cube mesh
type ColorCube struct {
	*engine.Mesh
	Size  float32
	Color render.Color
}

// NewCube creates a vertex colored cube mesh with a given size
func NewColorCube(color render.Color, size float32) *ColorCube {
	mat := assets.GetMaterialCached("vertex_color")
	cube := &ColorCube{
		Mesh:  engine.NewMesh(mat),
		Size:  size,
		Color: color,
	}
	cube.Passes.Set(render.Forward)
	cube.generate()
	return cube
}

func (c *ColorCube) generate() {
	s := c.Size / 2
	data := ColorVertices{
		ColorVertex{vec3.New(s, -s, s), vec3.UnitX, c.Color},
		ColorVertex{vec3.New(s, -s, -s), vec3.UnitX, c.Color},
		ColorVertex{vec3.New(s, s, -s), vec3.UnitX, c.Color},
		ColorVertex{vec3.New(s, -s, s), vec3.UnitX, c.Color},
		ColorVertex{vec3.New(s, s, -s), vec3.UnitX, c.Color},
		ColorVertex{vec3.New(s, s, s), vec3.UnitX, c.Color},

		ColorVertex{vec3.New(-s, -s, s), vec3.UnitXN, c.Color},
		ColorVertex{vec3.New(-s, s, -s), vec3.UnitXN, c.Color},
		ColorVertex{vec3.New(-s, -s, -s), vec3.UnitXN, c.Color},
		ColorVertex{vec3.New(-s, -s, s), vec3.UnitXN, c.Color},
		ColorVertex{vec3.New(-s, s, s), vec3.UnitXN, c.Color},
		ColorVertex{vec3.New(-s, s, -s), vec3.UnitXN, c.Color},

		ColorVertex{vec3.New(-s, s, -s), vec3.UnitY, c.Color},
		ColorVertex{vec3.New(-s, s, s), vec3.UnitY, c.Color},
		ColorVertex{vec3.New(s, s, -s), vec3.UnitY, c.Color},
		ColorVertex{vec3.New(s, s, -s), vec3.UnitY, c.Color},
		ColorVertex{vec3.New(-s, s, s), vec3.UnitY, c.Color},
		ColorVertex{vec3.New(s, s, s), vec3.UnitY, c.Color},

		ColorVertex{vec3.New(-s, -s, -s), vec3.UnitYN, c.Color},
		ColorVertex{vec3.New(s, -s, -s), vec3.UnitYN, c.Color},
		ColorVertex{vec3.New(-s, -s, s), vec3.UnitYN, c.Color},
		ColorVertex{vec3.New(s, -s, -s), vec3.UnitYN, c.Color},
		ColorVertex{vec3.New(s, -s, s), vec3.UnitYN, c.Color},
		ColorVertex{vec3.New(-s, -s, s), vec3.UnitYN, c.Color},

		ColorVertex{vec3.New(-s, -s, s), vec3.UnitZ, c.Color},
		ColorVertex{vec3.New(s, -s, s), vec3.UnitZ, c.Color},
		ColorVertex{vec3.New(-s, s, s), vec3.UnitZ, c.Color},
		ColorVertex{vec3.New(s, -s, s), vec3.UnitZ, c.Color},
		ColorVertex{vec3.New(s, s, s), vec3.UnitZ, c.Color},
		ColorVertex{vec3.New(-s, s, s), vec3.UnitZ, c.Color},

		ColorVertex{vec3.New(-s, -s, -s), vec3.UnitZN, c.Color},
		ColorVertex{vec3.New(-s, s, -s), vec3.UnitZN, c.Color},
		ColorVertex{vec3.New(s, -s, -s), vec3.UnitZN, c.Color},
		ColorVertex{vec3.New(s, -s, -s), vec3.UnitZN, c.Color},
		ColorVertex{vec3.New(-s, s, -s), vec3.UnitZN, c.Color},
		ColorVertex{vec3.New(s, s, -s), vec3.UnitZN, c.Color},
	}
	c.Buffer("geometry", data)
}
