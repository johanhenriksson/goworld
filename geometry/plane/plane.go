package plane

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Plane is a colored, one segment, one-sided 3D plane
type T struct {
	*engine.Mesh
	Args
}

type Args struct {
	Size  float32
	Color render.Color
}

// New creates a new 3D plane of a given size and color.
func New(args Args) *T {
	parent := object.New("Plane")
	return Attach(parent, args)
}

func Attach(parent *object.T, args Args) *T {
	mat := assets.GetMaterialShared("color.f")
	plane := &T{
		Mesh: engine.NewMesh("Plane", mat),
		Args: args,
	}
	plane.Pass = render.Forward
	plane.generate()
	parent.Attach(plane)
	return plane
}

func (p *T) generate() {
	s := p.Size / 2
	y := float32(0.001)
	c := p.Color.Vec4()

	o1 := vertex.C{P: vec3.New(-s, y, -s), N: vec3.UnitY, C: c}
	x1 := vertex.C{P: vec3.New(s, y, -s), N: vec3.UnitY, C: c}
	z1 := vertex.C{P: vec3.New(-s, y, s), N: vec3.UnitY, C: c}
	d1 := vertex.C{P: vec3.New(s, y, s), N: vec3.UnitY, C: c}

	o2 := vertex.C{P: vec3.New(-s, -y, -s), N: vec3.UnitYN, C: c}
	x2 := vertex.C{P: vec3.New(s, -y, -s), N: vec3.UnitYN, C: c}
	z2 := vertex.C{P: vec3.New(-s, -y, s), N: vec3.UnitYN, C: c}
	d2 := vertex.C{P: vec3.New(s, -y, s), N: vec3.UnitYN, C: c}

	data := []vertex.C{
		o1, z1, x1, x1, z1, d1,
		x2, z2, o2, d2, z2, x2,
	}
	p.Buffer(data)
}

func (p *T) DrawForward(args engine.DrawArgs) {
	p.Mesh.DrawForward(args)
}
