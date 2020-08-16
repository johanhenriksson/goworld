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
	vao         *render.VertexArray
	vbo         *render.VertexBuffer
}

func NewImageQuad(mat *render.Material, w, h, z float32) *ImageQuad {
	return NewImageQuadAt(mat, 0, 0, w, h, z)
}

func NewImageQuadAt(mat *render.Material, x, y, w, h, z float32) *ImageQuad {
	q := &ImageQuad{
		Material:    mat,
		TopLeft:     Vertex{X: x, Y: y + h, Z: z, U: 0, V: 0},
		TopRight:    Vertex{X: x + w, Y: y + h, Z: z, U: 1, V: 0},
		BottomLeft:  Vertex{X: x, Y: y, Z: z, U: 0, V: 1},
		BottomRight: Vertex{X: x + w, Y: y, Z: z, U: 1, V: 1},
		vao:         render.CreateVertexArray(),
		vbo:         render.CreateVertexBuffer(),
	}
	q.compute()
	return q
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

func (q *ImageQuad) FlipY() {
	q.TopLeft.V = 1.0 - q.TopLeft.V
	q.TopRight.V = 1.0 - q.TopRight.V
	q.BottomLeft.V = 1.0 - q.BottomLeft.V
	q.BottomRight.V = 1.0 - q.BottomRight.V
	q.compute()
}

func (q *ImageQuad) Draw(args render.DrawArgs) {
	if q.Material != nil {
		q.Material.Use()
		q.Material.Mat4f("model", args.Transform)
		q.Material.Mat4f("viewport", args.Projection)
	}
	q.vao.DrawElements()
}
