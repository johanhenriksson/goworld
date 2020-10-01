package vec3

import (
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec2"
)

func New(x, y, z float32) T {
	return T{x, y, z}
}

func NewI(x, y, z int) T {
	return T{float32(x), float32(y), float32(z)}
}

func Extend(v vec2.T, z float32) T {
	return T{v.X, v.Y, z}
}

// Dot returns the dot product of two vectors.
func Dot(a, b T) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross returns the cross product of two vectors.
func Cross(a, b *T) T {
	return T{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func Random(min, max T) T {
	return T{
		random.Range(min.X, max.X),
		random.Range(min.Y, max.Y),
		random.Range(min.Z, max.Z),
	}
}