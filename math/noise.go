package math

import (
	opensimplex "github.com/ojrac/opensimplex-go"
)

type Noise struct {
	opensimplex.Noise
	Seed int
	Freq float64
}

func NewNoise(seed int, freq float64) *Noise {
	return &Noise{
		Noise: opensimplex.New(int64(seed)),
		Seed:  seed,
		Freq:  freq,
	}
}

func (n *Noise) Sample(x, y, z int) float64 {
	fx, fy, fz := float64(x)*n.Freq, float64(y)*n.Freq, float64(z)*n.Freq
	return n.Eval3(fx, fy, fz)
}