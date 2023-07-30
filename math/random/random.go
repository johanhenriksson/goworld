package random

import (
	"math/rand"
	"time"
)

func init() {
	s := int64(time.Now().Nanosecond())
	rand.Seed(s)
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
