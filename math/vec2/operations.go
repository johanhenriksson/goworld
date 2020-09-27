package vec2

func New(x, y float32) T {
	return T{X: x, Y: y}
}

func NewI(x, y int) T {
	return T{X: float32(x), Y: float32(y)}
}

// Dot returns the dot product of two vectors.
func Dot(a, b T) float32 {
	return a.X*b.X + a.Y*b.Y
}
