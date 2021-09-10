package vec2

import (
	"fmt"

	"github.com/johanhenriksson/goworld/math"
)

var (
	// Zero is the zero vector
	Zero = T{0, 0}

	// One is the one vector
	One = T{1, 1}

	// UnitX is the unit vector in the X direction
	UnitX = T{1, 0}

	// UnitY is the unit vector in the Y direction
	UnitY = T{0, 1}
)

// T holds a 2-component vector of 32-bit floats
type T struct {
	X, Y float32
}

// Slice converts the vector into a 2-element slice of float32
func (v T) Slice() [2]float32 {
	return [2]float32{v.X, v.Y}
}

// Length returns the length of the vector.
// See also LengthSqr and Normalize.
func (v T) Length() float32 {
	return math.Sqrt(v.LengthSqr())
}

// LengthSqr returns the squared length of the vector.
// See also Length and Normalize.
func (v T) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y
}

// Abs sets every component of the vector to its absolute value.
func (v T) Abs() T {
	return T{math.Abs(v.X), math.Abs(v.Y)}
}

// Normalize normalizes the vector to unit length.
func (v *T) Normalize() {
	sl := v.LengthSqr()
	if sl == 0 || sl == 1 {
		return
	}
	s := 1 / math.Sqrt(sl)
	v.X *= s
	v.Y *= s
}

// Normalized returns a unit length normalized copy of the vector.
func (v T) Normalized() T {
	v.Normalize()
	return v
}

// Scaled returns a scaled copy of the vector.
func (v T) Scaled(f float32) T {
	return T{v.X * f, v.Y * f}
}

// Scale the vector by a constant (in-place)
func (v *T) Scale(f float32) {
	v.X *= f
	v.Y *= f
}

// Swap returns a new vector with components swapped.
func (v T) Swap() T {
	return T{v.Y, v.X}
}

// Invert components in place
func (v *T) Invert() {
	v.X = -v.X
	v.Y = -v.Y
}

// Inverted returns a new vector with inverted components
func (v T) Inverted() T {
	return T{-v.X, -v.Y}
}

// Add each element of the vector with the corresponding element of another vector
func (v T) Add(v2 T) T {
	return T{v.X + v2.X, v.Y + v2.Y}
}

// Sub subtracts each element of the vector with the corresponding element of another vector
func (v T) Sub(v2 T) T {
	return T{v.X - v2.X, v.Y - v2.Y}
}

// Mul multiplies each element of the vector with the corresponding element of another vector
func (v T) Mul(v2 T) T {
	return T{v.X * v2.X, v.Y * v2.Y}
}

// Div divides each element of the vector with the corresponding element of another vector
func (v T) Div(v2 T) T {
	return T{v.X / v2.X, v.Y / v2.Y}
}

func (v T) ApproxEqual(v2 T) bool {
	epsilon := float32(0.0001)
	return Distance(v, v2) < epsilon
}

func (v T) String() string {
	return fmt.Sprintf("%.3f,%.3f", v.X, v.Y)
}
