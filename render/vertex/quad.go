package vertex

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Full-screen quad helper
func ScreenQuad(key string) Mesh {
	return NewTriangles(key, []Vertex{
		T(vec3.New(-1, -1, 0), vec3.Zero, vec2.New(0, 0)),
		T(vec3.New(1, 1, 0), vec3.Zero, vec2.New(1, 1)),
		T(vec3.New(-1, 1, 0), vec3.Zero, vec2.New(0, 1)),
		T(vec3.New(1, -1, 0), vec3.Zero, vec2.New(1, 0)),
	}, []uint16{
		0, 1, 2,
		0, 3, 1,
	})
}
