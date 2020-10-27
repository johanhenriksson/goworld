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
	y := float32(0.001)

	o1 := vertex.C{P: vec3.New(-s, y, -s), N: vec3.UnitY, C: p.Color.Vec4()}
	x1 := vertex.C{P: vec3.New(s, y, -s), N: vec3.UnitY, C: p.Color.Vec4()}
	z1 := vertex.C{P: vec3.New(-s, y, s), N: vec3.UnitY, C: p.Color.Vec4()}
	d1 := vertex.C{P: vec3.New(s, y, s), N: vec3.UnitY, C: p.Color.Vec4()}

	o2 := vertex.C{P: vec3.New(-s, -y, -s), N: vec3.UnitYN, C: p.Color.Vec4()}
	x2 := vertex.C{P: vec3.New(s, -y, -s), N: vec3.UnitYN, C: p.Color.Vec4()}
	z2 := vertex.C{P: vec3.New(-s, -y, s), N: vec3.UnitYN, C: p.Color.Vec4()}
	d2 := vertex.C{P: vec3.New(s, -y, s), N: vec3.UnitYN, C: p.Color.Vec4()}

	data := []vertex.C{
		o1, z1, x1, x1, z1, d1,
		x2, z2, o2, d2, z2, x2,
	}
	p.Buffer(data)
}
