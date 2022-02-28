package screen_quad

import (
	"log"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vertex_array"
)

type T interface {
	Draw()
}

// quad is a fullscreen quad used for render passes
type quad struct {
	vao vertex_array.T
}

// NewQuad creates a new quad with a given material
func New(shader shader.T) T {
	q := &quad{
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

	ptrs := vertex.ParsePointers(vtx)
	ptrs.Bind(shader)
	log.Println(ptrs)
	q.vao.BufferTo(ptrs, vtx)

	return q
}

func (q *quad) Draw() {
	q.vao.Draw()
}
