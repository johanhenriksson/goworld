package quad

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	Size() vec2.T
	Position() vec2.T
	Material() material.T
	Color() color.T
	Update(Props)
	Draw(render.Args)
	Destroy()
}

type Props struct {
	Size     vec2.T
	Position vec2.T
	Color    color.T
	UVs      UV
}

type quad struct {
	props Props
	vao   vertex.Array
	mat   material.T
}

func New(mat material.T, props Props) T {
	q := &quad{
		props: props,
		mat:   mat,
		vao:   gl_vertex_array.New(vertex.Triangles),
	}
	q.compute()
	return q
}

func (q *quad) Material() material.T { return q.mat }
func (q *quad) Position() vec2.T     { return q.props.Position }
func (q *quad) Size() vec2.T         { return q.props.Size }
func (q *quad) Color() color.T       { return q.props.Color }

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

		TopRight,
		BottomRight,
		BottomLeft,
	}

	ptrs := vertex.ParsePointers(vtx)
	ptrs.Bind(q.mat)

	q.vao.SetPointers(ptrs)
	q.vao.SetIndexSize(0)
	q.vao.SetElements(len(vtx))

	q.vao.Buffer("vertex", vtx)
}

func (q *quad) Draw(args render.Args) {
	q.mat.Use()
	q.mat.Mat4("model", args.Transform)
	q.mat.Mat4("viewport", args.VP)
	q.vao.Draw()
}

func (q *quad) Destroy() {
	q.vao.Delete()
}
