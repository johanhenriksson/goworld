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
	C vec4.T `vtx:"color_0,float,4"`
}

func (v C) Position() vec3.T { return v.P }

// T - Textured Vertex
type T struct {
	P vec3.T `vtx:"position,float,3"`
	N vec3.T `vtx:"normal,float,3"`
	T vec2.T `vtx:"texcoord_0,float,2"`
}

func (v T) Position() vec3.T { return v.P }

type UI struct {
	P vec3.T  `vtx:"position,float,3"`
	C color.T `vtx:"color_0,float,4"`
	T vec2.T  `vtx:"texcoord_0,float,2"`
}

func (v UI) Position() vec3.T { return v.P }

func Min[V Vertex](vertices []V) vec3.T {
	if len(vertices) == 0 {
		return vec3.Zero
	}
	min := vec3.InfPos
	for _, v := range vertices {
		min = vec3.Min(min, v.Position())
	}
	return min
}

func Max[V Vertex](vertices []V) vec3.T {
	if len(vertices) == 0 {
		return vec3.Zero
	}
	max := vec3.InfNeg
	for _, v := range vertices {
		max = vec3.Max(max, v.Position())
	}
	return max
}
