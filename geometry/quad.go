package geometry

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

// Rect with support for borders and rounded corners.
type Rect struct {
	Width    float32
	Height   float32
	Material *render.Material

	segments int
	border   float32
	vao      *render.VertexArray
}

func NewQuad(mat *render.Material, size vec2.T) *Rect {
	q := &Rect{
		Material: mat,
		Width:    size.X,
		Height:   size.Y,
		segments: 5,
		border:   0,

		vao: render.CreateVertexArray(render.Triangles, "geometry"),
	}
	mat.SetupVertexPointers()
	q.compute()
	return q
}

func (q *Rect) BorderWidth() float32 { return q.border }
func (q *Rect) SetBorderWidth(width float32) {
	q.border = width
	q.compute()
}

func (q *Rect) appendCorner(vtx *Vertices, origin Vertex, offset float32) {
	r := q.border
	n := q.segments

	bw, bh := float32(0), float32(0)
	if tex := q.texture(); tex != nil {
		b := float32(tex.Border)
		bw, bh = b/float32(tex.Width), b/float32(tex.Height)
	}
	bw, bh = float32(128.0/1024.0), float32(128.0/1024.0)

	if n == 0 {
		first := Vertex{
			X: origin.X + r*math.Cos(offset),
			Y: origin.Y + r*math.Sin(offset),
			Z: 0,
			U: origin.U + bw*math.Cos(offset),
			V: origin.V + bh*math.Sin(offset),
		}
		corner := Vertex{
			X: origin.X + r*math.Cos(offset+math.Pi/4)*math.Sqrt2,
			Y: origin.Y + r*math.Sin(offset+math.Pi/4)*math.Sqrt2,
			Z: 0,
			U: origin.U + bw*math.Cos(offset+math.Pi/4)*math.Sqrt2,
			V: origin.V + bh*math.Sin(offset+math.Pi/4)*math.Sqrt2,
		}
		second := Vertex{
			X: origin.X + r*math.Cos(offset+math.Pi/2),
			Y: origin.Y + r*math.Sin(offset+math.Pi/2),
			Z: 0,
			U: origin.U + bw*math.Cos(offset+math.Pi/2),
			V: origin.V + bh*math.Sin(offset+math.Pi/2),
		}
		*vtx = append(*vtx, first, origin, corner)
		*vtx = append(*vtx, corner, origin, second)
	} else {
		/* Rounded corner */

		v := (math.Pi / 2.0) / float32(n)
		var prev Vertex
		for i := 0; i <= n; i++ {
			x := math.Cos(offset + float32(i)*v)
			y := math.Sin(offset + float32(i)*v)
			p := Vertex{
				X: origin.X + r*x,
				Y: origin.Y + r*y,
				Z: origin.Z,
				U: origin.U + bw*x,
				V: origin.V + bh*y,
			}

			if i > 0 {
				*vtx = append(*vtx, origin, prev, p)
			}

			prev = p
		}
	}
}

func (q *Rect) SetSize(size vec2.T) {
	q.Width = size.X
	q.Height = size.Y
	q.compute()
}

func (q *Rect) texture() *render.Texture {
	for _, tex := range q.Material.Textures {
		return tex
	}
	return nil
}

func (q *Rect) compute() {
	b := q.border

	w, h := q.Width, q.Height
	bw, bh := float32(0), float32(0)
	if tex := q.texture(); tex != nil {
		tb := float32(tex.Border)
		bw, bh = tb/float32(tex.Width), tb/float32(tex.Height)
	}

	bw, bh = float32(128.0/1024.0), float32(128.0/1024.0)

	TopLeft := Vertex{X: b, Y: h - b, Z: 0, U: bw, V: 1 - bh}
	TopRight := Vertex{X: w - b, Y: h - b, Z: 0, U: 1 - bw, V: 1 - bh}
	BottomLeft := Vertex{X: b, Y: b, Z: 0, U: bw, V: bh}
	BottomRight := Vertex{X: w - b, Y: b, Z: 0, U: 1 - bw, V: bh}

	vtx := Vertices{
		BottomLeft, TopRight, TopLeft,
		BottomLeft, BottomRight, TopRight,
	}

	/* If we have a positive border width, tesselate border */
	if b > 0.0 {
		q.appendCorner(&vtx, TopRight, 0.0)
		q.appendCorner(&vtx, TopLeft, math.Pi/2.0)
		q.appendCorner(&vtx, BottomLeft, math.Pi)
		q.appendCorner(&vtx, BottomRight, 3.0*math.Pi/2.0)

		/* Top Border Box */
		topTopLeft := TopLeft
		topTopLeft.Y += b
		topTopLeft.V = 1
		topTopRight := TopRight
		topTopRight.Y += b
		topTopRight.V = 1
		vtx = append(vtx, TopLeft, topTopRight, topTopLeft,
			TopLeft, TopRight, topTopRight)

		/* Bottom border box */
		bottomBottomLeft := BottomLeft
		bottomBottomLeft.Y -= b
		bottomBottomLeft.V = 0
		bottomBottomRight := BottomRight
		bottomBottomRight.Y -= b
		bottomBottomRight.V = 0
		vtx = append(vtx, bottomBottomLeft, BottomRight, BottomLeft,
			bottomBottomLeft, bottomBottomRight, BottomRight)

		/* Right border box */
		rightTopRight := TopRight
		rightTopRight.X += b
		rightTopRight.U = 1
		rightBottomRight := BottomRight
		rightBottomRight.X += b
		rightBottomRight.U = 1
		vtx = append(vtx, BottomRight, rightTopRight, TopRight,
			BottomRight, rightBottomRight, rightTopRight)

		/* Left border box */
		leftTopLeft := TopLeft
		leftTopLeft.X -= b
		leftTopLeft.U = 0
		leftBottomLeft := BottomLeft
		leftBottomLeft.X -= b
		leftBottomLeft.U = 0
		vtx = append(vtx, leftBottomLeft, TopLeft, leftTopLeft,
			leftBottomLeft, BottomLeft, TopLeft)
	}

	/* Setup VAO */
	q.vao.Bind()
	q.vao.Buffer("geometry", vtx)
}

func (q *Rect) Draw(args engine.DrawArgs) {
	q.Material.Use()
	q.Material.Mat4("model", &args.Transform)
	q.Material.Mat4("viewport", &args.Projection)
	q.vao.Draw()
}
