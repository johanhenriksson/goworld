package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
)

// Cube mesh, textured
type Cube struct {
	*engine.Mesh
	Size float32
}

// NewCube creates a new textured cube mesh with a given size
func NewCube(size float32) *Cube {
	mat := assets.GetMaterialCached("default")
	cube := &Cube{
		Mesh: engine.NewMesh(mat),
		Size: size,
	}
	cube.generate()
	return cube
}

func (c *Cube) generate() {
	s := c.Size / 2
	data := Vertices{
		// XP
		Vertex{s, -s, s, 1, 0, 0, 1, 1},
		Vertex{s, -s, -s, 1, 0, 0, 1, 0},
		Vertex{s, s, -s, 1, 0, 0, 0, 0},
		Vertex{s, -s, s, 1, 0, 0, 1, 1},
		Vertex{s, s, -s, 1, 0, 0, 0, 0},
		Vertex{s, s, s, 1, 0, 0, 0, 1},

		// XN
		Vertex{-s, -s, s, -1, 0, 0, 0, 1},
		Vertex{-s, s, -s, -1, 0, 0, 1, 0},
		Vertex{-s, -s, -s, -1, 0, 0, 0, 0},
		Vertex{-s, -s, s, -1, 0, 0, 0, 1},
		Vertex{-s, s, s, -1, 0, 0, 1, 1},
		Vertex{-s, s, -s, -1, 0, 0, 1, 0},

		// YP
		Vertex{-s, s, -s, 0, 1, 0, 0, 0},
		Vertex{-s, s, s, 0, 1, 0, 0, 1},
		Vertex{s, s, -s, 0, 1, 0, 1, 0},
		Vertex{s, s, -s, 0, 1, 0, 1, 0},
		Vertex{-s, s, s, 0, 1, 0, 0, 1},
		Vertex{s, s, s, 0, 1, 0, 1, 1},

		// YN
		Vertex{-s, -s, -s, 0, -1, 0, 0, 0},
		Vertex{s, -s, -s, 0, -1, 0, 1, 0},
		Vertex{-s, -s, s, 0, -1, 0, 0, 1},
		Vertex{s, -s, -s, 0, -1, 0, 1, 0},
		Vertex{s, -s, s, 0, -1, 0, 1, 1},
		Vertex{-s, -s, s, 0, -1, 0, 0, 1},

		// ZP
		Vertex{-s, -s, s, 0, 0, 1, 1, 0},
		Vertex{s, -s, s, 0, 0, 1, 0, 0},
		Vertex{-s, s, s, 0, 0, 1, 1, 1},
		Vertex{s, -s, s, 0, 0, 1, 0, 0},
		Vertex{s, s, s, 0, 0, 1, 0, 1},
		Vertex{-s, s, s, 0, 0, 1, 1, 1},

		// ZN
		Vertex{-s, -s, -s, 0, 0, -1, 0, 0},
		Vertex{-s, s, -s, 0, 0, -1, 0, 1},
		Vertex{s, -s, -s, 0, 0, -1, 1, 0},
		Vertex{s, -s, -s, 0, 0, -1, 1, 0},
		Vertex{-s, s, -s, 0, 0, -1, 0, 1},
		Vertex{s, s, -s, 0, 0, -1, 1, 1},
	}
	c.Buffer("geometry", data)
}
