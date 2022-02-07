package vec2

import "github.com/johanhenriksson/goworld/math"

// New returns a vec2 from its components
func New(x, y float32) T {
	return T{X: x, Y: y}
}

// NewI returns a vec2 from integer components
func NewI(x, y int) T {
	return T{X: float32(x), Y: float32(y)}
}

// Dot returns the dot product of two vectors.
func Dot(a, b T) float32 {
	return a.X*b.X + a.Y*b.Y
}

// Distance returns the euclidian distance between two points.
func Distance(a, b T) float32 {
	return a.Sub(b).Length()
}

func Min(a, b T) T {
	return T{
		X: math.Min(a.X, b.X),
		Y: math.Min(a.Y, b.Y),
	}
}

func Max(a, b T) T {
	return T{
		X: math.Max(a.X, b.X),
		Y: math.Max(a.Y, b.Y),
	}
}
