package math

import (
	opensimplex "github.com/ojrac/opensimplex-go"
)

type Noise struct {
	opensimplex.Noise
	Seed int
	Freq float32
}

func NewNoise(seed int, freq float32) *Noise {
	return &Noise{
		Noise: opensimplex.New(int64(seed)),
		Seed:  seed,
		Freq:  freq,
	}
}

func (n *Noise) Sample(x, y, z int) float32 {
	// jeez
	fx, fy, fz := float64(float32(x)*n.Freq), float64(float32(y)*n.Freq), float64(float32(z)*n.Freq)
	return float32(n.Eval3(fx, fy, fz))
}
