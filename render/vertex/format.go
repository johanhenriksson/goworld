package vertex

import (
	"github.com/johanhenriksson/goworld/math/byte4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// C - Colored Vertex
type C struct {
	P vec3.T  `vtx:"position,float,3"`
	N vec3.T  `vtx:"normal,float,3"`
	C byte4.T `vtx:"color,byte,4,normalize"`
}

// T - Textured Vertex
type T struct {
	P vec3.T `vtx:"position,float,3"`
	N vec3.T `vtx:"normal,float,3"`
	T vec2.T `vtx:"texcoord,float,2"`
}
