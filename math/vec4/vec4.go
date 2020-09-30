package vec4

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

var (
	One   = T{1, 1, 1, 1}
	Zero  = T{0, 0, 0, 0}
	UnitX = T{1, 0, 0, 0}
	UnitY = T{0, 1, 0, 0}
	UnitZ = T{0, 0, 1, 0}
	UnitW = T{0, 0, 0, 1}
)

// T holds a 4-component vector of 32-bit floats
type T struct {
	X, Y, Z, W float32
}

func (v T) Slice() [4]float32 {
	return [4]float32{v.X, v.Y, v.Z, v.W}
}

// Length returns the length of the vector.
// See also LengthSqr and Normalize.
func (v T) Length() float32 {
	return math.Sqrt(v.LengthSqr())
}

// LengthSqr returns the squared length of the vector.
// See also Length and Normalize.
func (v T) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
}

// Abs sets every component of the vector to its absolute value.
func (v T) Abs() T {
	return T{
		math.Abs(v.X),
		math.Abs(v.Y),
		math.Abs(v.Z),
		math.Abs(v.W),
	}
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
	v.Z *= s
	v.W *= s
}

// Normalized returns a unit length normalized copy of the vector.
func (v T) Normalized() T {
	v.Normalize()
	return v
}

// Scaled the vector
func (v T) Scaled(f float32) T {
	return T{v.X * f, v.Y * f, v.Z * f, v.W * f}
}

// Scale the vector by a constant (in-place)
func (v *T) Scale(f float32) {
	v.X *= f
	v.Y *= f
	v.Z *= f
	v.W *= f
}

// Invert the vector components
func (v *T) Invert() {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	v.W = -v.W
}

// Inverted returns an inverted vector
func (v T) Inverted() T {
	v.Invert()
	return v
}

// Add each element of the vector with the corresponding element of another vector
func (v T) Add(v2 T) T {
	return T{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z, v.W + v2.W}
}

// Sub subtracts each element of the vector with the corresponding element of another vector
func (v T) Sub(v2 T) T {
	return T{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z, v.W - v2.W}
}

// Mul multiplies each element of the vector with the corresponding element of another vector
func (v T) Mul(v2 T) T {
	return T{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z, v.W * v2.W}
}

// XY returns a 2-component vector with the X, Y components of this vector
func (v T) XY() vec2.T {
	return vec2.T{X: v.X, Y: v.Y}
}

func (v T) XYZ() vec3.T {
	return vec3.T{X: v.X, Y: v.Y, Z: v.Z}
}

// Div divides each element of the vector with the corresponding element of another vector
func (v T) Div(v2 T) T {
	return T{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z, v.W / v2.W}
}
