package vertex

import (
	"encoding/gob"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

func init() {
	gob.Register(Vertex{})
}

// Standard vertex format
type Vertex struct {
	P  vec3.T  `vtx:"position,float,3"`
	Tx float32 `vtx:"tex_x,float,1"`
	N  vec3.T  `vtx:"normal,float,3"`
	Ty float32 `vtx:"tex_y,float,1"`
	C  color.T `vtx:"color,float,4"`
}

func (v Vertex) Position() vec3.T { return v.P }

// New creates a new vertex with position, normal, texture coordinates and color
func New(p vec3.T, n vec3.T, t vec2.T, c color.T) Vertex {
	return Vertex{P: p, Tx: t.X, N: n, Ty: t.Y, C: c}
}

// P defines a vertex with a position
func P(p vec3.T) Vertex {
	return Vertex{P: p, C: color.White}
}

// C defines a vertex with a position and color
func C(p vec3.T, n vec3.T, c color.T) Vertex {
	return Vertex{P: p, N: n, C: c}
}

// T defines a vertex with a position, normal and texture coordinates
func T(p vec3.T, n vec3.T, t vec2.T) Vertex {
	return Vertex{P: p, N: n, Tx: t.X, Ty: t.Y, C: color.White}
}

func Min[V VertexFormat](vertices []V) vec3.T {
	if len(vertices) == 0 {
		return vec3.Zero
	}
	min := vec3.InfPos
	for _, v := range vertices {
		min = vec3.Min(min, v.Position())
	}
	return min
}

func Max[V VertexFormat](vertices []V) vec3.T {
	if len(vertices) == 0 {
		return vec3.Zero
	}
	max := vec3.InfNeg
	for _, v := range vertices {
		max = vec3.Max(max, v.Position())
	}
	return max
}
