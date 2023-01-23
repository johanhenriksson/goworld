package util

import (
	"log"
	"time"
)

var timers = map[string]time.Time{}

func Timer(name string) {
	timers[name] = time.Now()
}

func Elapsed(name string) float32 {
	if start, exists := timers[name]; exists {
		dt := float32(time.Since(start).Seconds())
		log.Printf("Elapsed %s=%.2fms\n", name, dt*1000)
		return dt
	}
	return 0
}
