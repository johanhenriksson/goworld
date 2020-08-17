package geometry

import (
	"github.com/johanhenriksson/goworld/render"
)

/** Not exactly a quad anymore is it? */
type ImageQuad struct {
	Material    *render.Material
	TopLeft     Vertex
	TopRight    Vertex
	BottomLeft  Vertex
	BottomRight Vertex
	InvertY     bool
	vao         *render.VertexArray
	vbo         *render.VertexBuffer
}

func NewImageQuad(mat *render.Material, w, h float32, invert bool) *ImageQuad {
	q := &ImageQuad{
		Material:    mat,
		InvertY:     invert,
		TopLeft:     Vertex{X: 0, Y: h, Z: 0, U: 0, V: 0},
		TopRight:    Vertex{X: w, Y: h, Z: 0, U: 1, V: 0},
		BottomLeft:  Vertex{X: 0, Y: 0, Z: 0, U: 0, V: 1},
		BottomRight: Vertex{X: w, Y: 0, Z: 0, U: 1, V: 1},
		vao:         render.CreateVertexArray(),
		vbo:         render.CreateVertexBuffer(),
	}
	if q.InvertY {
		q.TopLeft.V = 1 - q.TopLeft.V
		q.TopRight.V = 1 - q.TopRight.V
		q.BottomLeft.V = 1 - q.BottomLeft.V
		q.BottomRight.V = 1 - q.BottomRight.V
	}
	q.compute()
	return q
}

func (q *ImageQuad) SetSize(w, h float32) {
	z := q.TopLeft.Z
	q.TopLeft = Vertex{X: 0, Y: h, Z: z, U: 0, V: 0}
	q.TopRight = Vertex{X: w, Y: h, Z: z, U: 1, V: 0}
	q.BottomLeft = Vertex{X: 0, Y: 0, Z: z, U: 0, V: 1}
	q.BottomRight = Vertex{X: w, Y: 0, Z: z, U: 1, V: 1}

	if q.InvertY {
		q.TopLeft.V = 1 - q.TopLeft.V
		q.TopRight.V = 1 - q.TopRight.V
		q.BottomLeft.V = 1 - q.BottomLeft.V
		q.BottomRight.V = 1 - q.BottomRight.V
	}

	q.compute()
}

func (q *ImageQuad) compute() {
	vtx := Vertices{
		q.BottomLeft, q.TopRight, q.TopLeft,
		q.BottomLeft, q.BottomRight, q.TopRight,
	}

	/* Setup VAO */
	q.vao.Length = int32(len(vtx))
	q.vao.Bind()
	q.vbo.Buffer(vtx)
	if q.Material != nil {
		q.Material.SetupVertexPointers()
	}
}

func (q *ImageQuad) Draw(args render.DrawArgs) {
	if q.Material != nil {
		q.Material.Use()
		q.Material.Mat4f("model", args.Transform)
		q.Material.Mat4f("viewport", args.Projection)
	}
	q.vao.DrawElements()
}
