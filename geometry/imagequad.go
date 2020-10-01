package geometry

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type ImageQuad struct {
	Material *render.Material
	Width    float32
	Height   float32
	U        float32
	V        float32
	InvertY  bool
	vao      *render.VertexArray
}

func NewImageQuad(mat *render.Material, size vec2.T, invert bool) *ImageQuad {
	q := &ImageQuad{
		Material: mat,
		InvertY:  invert,
		Width:    size.X,
		Height:   size.Y,
		U:        1,
		V:        1,
		vao:      render.CreateVertexArray(render.Triangles, "geometry"),
	}
	q.compute()
	return q
}

func (q *ImageQuad) SetSize(size vec2.T) {
	q.Width = size.X
	q.Height = size.Y
	q.compute()
}

func (q *ImageQuad) SetUV(u, v float32) {
	q.U = u
	q.V = v
	q.compute()
}

func (q *ImageQuad) compute() {
	TopLeft := Vertex{X: 0, Y: q.Height, Z: 0, U: 0, V: 0}
	TopRight := Vertex{X: q.Width, Y: q.Height, Z: 0, U: q.U, V: 0}
	BottomLeft := Vertex{X: 0, Y: 0, Z: 0, U: 0, V: q.V}
	BottomRight := Vertex{X: q.Width, Y: 0, Z: 0, U: q.U, V: q.V}

	if q.InvertY {
		TopLeft.V = 1 - TopLeft.V
		TopRight.V = 1 - TopRight.V
		BottomLeft.V = 1 - BottomLeft.V
		BottomRight.V = 1 - BottomRight.V
	}

	vtx := Vertices{
		BottomLeft, TopRight, TopLeft,
		BottomLeft, BottomRight, TopRight,
	}

	/* Setup VAO */
	q.vao.Bind()
	q.vao.Buffer("geometry", vtx)
	if q.Material != nil {
		q.Material.SetupVertexPointers()
	}
}

func (q *ImageQuad) Draw(args render.DrawArgs) {
	if q.Material != nil {
		q.Material.Use()
		q.Material.Mat4f("model", &args.Transform)
		q.Material.Mat4f("viewport", &args.Projection)
	}
	q.vao.Draw()
}
