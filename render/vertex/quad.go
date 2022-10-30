package vertex

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Full-screen quad helper
func ScreenQuad() Mesh {
	return NewTriangles("screen_quad", []T{
		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
		{P: vec3.New(-1, 1, 0), T: vec2.New(0, 1)},
		{P: vec3.New(1, -1, 0), T: vec2.New(1, 0)},
	}, []uint16{
		0, 1, 2,
		0, 3, 1,
	})
}
