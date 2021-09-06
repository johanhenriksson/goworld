package vec3

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
)

var (
	// Zero is the zero vector
	Zero = T{0, 0, 0}

	// One is the unit vector
	One = T{1, 1, 1}

	// UnitX is the unit vector in the X direction
	UnitX = T{1, 0, 0}

	// UnitXN is the unit vector in the negative X direction
	UnitXN = T{-1, 0, 0}

	// UnitY is the unit vector in the Y direction
	UnitY = T{0, 1, 0}

	// UnitYN is the unit vector in the negative Y direction
	UnitYN = T{0, -1, 0}

	// UnitZ is the unit vector in the Z direction
	UnitZ = T{0, 0, 1}

	// UnitZN is the unit vector in the negative Z direction
	UnitZN = T{0, 0, -1}
)

// T holds a 3-component vector of 32-bit floats
type T struct {
	X, Y, Z float32
}

// Slice converts the vector into a 3-element slice of float32
func (v T) Slice() [3]float32 {
	return [3]float32{v.X, v.Y, v.Z}
}

// Length returns the length of the vector.
// See also LengthSqr and Normalize.
func (v T) Length() float32 {
	return math.Sqrt(v.LengthSqr())
}

// LengthSqr returns the squared length of the vector.
// See also Length and Normalize.
func (v T) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Abs returns a copy containing the absolute values of the vector components.
func (v T) Abs() T {
	return T{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)}
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
}

// Normalized returns a unit length normalized copy of the vector.
func (v T) Normalized() T {
	v.Normalize()
	return v
}

// Scale the vector by a constant (in-place)
func (v *T) Scale(f float32) {
	v.X *= f
	v.Y *= f
	v.Z *= f
}

// Scaled returns a scaled vector
func (v T) Scaled(f float32) T {
	return T{v.X * f, v.Y * f, v.Z * f}
}

// ScaleI returns a vector scaled by an integer factor
func (v T) ScaleI(i int) T {
	return v.Scaled(float32(i))
}

// Invert the vector components
func (v *T) Invert() {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
}

// Inverted returns an inverted vector
func (v *T) Inverted() T {
	i := *v
	i.Invert()
	return i
}

// Floor each components of the vector
func (v T) Floor() T {
	return T{math.Floor(v.X), math.Floor(v.Y), math.Floor(v.Z)}
}

// Ceil each component of the vector
func (v T) Ceil() T {
	return T{math.Ceil(v.X), math.Ceil(v.Y), math.Ceil(v.Z)}
}

// Add each element of the vector with the corresponding element of another vector
func (v T) Add(v2 T) T {
	return T{
		v.X + v2.X,
		v.Y + v2.Y,
		v.Z + v2.Z,
	}
}

// Sub subtracts each element of the vector with the corresponding element of another vector
func (v T) Sub(v2 T) T {
	return T{
		v.X - v2.X,
		v.Y - v2.Y,
		v.Z - v2.Z,
	}
}

// Mul multiplies each element of the vector with the corresponding element of another vector
func (v T) Mul(v2 T) T {
	return T{
		v.X * v2.X,
		v.Y * v2.Y,
		v.Z * v2.Z,
	}
}

// XY returns a 2-component vector with the X, Y components of this vector
func (v T) XY() vec2.T {
	return vec2.T{X: v.X, Y: v.Y}
}

// XZ returns a 2-component vector with the X, Z components of this vector
func (v T) XZ() vec2.T {
	return vec2.T{X: v.X, Y: v.Z}
}

// YZ returns a 2-component vector with the Y, Z components of this vector
func (v T) YZ() vec2.T {
	return vec2.T{X: v.Y, Y: v.Z}
}

// Div divides each element of the vector with the corresponding element of another vector
func (v T) Div(v2 T) T {
	return T{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z}
}

// WithX returns a new vector with the X component set to a given value
func (v T) WithX(x float32) T {
	return T{x, v.Y, v.Z}
}

// WithY returns a new vector with the Y component set to a given value
func (v T) WithY(y float32) T {
	return T{v.X, y, v.Z}
}

// WithZ returns a new vector with the Z component set to a given value
func (v T) WithZ(z float32) T {
	return T{v.X, v.Y, z}
}

func (v T) ApproxEqual(v2 T) bool {
	epsilon := float32(0.0001)
	return Distance(v, v2) < epsilon
}
