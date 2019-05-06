package geometry

import (
	"github.com/johanhenriksson/goworld/render"
)

/** Not exactly a quad anymore is it? */
type ImageQuad struct {
	Material    *render.Material
	TopLeft     ImageVertex
	TopRight    ImageVertex
	BottomLeft  ImageVertex
	BottomRight ImageVertex
	vao         *render.VertexArray
	vbo         *render.VertexBuffer
}

func NewImageQuad(mat *render.Material, w, h, z float32) *ImageQuad {
	return NewImageQuadAt(mat, 0, 0, w, h, z)
}

func NewImageQuadAt(mat *render.Material, x, y, w, h, z float32) *ImageQuad {
	q := &ImageQuad{
		Material:    mat,
		TopLeft:     ImageVertex{X: x, Y: y + h, Z: z, Tx: 0, Ty: 0},
		TopRight:    ImageVertex{X: x + w, Y: y + h, Z: z, Tx: 1, Ty: 0},
		BottomLeft:  ImageVertex{X: x, Y: y, Z: z, Tx: 0, Ty: 1},
		BottomRight: ImageVertex{X: x + w, Y: y, Z: z, Tx: 1, Ty: 1},
		vao:         render.CreateVertexArray(),
		vbo:         render.CreateVertexBuffer(),
	}
	q.compute()
	return q
}

func (q *ImageQuad) compute() {
	vtx := ImageVertices{
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
	q.TopLeft.Ty = 1.0 - q.TopLeft.Ty
	q.TopRight.Ty = 1.0 - q.TopRight.Ty
	q.BottomLeft.Ty = 1.0 - q.BottomLeft.Ty
	q.BottomRight.Ty = 1.0 - q.BottomRight.Ty
	q.compute()
}

func (q *ImageQuad) Draw(args render.DrawArgs) {
	if q.Material != nil {
		q.Material.Use()
		q.Material.Shader.Matrix4f("model", &args.Transform[0])
		q.Material.Shader.Matrix4f("viewport", &args.Projection[0])
	}
	q.vao.Draw()
}
