package screen_quad

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	Draw()
}

// quad is a fullscreen quad used for render passes
type quad struct {
	vao vertex.Array
}

// NewQuad creates a new quad with a given material
func New(shader shader.T) T {
	q := &quad{
		vao: gl_vertex_array.New(vertex.Triangles),
	}

	mesh := vertex.NewTriangles("screen_quad", []vertex.T{
		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
		{P: vec3.New(-1, 1, 0), T: vec2.New(0, 1)},
		{P: vec3.New(1, -1, 0), T: vec2.New(1, 0)},
	}, []uint8{
		0, 1, 2,
		0, 3, 1,
	})

	q.vao.Buffer("vertex", mesh.VertexData())
	q.vao.Buffer("index", mesh.IndexData())

	ptrs := mesh.Pointers()
	ptrs.Bind(shader)

	q.vao.SetPointers(ptrs)
	q.vao.SetIndexSize(1)
	q.vao.SetElements(6)

	return q
}

func (q *quad) Draw() {
	q.vao.Draw()
}
