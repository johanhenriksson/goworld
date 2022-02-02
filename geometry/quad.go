package geometry

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vertex_array"
)

// Quad with support for centered backgrounds
type Quad struct {
	Width    float32
	Height   float32
	Invert   bool
	Depth    bool
	Material material.T

	size vec2.T
	vao  vertex_array.T
}

func NewQuad(mat material.T, size vec2.T) *Quad {
	q := &Quad{
		Material: mat,
		Invert:   false,
		Depth:    false,

		size: size,
		vao:  gl_vertex_array.New(render.Triangles),
	}
	q.compute()
	return q
}

func (q *Quad) Size() vec2.T { return q.size }
func (q *Quad) SetSize(size vec2.T) {
	q.size = size
	q.compute()
}

func (q *Quad) compute() {
	w, h := q.size.X, q.size.Y

	TopLeft := vertex.T{
		P: vec3.New(0, 0, -1),
		T: vec2.New(0, 1),
	}
	TopRight := vertex.T{
		P: vec3.New(w, 0, -1),
		T: vec2.New(1, 1),
	}
	BottomLeft := vertex.T{
		P: vec3.New(0, h, -1),
		T: vec2.New(0, 0),
	}
	BottomRight := vertex.T{
		P: vec3.New(w, h, -1),
		T: vec2.New(1, 0),
	}

	vtx := []vertex.T{
		TopLeft,
		TopRight,
		BottomLeft,

		TopRight,
		BottomRight,
		BottomLeft,
	}

	// draw vertex array
	q.vao.Bind()
	ptr := q.Material.VertexPointers(vtx)
	q.vao.BufferTo(ptr, vtx)
}

func (q *Quad) Draw(args render.Args) {
	q.Material.Use()
	q.Material.Mat4("model", args.Transform)
	q.Material.Mat4("viewport", args.VP)
	q.Material.Bool("invert", q.Invert)
	q.Material.Bool("depth", q.Depth)
	q.vao.Draw()
}
