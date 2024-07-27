package terrain

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
)

type Vertex struct {
	P vec3.T   `vtx:"position,float,3"`
	N vec3.T   `vtx:"normal,float,3"`
	T vec2.T   `vtx:"texcoord_0,float,2"`
	W vec4.T   `vtx:"weights,float,4"`
	I [4]uint8 `vtx:"indices,uint8,4"`
}

func (v Vertex) Position() vec3.T { return v.P }
