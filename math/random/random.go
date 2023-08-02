package random

import (
	"math/rand"
	"time"
)

func init() {
	seed := time.Now().Nanosecond()
	Seed(seed)
}

func Seed(seed int) {
	rand.Seed(int64(seed))
}

func Range(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func Chance(chance float32) bool {
	return Range(0, 1) <= chance
}

func Choice[T any](slice []T) T {
	idx := rand.Intn(len(slice))
	return slice[idx]
}
