package math

import (
	opensimplex "github.com/ojrac/opensimplex-go"
)

// Noise utility to sample simplex noise
type Noise struct {
	opensimplex.Noise
	Seed int
	Freq float32
}

// NewNoise creates a new noise struct from a seed value and a frequency factor.
func NewNoise(seed int, freq float32) *Noise {
	return &Noise{
		Noise: opensimplex.New(int64(seed)),
		Seed:  seed,
		Freq:  freq,
	}
}

// Sample the noise at a certain position
func (n *Noise) Sample(x, y, z int) float32 {
	// jeez
	fx, fy, fz := float64(float32(x)*n.Freq), float64(float32(y)*n.Freq), float64(float32(z)*n.Freq)
	return float32(n.Eval3(fx, fy, fz))
}
