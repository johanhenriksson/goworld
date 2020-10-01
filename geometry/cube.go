package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
)

type Cube struct {
	*engine.Mesh
}

func NewCube(parent *engine.Object) *Cube {
	mat := assets.GetMaterialCached("default")
	cube := &Cube{
		Mesh: engine.NewMesh(parent, mat),
	}
	cube.generate()
	return cube
}

func (c *Cube) generate() {
	data := Vertices{
		Vertex{1, -1, 1, 1, 0, 0, 1, 1},
		Vertex{1, -1, -1, 1, 0, 0, 1, 0},
		Vertex{1, 1, -1, 1, 0, 0, 0, 0},
		Vertex{1, -1, 1, 1, 0, 0, 1, 1},
		Vertex{1, 1, -1, 1, 0, 0, 0, 0},
		Vertex{1, 1, 1, 1, 0, 0, 0, 1},

		Vertex{-1, -1, 1, -1, 0, 0, 0, 1},
		Vertex{-1, 1, -1, -1, 0, 0, 1, 0},
		Vertex{-1, -1, -1, -1, 0, 0, 0, 0},
		Vertex{-1, -1, 1, -1, 0, 0, 0, 1},
		Vertex{-1, 1, 1, -1, 0, 0, 1, 1},
		Vertex{-1, 1, -1, -1, 0, 0, 1, 0},

		Vertex{-1, 1, -1, 0, 1, 0, 0, 0},
		Vertex{-1, 1, 1, 0, 1, 0, 0, 1},
		Vertex{1, 1, -1, 0, 1, 0, 1, 0},
		Vertex{1, 1, -1, 0, 1, 0, 1, 0},
		Vertex{-1, 1, 1, 0, 1, 0, 0, 1},
		Vertex{1, 1, 1, 0, 1, 0, 1, 1},

		Vertex{-1, -1, -1, 0, -1, 0, 0, 0},
		Vertex{1, -1, -1, 0, -1, 0, 1, 0},
		Vertex{-1, -1, 1, 0, -1, 0, 0, 1},
		Vertex{1, -1, -1, 0, -1, 0, 1, 0},
		Vertex{1, -1, 1, 0, -1, 0, 1, 1},
		Vertex{-1, -1, 1, 0, -1, 0, 0, 1},

		Vertex{-1, -1, 1, 0, 0, 1, 1, 0},
		Vertex{1, -1, 1, 0, 0, 1, 0, 0},
		Vertex{-1, 1, 1, 0, 0, 1, 1, 1},
		Vertex{1, -1, 1, 0, 0, 1, 0, 0},
		Vertex{1, 1, 1, 0, 0, 1, 0, 1},
		Vertex{-1, 1, 1, 0, 0, 1, 1, 1},

		Vertex{-1, -1, -1, 0, 0, -1, 0, 0},
		Vertex{-1, 1, -1, 0, 0, -1, 0, 1},
		Vertex{1, -1, -1, 0, 0, -1, 1, 0},
		Vertex{1, -1, -1, 0, 0, -1, 1, 0},
		Vertex{-1, 1, -1, 0, 0, -1, 0, 1},
		Vertex{1, 1, -1, 0, 0, -1, 1, 1},
	}
	c.Buffer("geometry", data)
}
