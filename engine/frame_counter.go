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
	start   time.Time
	now     time.Time
	elapsed float32
	delta   float32
}

func NewFrameCounter(samples int) *framecounter {
	return &framecounter{
		samples: samples,
		last:    time.Now().UnixNano(),
		frames:  make([]int64, 0, samples),
		start:   time.Now(),
		now:     time.Now(),
	}
}

type Timing struct {
	Current float32
	Average float32
	Max     float32
}

func (fc *framecounter) Elapsed() float32 {
	return fc.elapsed
}

func (fc *framecounter) Delta() float32 {
	return fc.delta
}

func (fc *framecounter) Update() {
	// clock
	now := time.Now()
	fc.delta = float32(now.Sub(fc.now).Seconds())
	fc.now = now
	fc.elapsed = float32(fc.now.Sub(fc.start).Seconds())

	// fps
	ft := fc.now.UnixNano()
	ns := ft - fc.last
	fc.last = ft
	if len(fc.frames) < fc.samples {
		fc.frames = append(fc.frames, ns)
	} else {
		fc.frames[fc.next%fc.samples] = ns
	}
	fc.next++
}

func (fc *framecounter) Sample() Timing {
	tot := int64(0)
	max := int64(0)
	for _, f := range fc.frames {
		tot += f
		max = math.Max(max, f)
	}

	current := fc.frames[(fc.next-1)%fc.samples]
	return Timing{
		Average: float32(tot) / float32(len(fc.frames)) / 1e9,
		Max:     float32(max) / 1e9,
		Current: float32(current) / 1e9,
	}
}
