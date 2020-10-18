package vec4

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// New returns a new vec4 from its components
func New(x, y, z, w float32) T {
	return T{x, y, z, w}
}

// Extend a vec3 to a vec4 by adding a W component
func Extend(v vec3.T, w float32) T {
	return T{v.X, v.Y, v.Z, w}
}

// Extend2 a vec2 to a vec4 by adding the Z and W components
func Extend2(v vec2.T, z, w float32) T {
	return T{v.X, v.Y, z, w}
}

// Dot returns the dot product of two vectors.
func Dot(a, b T) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}
