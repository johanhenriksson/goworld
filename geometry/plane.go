package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Plane is a colored, one segment, one-sided 3D plane
type Plane struct {
	*engine.Mesh
	Size  float32
	Color render.Color
}

// NewPlane creates a new 3D plane of a given size and color.
func NewPlane(size float32, color render.Color) *Plane {
	mat := assets.GetMaterialCached("vertex_color")
	plane := &Plane{
		Mesh:  engine.NewMesh(mat),
		Size:  size,
		Color: color,
	}
	plane.Pass = engine.DrawForward
	plane.generate()
	return plane
}

func (p *Plane) generate() {
	s := p.Size / 2
	o := ColorVertex{
		Position: vec3.New(-s, 0, -s),
		Normal:   vec3.UnitY,
		Color:    p.Color,
	}
	x := ColorVertex{
		Position: vec3.New(s, 0, -s),
		Normal:   vec3.UnitY,
		Color:    p.Color,
	}
	z := ColorVertex{
		Position: vec3.New(-s, 0, s),
		Normal:   vec3.UnitY,
		Color:    p.Color,
	}
	d := ColorVertex{
		Position: vec3.New(s, 0, s),
		Normal:   vec3.UnitY,
		Color:    p.Color,
	}
	data := ColorVertices{
		o, z, x, x, z, d,
		x, z, o, d, z, x,
	}
	p.Buffer("geometry", data)
}
