package geometry

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vertex_array"
)

// Rect with support for borders and rounded corners.
type Rect struct {
	Width    float32
	Height   float32
	Invert   bool
	Depth    bool
	Material material.T

	segments int
	border   float32
	vao      vertex_array.T
}

func NewRect(mat material.T, size vec2.T) *Rect {
	q := &Rect{
		Material: mat,
		Width:    size.X,
		Height:   size.Y,
		segments: 5,
		border:   0,
		Invert:   false,
		Depth:    false,

		vao: gl_vertex_array.New(render.Triangles),
	}
	q.compute()
	return q
}

func (q *Rect) BorderWidth() float32 { return q.border }
func (q *Rect) SetBorderWidth(width float32) {
	q.border = width
	q.compute()
}

func (q *Rect) appendCorner(vtx *[]vertex.T, origin vertex.T, offset float32) {
	r := q.border
	n := q.segments

	bw, bh := float32(0), float32(0)
	if tex := q.texture(); tex != nil {
		b := float32(0) // float32(tex.Border)
		bw, bh = b/float32(tex.Width()), b/float32(tex.Height())
	}
	bw, bh = float32(128.0/1024.0), float32(128.0/1024.0)

	if n == 0 {
		first := vertex.T{
			P: origin.P.Add(vec3.New(
				r*math.Cos(offset),
				r*math.Sin(offset),
				0)),
			T: origin.T.Add(vec2.New(
				bw*math.Cos(offset),
				bh*math.Sin(offset))),
		}
		corner := vertex.T{
			P: origin.P.Add(vec3.New(
				r*math.Cos(offset+math.Pi/4)*math.Sqrt2,
				r*math.Sin(offset+math.Pi/4)*math.Sqrt2,
				0)),
			T: origin.T.Add(vec2.New(
				bw*math.Cos(offset+math.Pi/4)*math.Sqrt2,
				bh*math.Sin(offset+math.Pi/4)*math.Sqrt2)),
		}
		second := vertex.T{
			P: origin.P.Add(vec3.New(
				r*math.Cos(offset+math.Pi/2),
				r*math.Sin(offset+math.Pi/2),
				0)),
			T: origin.T.Add(vec2.New(
				bw*math.Cos(offset+math.Pi/2),
				bh*math.Sin(offset+math.Pi/2))),
		}
		*vtx = append(*vtx, first, origin, corner)
		*vtx = append(*vtx, corner, origin, second)
	} else {
		/* Rounded corner */

		v := (math.Pi / 2.0) / float32(n)
		var prev vertex.T
		for i := 0; i <= n; i++ {
			x := math.Cos(offset + float32(i)*v)
			y := math.Sin(offset + float32(i)*v)
			p := vertex.T{
				P: origin.P.Add(vec3.New(r*x, r*y, 0)),
				T: origin.T.Add(vec2.New(bw*x, bh*y)),
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

func (q *Rect) texture() texture.T {
	return q.Material.TextureSlot(0)
}

func (q *Rect) compute() {
	b := q.border

	w, h := q.Width, q.Height
	bw, bh := float32(0), float32(0)
	if tex := q.texture(); tex != nil {
		tb := float32(0.0) // float32(tex.Border)
		bw, bh = tb/float32(tex.Width()), tb/float32(tex.Height())
	}

	// bw, bh = float32(128.0/1024.0), float32(128.0/1024.0)

	TopLeft := vertex.T{P: vec3.New(b, h-b, 0), T: vec2.New(bw, 1-bh)}
	TopRight := vertex.T{P: vec3.New(w-b, h-b, 0), T: vec2.New(1-bw, 1-bh)}
	BottomLeft := vertex.T{P: vec3.New(b, b, 0), T: vec2.New(bw, bh)}
	BottomRight := vertex.T{P: vec3.New(w-b, b, 0), T: vec2.New(1-bw, bh)}

	vtx := []vertex.T{
		BottomLeft, TopRight, TopLeft,
		BottomLeft, BottomRight, TopRight,
	}

	// if we have a positive border width, tesselate border
	if b > 0.0 {
		q.appendCorner(&vtx, TopRight, 0.0)
		q.appendCorner(&vtx, TopLeft, math.Pi/2.0)
		q.appendCorner(&vtx, BottomLeft, math.Pi)
		q.appendCorner(&vtx, BottomRight, 3.0*math.Pi/2.0)

		// top border box
		topTopLeft := TopLeft
		topTopLeft.P.Y += b
		topTopLeft.T.Y = 1
		topTopRight := TopRight
		topTopRight.P.Y += b
		topTopRight.T.Y = 1
		vtx = append(vtx, TopLeft, topTopRight, topTopLeft,
			TopLeft, TopRight, topTopRight)

		// bottom border box
		bottomBottomLeft := BottomLeft
		bottomBottomLeft.P.Y -= b
		bottomBottomLeft.T.Y = 0
		bottomBottomRight := BottomRight
		bottomBottomRight.P.Y -= b
		bottomBottomRight.T.Y = 0
		vtx = append(vtx, bottomBottomLeft, BottomRight, BottomLeft,
			bottomBottomLeft, bottomBottomRight, BottomRight)

		// right border box
		rightTopRight := TopRight
		rightTopRight.P.X += b
		rightTopRight.T.X = 1
		rightBottomRight := BottomRight
		rightBottomRight.P.X += b
		rightBottomRight.T.X = 1
		vtx = append(vtx, BottomRight, rightTopRight, TopRight,
			BottomRight, rightBottomRight, rightTopRight)

		// left border box
		leftTopLeft := TopLeft
		leftTopLeft.P.X -= b
		leftTopLeft.T.X = 0
		leftBottomLeft := BottomLeft
		leftBottomLeft.P.X -= b
		leftBottomLeft.T.X = 0
		vtx = append(vtx, leftBottomLeft, TopLeft, leftTopLeft,
			leftBottomLeft, BottomLeft, TopLeft)
	}

	ptrs := vertex.ParsePointers(vtx)
	ptrs.Bind(q.Material)

	q.vao.Buffer("vertex", vtx)
	q.vao.SetPointers(ptrs)
}

func (q *Rect) Draw(args render.Args) {
	q.Material.Use()
	q.Material.Mat4("model", args.Transform)
	q.Material.Mat4("viewport", args.VP)
	q.Material.Bool("invert", q.Invert)
	q.Material.Bool("depth", q.Depth)
	q.vao.Draw()
}
