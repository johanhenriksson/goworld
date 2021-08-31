package cone

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// A Cone is a forward rendered colored cone mesh
type T struct {
	*engine.Mesh
	Args
}

type Args struct {
	Radius   float32
	Height   float32
	Segments int
	Color    render.Color
}

// NewObject creates a new Cone attached to a Game Object
func NewObject(args Args) *T {
	parent := object.New("Cone")
	return Attach(parent, args)
}

func Builder(out **T, args Args) *object.Builder {
	var tmp *T = nil
	if out == nil {
		out = &tmp
	}
	b := object.Build("Cone")
	*out = New(args)
	return b.Attach(*out)
}

func New(args Args) *T {
	mat := assets.GetMaterialShared("color.f")
	cone := &T{
		Mesh: engine.NewMesh(mat),
		Args: args,
	}
	cone.Pass = render.Forward
	cone.generate()
	return cone
}

func Attach(parent object.T, args Args) *T {
	plane := New(args)
	parent.Attach(plane)
	return plane
}

func (c *T) generate() {
	data := make([]vertex.C, 6*c.Segments)

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
		data[o+0] = vertex.C{P: v2, N: n, C: c.Color.Vec4()}
		data[o+1] = vertex.C{P: top, N: n, C: c.Color.Vec4()}
		data[o+2] = vertex.C{P: v1, N: n, C: c.Color.Vec4()}
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
		data[o+0] = vertex.C{P: v1, N: n, C: c.Color.Vec4()}
		data[o+1] = vertex.C{P: base, N: n, C: c.Color.Vec4()}
		data[o+2] = vertex.C{P: v2, N: n, C: c.Color.Vec4()}
	}

	c.Buffer(data)
}
