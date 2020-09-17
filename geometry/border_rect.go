package geometry

/*
import (
	"math"

	"github.com/johanhenriksson/goworld/render"
)

type BorderRect struct {
	Width    float32
	Height   float32
	Material *render.Material

	border float32
	vao    *render.VertexArray
	vbo    *render.VertexBuffer
	vbi    *render.VertexBuffer
}

func NewBorderRect(mat *render.Material, w, h float32) *Quad {
	q := &Quad{
		Material: mat,
		Width:    w,
		Height:   h,
		border:   0,

		vao: render.CreateVertexArray(),
		vbo: render.CreateVertexBuffer(),
		vbi: render.CreateIndexBuffer(),
	}
	q.vao.Bind()
	q.vbo.Bind()
	q.vbi.Bind()
	mat.SetupVertexPointers()
	q.compute()
	return q
}

func (q *BorderRect) BorderWidth() float32 { return q.border }
func (q *BorderRect) SetBorderWidth(width float32) {
	q.border = width
	q.compute()
}

func (q *BorderRect) SetSize(w, h float32) {
	q.Width = w
	q.Height = h
	q.compute()
}

func (q *BorderRect) compute() {
	// pick texture dimensions from first texture
	tw, th := 128.0, 128.0
	for _, t := q.Material.Textures {
		tw = float32(t.Width)
		th = float32(t.Height)
		break
	}

	w, h, b := q.Width, q.Height, q.border
	bw, bh := q.border/tw, q.border/th

	vtx := ImageVertices{
		ImageVertex{X: 0, Y: h, Z: 0, Tx: 0, Ty: 1},
		ImageVertex{X: b, Y: h, Z: 0, Tx: bw, Ty: 1},
		ImageVertex{X: w - b, Y: h, Z: 0, Tx: 1 - bw, Ty: 1},
		ImageVertex{X: w, Y: h, Z: 0, Tx: 1, Ty: 1},

		ImageVertex{X: 0, Y: h - b, Z: 0, Tx: 0, Ty: 1},
		ImageVertex{X: b, Y: h - b, Z: 0, Tx: bw, Ty: 1},
		ImageVertex{X: w - b, Y: h - b, Z: 0, Tx: 1 - bw, Ty: 1},
		ImageVertex{X: w, Y: h - b, Z: 0, Tx: 1, Ty: 1},

		ImageVertex{X: 0, Y: b, Z: 0, Tx: 0, Ty: 1},
		ImageVertex{X: b, Y: b, Z: 0, Tx: bw, Ty: 1},
		ImageVertex{X: w - b, Y: b, Z: 0, Tx: 1 - bw, Ty: 1},
		ImageVertex{X: w, Y: b, Z: 0, Tx: 1, Ty: 1},

		ImageVertex{X: 0, Y: 0, Z: 0, Tx: 0, Ty: 1},
		ImageVertex{X: b, Y: 0, Z: 0, Tx: bw, Ty: 1},
		ImageVertex{X: w - b, Y: 0, Z: 0, Tx: 1 - bw, Ty: 1},
		ImageVertex{X: w, Y: 0, Z: 0, Tx: 1, Ty: 1},
	}

	idx := []uint{
		0, 1, 4, 1, 4, 5,
		1, 5, 2, 2, 5, 6,
		2, 6, 3, 3, 6, 7,
		4, 8, 5, 5, 8, 9,
		5, 8, 6, 6, 9, 10,
		6, 10, 7, 7, 10, 11,
		8, 12, 9, 9, 12, 13,
		9, 13, 10, 10, 13, 14,
		10, 14, 11, 11, 14, 15,
	}

	q.vao.Length = int32(len(vtx))
	q.vao.Bind()
	q.vbo.Buffer(vtx)
	q.vbi.Buffer(idx)
}

func (q *BorderRect) Draw(args render.DrawArgs) {
	q.Material.Use()
	q.Material.Mat4f("model", args.Transform)
	q.Material.Mat4f("viewport", args.Projection)
	q.vao.DrawIndexed()
}

func (q *BorderRect) SetColor(color render.Color) {
	q.Color = color
	q.compute()
}
*/
