package vertex

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
)

// C - Colored Vertex
type C struct {
	P vec3.T `vtx:"position,float,3"`
	N vec3.T `vtx:"normal,float,3"`
	C vec4.T `vtx:"color,float,4"`
}

// T - Textured Vertex
type T struct {
	P vec3.T `vtx:"position,float,3"`
	N vec3.T `vtx:"normal,float,3"`
	T vec2.T `vtx:"texcoord,float,2"`
}

type UI struct {
	P vec3.T  `vtx:"position,float,3"`
	C color.T `vtx:"color,float,4"`
	T vec2.T  `vtx:"texcoord,float,2"`
}
