package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Cone mesh
type Cone struct {
	*engine.Mesh
	Radius   float32
	Height   float32
	Segments int
	Color    render.Color
}

// NewCone generates a new parameterized cone mesh
func NewCone(radius, height float32, segments int, color render.Color) *Cone {
	mat := assets.GetMaterialCached("vertex_color")
	cone := &Cone{
		Mesh:     engine.NewMesh(mat),
		Radius:   radius,
		Height:   height,
		Segments: segments,
		Color:    color,
	}
	cone.Passes.Set(render.Forward)
	cone.generate()
	return cone
}

func (c *Cone) generate() {
	data := make(ColorVertices, 6*c.Segments)

	// cone
	top := vec3.New(0, c.Height, 0)
	sangle := 2 * math.Pi / float32(c.Segments)
	for i := 0; i < c.Segments; i++ {
		a1 := sangle * (float32(i) + 0.5)
		a2 := sangle * (float32(i) + 1.5)
		v1 := vec3.New(math.Cos(a1), 0, -math.Sin(a1)).Scaled(c.Radius)
		v2 := vec3.New(math.Cos(a2), 0, -math.Sin(a2)).Scaled(c.Radius)
		v1t, v2t := top.Sub(v1), top.Sub(v2)
		n := vec3.Cross(&v1t, &v2t).Normalized()

		o := 3 * i
		data[o+0] = ColorVertex{Position: v2, Normal: n, Color: c.Color}
		data[o+1] = ColorVertex{Position: top, Normal: n, Color: c.Color}
		data[o+2] = ColorVertex{Position: v1, Normal: n, Color: c.Color}
	}

	// bottom
	base := vec3.Zero
	n := vec3.New(0, -1, 0)
	for i := 0; i < c.Segments; i++ {
		a1 := sangle * (float32(i) + 0.5)
		a2 := sangle * (float32(i) + 1.5)
		v1 := vec3.New(math.Cos(a1), 0, -math.Sin(a1)).Scaled(c.Radius)
		v2 := vec3.New(math.Cos(a2), 0, -math.Sin(a2)).Scaled(c.Radius)
		o := 3 * (i + c.Segments)
		data[o+0] = ColorVertex{Position: v1, Normal: n, Color: c.Color}
		data[o+1] = ColorVertex{Position: base, Normal: n, Color: c.Color}
		data[o+2] = ColorVertex{Position: v2, Normal: n, Color: c.Color}
	}

	c.Buffer("geometry", data)
}
