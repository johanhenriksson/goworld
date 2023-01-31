package engine

import (
	"time"

	"github.com/johanhenriksson/goworld/math"
)

type framecounter struct {
	next    int
	samples int
	last    int64
	frames  []int64
}

func NewFrameCounter(samples int) *framecounter {
	return &framecounter{
		samples: samples,
		last:    time.Now().UnixNano(),
		frames:  make([]int64, samples),
	}
}

type Timing struct {
	Current float32
	Average float32
	Max     float32
}

func (fc *framecounter) Sample() Timing {
	ft := time.Now().UnixNano()
	ns := ft - fc.last
	fc.last = ft
	fc.frames[fc.next%fc.samples] = ns
	fc.next++
	tot := int64(0)
	max := int64(0)
	for _, f := range fc.frames {
		tot += f
		max = math.Max(max, f)
	}
	return Timing{
		Average: float32(tot) / float32(fc.samples) / 1e9,
		Max:     float32(max) / 1e9,
		Current: float32(ns) / 1e9,
	}
}
