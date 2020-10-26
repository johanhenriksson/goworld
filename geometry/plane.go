package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Plane is a colored, one segment, one-sided 3D plane
type Plane struct {
	*engine.Mesh
	Size  float32
	Color render.Color
}

type MeshOptions struct {
	Material       string
	SharedMaterial bool
	Pass           render.Pass
}

// NewPlane creates a new 3D plane of a given size and color.
func NewPlane(size float32, color render.Color) *Plane {
	mat := assets.GetMaterialShared("color.f")
	plane := &Plane{
		Mesh:  engine.NewMesh("Plane", mat),
		Size:  size,
		Color: color,
	}
	plane.Pass = render.Forward
	plane.generate()
	return plane
}

func (p *Plane) generate() {
	s := p.Size / 2
	o := vertex.C{
		P: vec3.New(-s, 0, -s),
		N: vec3.UnitY,
		C: p.Color.Vec4(),
	}
	x := vertex.C{
		P: vec3.New(s, 0, -s),
		N: vec3.UnitY,
		C: p.Color.Vec4(),
	}
	z := vertex.C{
		P: vec3.New(-s, 0, s),
		N: vec3.UnitY,
		C: p.Color.Vec4(),
	}
	d := vertex.C{
		P: vec3.New(s, 0, s),
		N: vec3.UnitY,
		C: p.Color.Vec4(),
	}
	data := []vertex.C{
		o, z, x, x, z, d,
		x, z, o, d, z, x,
	}
	p.Buffer(data)
}
