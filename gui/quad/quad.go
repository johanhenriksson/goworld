package quad

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	Size() vec2.T
	Position() vec2.T
	Color() color.T
	Update(Props)
	Mesh() vertex.Mesh
}

type Props struct {
	Size     vec2.T
	Position vec2.T
	Color    color.T
	UVs      UV
}

type quad struct {
	props Props
	mesh  vertex.MutableMesh[vertex.UI, uint16]
}

func New(props Props) T {
	q := &quad{
		props: props,
		mesh:  vertex.NewTriangles[vertex.UI, uint16]("quad", nil, nil),
	}
	q.compute()
	return q
}

func (q *quad) Position() vec2.T { return q.props.Position }
func (q *quad) Size() vec2.T     { return q.props.Size }
func (q *quad) Color() color.T   { return q.props.Color }

func (q *quad) Update(props Props) {
	q.props = props
	q.compute()
}

func (q *quad) compute() {
	x, y := q.props.Position.X, q.props.Position.Y
	w, h := q.props.Size.X, q.props.Size.Y
	ax, ay := q.props.UVs.A.X, q.props.UVs.A.Y
	bx, by := q.props.UVs.B.X, q.props.UVs.B.Y

	TopLeft := vertex.UI{
		P: vec3.New(x, y, -1),
		T: vec2.New(ax, ay),
		C: q.props.Color,
	}
	TopRight := vertex.UI{
		P: vec3.New(x+w, y, -1),
		T: vec2.New(bx, ay),
		C: q.props.Color,
	}
	BottomLeft := vertex.UI{
		P: vec3.New(x, y+h, -1),
		T: vec2.New(ax, by),
		C: q.props.Color,
	}
	BottomRight := vertex.UI{
		P: vec3.New(x+w, y+h, -1),
		T: vec2.New(bx, by),
		C: q.props.Color,
	}

	vtx := []vertex.UI{
		TopLeft,
		TopRight,
		BottomLeft,
		BottomRight,
	}
	idx := []uint16{
		0, 1, 2,
		1, 3, 2,
	}
	q.mesh.Update(vtx, idx)
}

func (q *quad) Mesh() vertex.Mesh {
	return q.mesh
}
