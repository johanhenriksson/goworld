package random

import (
	"math/rand"
)

func Range(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func Chance(chance float32) bool {
	return Range(0, 1) <= chance
}
