package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// ColorCube is a vertex colored cube mesh
type ColorCube struct {
	*engine.Mesh
	Size  float32
	Color render.Color
}

// NewColorCube creates a vertex colored cube mesh with a given size
func NewColorCube(color render.Color, size float32) *ColorCube {
	mat := assets.GetMaterialShared("color.f")
	cube := &ColorCube{
		Mesh:  engine.NewMesh("ColorCube", mat),
		Size:  size,
		Color: color,
	}
	cube.Pass = render.Forward
	cube.generate()
	return cube
}

func (c *ColorCube) generate() {
	s := c.Size / 2
	co := c.Color.Vec4()
	data := []vertex.C{
		{P: vec3.New(s, -s, s), N: vec3.UnitX, C: co},
		{P: vec3.New(s, -s, -s), N: vec3.UnitX, C: co},
		{P: vec3.New(s, s, -s), N: vec3.UnitX, C: co},
		{P: vec3.New(s, -s, s), N: vec3.UnitX, C: co},
		{P: vec3.New(s, s, -s), N: vec3.UnitX, C: co},
		{P: vec3.New(s, s, s), N: vec3.UnitX, C: co},

		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, C: co},
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, C: co},
		{P: vec3.New(-s, -s, -s), N: vec3.UnitXN, C: co},
		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, C: co},
		{P: vec3.New(-s, s, s), N: vec3.UnitXN, C: co},
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, C: co},

		{P: vec3.New(-s, s, -s), N: vec3.UnitY, C: co},
		{P: vec3.New(-s, s, s), N: vec3.UnitY, C: co},
		{P: vec3.New(s, s, -s), N: vec3.UnitY, C: co},
		{P: vec3.New(s, s, -s), N: vec3.UnitY, C: co},
		{P: vec3.New(-s, s, s), N: vec3.UnitY, C: co},
		{P: vec3.New(s, s, s), N: vec3.UnitY, C: co},

		{P: vec3.New(-s, -s, -s), N: vec3.UnitYN, C: co},
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, C: co},
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, C: co},
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, C: co},
		{P: vec3.New(s, -s, s), N: vec3.UnitYN, C: co},
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, C: co},

		{P: vec3.New(-s, -s, s), N: vec3.UnitZ, C: co},
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, C: co},
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, C: co},
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, C: co},
		{P: vec3.New(s, s, s), N: vec3.UnitZ, C: co},
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, C: co},

		{P: vec3.New(-s, -s, -s), N: vec3.UnitZN, C: co},
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, C: co},
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, C: co},
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, C: co},
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, C: co},
		{P: vec3.New(s, s, -s), N: vec3.UnitZN, C: co},
	}
	c.Buffer(data)
}
