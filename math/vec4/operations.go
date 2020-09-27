package vec4

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func Extend(v vec3.T, w float32) T {
	return T{v.X, v.Y, v.Z, w}
}

func Extend2(v vec2.T, z, w float32) T {
	return T{v.X, v.Y, z, w}
}

// Dot returns the dot product of two vectors.
func Dot(a, b T) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}
