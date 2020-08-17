package geometry

import (
	"math"

	"github.com/johanhenriksson/goworld/render"
)

type Quad struct {
	Width    float32
	Height   float32
	Color    render.Color
	Material *render.Material

	segments int
	border   float32
	vao      *render.VertexArray
	vbo      *render.VertexBuffer
}

func NewQuad(mat *render.Material, color render.Color, w, h float32) *Quad {
	q := &Quad{
		Material: mat,
		Width:    w,
		Height:   h,
		Color:    color,
		segments: 5,
		border:   0,

		vao: render.CreateVertexArray(),
		vbo: render.CreateVertexBuffer(),
	}
	q.vao.Bind()
	q.vbo.Bind()
	mat.SetupVertexPointers()
	q.compute()
	return q
}

func (q *Quad) BorderWidth() float32 { return q.border }
func (q *Quad) SetBorderWidth(width float32) {
	q.border = width
	q.compute()
}

func (q *Quad) appendCorner(vtx *ColorVertices, origin ColorVertex, n int, r, offset float32) {
	if n == 0 {
		/* TODO: Square corner */
	} else {
		/* Rounded corner */
		v := (math.Pi / 2.0) / float64(n)
		var prev ColorVertex
		for i := 0; i <= n; i++ {
			p := ColorVertex{
				X:     origin.X + r*float32(math.Cos(float64(offset)+float64(i)*v)),
				Y:     origin.Y + r*float32(math.Sin(float64(offset)+float64(i)*v)),
				Z:     origin.Z,
				Color: origin.Color,
			}

			if i > 0 {
				*vtx = append(*vtx, origin, prev, p)
			}

			prev = p
		}
	}
}

func (q *Quad) SetSize(w, h float32) {
	q.Width = w
	q.Height = h
	q.compute()
}

func (q *Quad) compute() {
	b := q.border
	TopLeft := ColorVertex{X: b, Y: q.Height - b, Z: 0, Color: q.Color}
	TopRight := ColorVertex{X: q.Width - b, Y: q.Height - b, Z: 0, Color: q.Color}
	BottomLeft := ColorVertex{X: b, Y: b, Z: 0, Color: q.Color}
	BottomRight := ColorVertex{X: q.Width - b, Y: b, Z: 0, Color: q.Color}

	vtx := ColorVertices{
		BottomLeft, TopRight, TopLeft,
		BottomLeft, BottomRight, TopRight,
	}

	/* If we have a positive border width, tesselate border */
	if b > 0.0 {
		q.appendCorner(&vtx, TopRight, q.segments, q.border, 0.0)
		q.appendCorner(&vtx, TopLeft, q.segments, q.border, math.Pi/2.0)
		q.appendCorner(&vtx, BottomLeft, q.segments, q.border, math.Pi)
		q.appendCorner(&vtx, BottomRight, q.segments, q.border, 3.0*math.Pi/2.0)

		/* Top Border Box */
		topTopLeft := TopLeft
		topTopLeft.Y += b
		topTopRight := TopRight
		topTopRight.Y += b
		vtx = append(vtx, TopLeft, topTopRight, topTopLeft,
			TopLeft, TopRight, topTopRight)

		/* Bottom border box */
		bottomBottomLeft := BottomLeft
		bottomBottomLeft.Y -= b
		bottomBottomRight := BottomRight
		bottomBottomRight.Y -= b
		vtx = append(vtx, bottomBottomLeft, BottomRight, BottomLeft,
			bottomBottomLeft, bottomBottomRight, BottomRight)

		/* Right border box */
		rightTopRight := TopRight
		rightTopRight.X += b
		rightBottomRight := BottomRight
		rightBottomRight.X += b
		vtx = append(vtx, BottomRight, rightTopRight, TopRight,
			BottomRight, rightBottomRight, rightTopRight)

		/* Left border box */
		leftTopLeft := TopLeft
		leftTopLeft.X -= b
		leftBottomLeft := BottomLeft
		leftBottomLeft.X -= b
		vtx = append(vtx, leftBottomLeft, TopLeft, leftTopLeft,
			leftBottomLeft, BottomLeft, TopLeft)
	}

	/* Setup VAO */
	q.vao.Length = int32(len(vtx))
	q.vao.Bind()
	q.vbo.Buffer(vtx)
}

func (q *Quad) Draw(args render.DrawArgs) {
	q.Material.Use()
	q.Material.Mat4f("model", args.Transform)
	q.Material.Mat4f("viewport", args.Projection)
	q.vao.DrawElements()
}

func (q *Quad) SetColor(color render.Color) {
	q.Color = color
	q.compute()
}
