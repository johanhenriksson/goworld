package math

import (
	"math"
)

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
