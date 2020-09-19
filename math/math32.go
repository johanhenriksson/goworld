package math

import (
	"math"
)

const E = float32(math.E)
const Pi = float32(math.Pi)
const Sqrt2 = float32(math.Sqrt2)

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func Clamp(v, min, max float32) float32 {
	if v > max {
		return max
	}
	if v < min {
		return min
	}
	return v
}

func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Tan(x float32) float32 {
	return float32(math.Tan(float64(x)))
}

func Sign(x float32) float32 {
	if x > 0 {
		return 1
	}
	return -1
}
