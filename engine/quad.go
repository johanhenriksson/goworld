package engine

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vertex_array"
)

// Quad is a fullscreen quad used for render passes
type Quad struct {
	vao vertex_array.T
}

// NewQuad creates a new quad with a given material
func NewQuad(shader shader.T) *Quad {
	q := &Quad{
		vao: gl_vertex_array.New(render.Triangles),
	}

	vtx := []vertex.T{
		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
		{P: vec3.New(-1, 1, 0), T: vec2.New(0, 1)},

		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, -1, 0), T: vec2.New(1, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
	}

	ptr := shader.VertexPointers(vtx)
	q.vao.BufferTo(ptr, vtx)

	return q
}

func (q *Quad) Draw() {
	q.vao.Draw()
}
