package math

import (
	opensimplex "github.com/ojrac/opensimplex-go"
)

// Noise utility to sample simplex noise
type Noise struct {
	Seed  int
	Scale float32
	Freq  float32

	simplex opensimplex.Noise
}

// NewNoise creates a new noise struct from a seed value and a frequency factor.
func NewNoise(seed int, scale, freq float32) *Noise {
	return &Noise{
		Seed:  seed,
		Scale: scale,
		Freq:  freq,

		simplex: opensimplex.New(int64(seed)),
	}
}

// Sample the 2D noise at a certain position
func (n *Noise) Sample2(x, z int) float32 {
	fx, fz := float64(float32(x)*n.Freq), float64(float32(z)*n.Freq)
	return n.Scale * float32(n.simplex.Eval2(fx, fz))
}

// Sample the 3D noise at a certain position
func (n *Noise) Sample3(x, y, z int) float32 {
	// jeez
	fx, fy, fz := float64(float32(x)*n.Freq), float64(float32(y)*n.Freq), float64(float32(z)*n.Freq)
	return n.Scale * float32(n.simplex.Eval3(fx, fy, fz))
}
