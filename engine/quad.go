package engine

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Quad is a fullscreen quad used for render passes
type Quad struct {
	*Mesh
}

// NewQuad creates a new quad with a given material
func NewQuad(mat *render.Material) *Quad {
	q := &Quad{
		Mesh: NewMesh("screen_quad", mat),
	}

	vtx := []vertex.T{
		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
		{P: vec3.New(-1, 1, 0), T: vec2.New(0, 1)},

		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, -1, 0), T: vec2.New(1, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
	}

	ptr := mat.VertexPointers(vtx)
	q.vao.BufferTo(ptr, vtx)

	return q
}

func (q *Quad) Draw() {
	q.vao.Bind()
	if q.material != nil {
		q.material.Use()
	}
	q.vao.Draw()
}
