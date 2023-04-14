package math

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Various useful constants.
var (
	MinNormal = float32(1.1754943508222875e-38) // 1 / 2**(127 - 1)
	MinValue  = float32(math.SmallestNonzeroFloat32)
	MaxValue  = float32(math.MaxFloat32)

	InfPos = float32(math.Inf(1))
	InfNeg = float32(math.Inf(-1))
	NaN    = float32(math.NaN())

	E       = float32(math.E)
	Pi      = float32(math.Pi)
	PiOver2 = Pi / 2
	PiOver4 = Pi / 4
	Sqrt2   = float32(math.Sqrt2)

	Epsilon = float32(1e-10)
)

// Abs returns the absolute value of a number
func Abs[T constraints.Float | constraints.Integer](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// Min returns the smaller of two numbers
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the greater of two numbers
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Clamp a value between a minimum and a maximum value
func Clamp[T constraints.Ordered](v, min, max T) T {
	if v > max {
		return max
	}
	if v < min {
		return min
	}
	return v
}

// Ceil a number to the closest integer
func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

// Floor a number to the closest integer
func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

// Mod returns the remainder of a floating point division
func Mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}

// Sqrt returns the square root of a number
func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

// Sin computes the sine of x
func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

// Cos computes the cosine of x
func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

// Tan computes the tangent of x
func Tan(x float32) float32 {
	return float32(math.Tan(float64(x)))
}

// Sign returns the sign of x (-1 or 1)
func Sign(x float32) float32 {
	if x > 0 {
		return 1
	}
	return -1
}

// DegToRad converts degrees to radians
func DegToRad(deg float32) float32 {
	return Pi * deg / 180.0
}

// RadToDeg converts radians to degrees
func RadToDeg(rad float32) float32 {
	return 180.0 * rad / Pi
}

// Equal checks two floats for (approximate) equality
func Equal(a, b float32) bool {
	return EqualThreshold(a, b, Epsilon)
}

// EqualThreshold is a utility function to compare floats.
// It's Taken from http://floating-point-gui.de/errors/comparison/
//
// It is slightly altered to not call Abs when not needed.
//
// This differs from FloatEqual in that it lets you pass in your comparison threshold, so that you can adjust the comparison value to your specific needs
func EqualThreshold(a, b, epsilon float32) bool {
	if a == b { // Handles the case of inf or shortcuts the loop when no significant error has accumulated
		return true
	}

	diff := Abs(a - b)
	if a*b == 0 || diff < MinNormal { // If a or b are 0 or both are extremely close to it
		return diff < epsilon*epsilon
	}

	// Else compare difference
	return diff/(Abs(a)+Abs(b)) < epsilon
}

// Lerp performs linear interpolation between a and b
func Lerp(a, b, f float32) float32 {
	return a + f*(b-a)
}

func Round(f float32) float32 {
	return float32(math.Round(float64(f)))
}

func Snap(f, multiple float32) float32 {
	return Round(f/multiple) * multiple
}

func Pow(f, x float32) float32 {
	return float32(math.Pow(float64(f), float64(x)))
}
