package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
)

type Cube struct {
	*engine.Mesh
	Size float32
}

func NewCube(parent *engine.Object, size float32) *Cube {
	mat := assets.GetMaterialCached("default")
	cube := &Cube{
		Mesh: engine.NewMesh(parent, mat),
		Size: size,
	}
	cube.generate()
	return cube
}

func (c *Cube) generate() {
	s := c.Size / 2
	data := Vertices{
		Vertex{s, -s, s, 1, 0, 0, 1, 1},
		Vertex{s, -s, -s, 1, 0, 0, 1, 0},
		Vertex{s, s, -s, 1, 0, 0, 0, 0},
		Vertex{s, -s, s, 1, 0, 0, 1, 1},
		Vertex{s, s, -s, 1, 0, 0, 0, 0},
		Vertex{s, s, s, 1, 0, 0, 0, 1},

		Vertex{-s, -s, s, -1, 0, 0, 0, 1},
		Vertex{-s, s, -s, -1, 0, 0, 1, 0},
		Vertex{-s, -s, -s, -1, 0, 0, 0, 0},
		Vertex{-s, -s, s, -1, 0, 0, 0, 1},
		Vertex{-s, s, s, -1, 0, 0, 1, 1},
		Vertex{-s, s, -s, -1, 0, 0, 1, 0},

		Vertex{-s, s, -s, 0, 1, 0, 0, 0},
		Vertex{-s, s, s, 0, 1, 0, 0, 1},
		Vertex{s, s, -s, 0, 1, 0, 1, 0},
		Vertex{s, s, -s, 0, 1, 0, 1, 0},
		Vertex{-s, s, s, 0, 1, 0, 0, 1},
		Vertex{s, s, s, 0, 1, 0, 1, 1},

		Vertex{-s, -s, -s, 0, -1, 0, 0, 0},
		Vertex{s, -s, -s, 0, -1, 0, 1, 0},
		Vertex{-s, -s, s, 0, -1, 0, 0, 1},
		Vertex{s, -s, -s, 0, -1, 0, 1, 0},
		Vertex{s, -s, s, 0, -1, 0, 1, 1},
		Vertex{-s, -s, s, 0, -1, 0, 0, 1},

		Vertex{-s, -s, s, 0, 0, 1, 1, 0},
		Vertex{s, -s, s, 0, 0, 1, 0, 0},
		Vertex{-s, s, s, 0, 0, 1, 1, 1},
		Vertex{s, -s, s, 0, 0, 1, 0, 0},
		Vertex{s, s, s, 0, 0, 1, 0, 1},
		Vertex{-s, s, s, 0, 0, 1, 1, 1},

		Vertex{-s, -s, -s, 0, 0, -1, 0, 0},
		Vertex{-s, s, -s, 0, 0, -1, 0, 1},
		Vertex{s, -s, -s, 0, 0, -1, 1, 0},
		Vertex{s, -s, -s, 0, 0, -1, 1, 0},
		Vertex{-s, s, -s, 0, 0, -1, 0, 1},
		Vertex{s, s, -s, 0, 0, -1, 1, 1},
	}
	c.Buffer("geometry", data)
}
