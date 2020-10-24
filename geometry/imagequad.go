package geometry

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
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
		vao:      render.CreateVertexArray(render.Triangles),
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
	TopLeft := vertex.T{P: vec3.New(0, q.Height, 0), T: vec2.New(0, 0)}
	TopRight := vertex.T{P: vec3.New(q.Width, q.Height, 0), T: vec2.New( q.U, 0)}
	BottomLeft := vertex.T{P: vec3.New(0, 0, 0), T: vec2.New(0, q.V)}
	BottomRight := vertex.T{P: vec3.New(q.Width, 0, 0), T: vec2.New(q.U, q.V)}

	if q.InvertY {
		TopLeft.T.Y = 1 - TopLeft.T.Y
		TopRight.T.Y = 1 - TopRight.T.Y
		BottomLeft.T.Y = 1 - BottomLeft.T.Y
		BottomRight.T.Y = 1 - BottomRight.T.Y
	}

	vtx := []vertex.T{
		BottomLeft, TopRight, TopLeft,
		BottomLeft, BottomRight, TopRight,
	}

	ptr := q.Material.VertexPointers(vtx)

	q.vao.Bind()
	q.vao.BufferTo(ptr, vtx)
}

func (q *ImageQuad) Draw(args engine.DrawArgs) {
	if q.Material != nil {
		q.Material.Use()
		q.Material.Mat4("model", &args.Transform)
		q.Material.Mat4("viewport", &args.Projection)
	}
	q.vao.Draw()
}
