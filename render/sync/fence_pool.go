package sync

import (
	"slices"
	"time"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/samber/lo"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type FencePool struct {
	name      string
	device    *device.Device
	available []Fence
	waiting   []Fence
	callbacks []func()
	mutex     Mutex
}

func NewFencePool(device *device.Device, name string) *FencePool {
	return &FencePool{
		name:   name,
		device: device,
	}
}

// Poll checks all fences in the pool and runs their callbacks if they are done.
// If a fence is done, it is returned to the pool.
func (w *FencePool) Poll() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for i := 0; i < len(w.waiting); i++ {
		fence := w.waiting[i]
		if !fence.Done() {
			continue
		}

		w.callbacks[i]()
		w.waiting = slices.Delete(w.waiting, i, i+1)
		w.callbacks = slices.Delete(w.callbacks, i, i+1)

		// return the fence to the pool
		fence.Reset()
		w.available = append(w.available, fence)

		i--
	}
}

// Next returns a fence from the pool.
// Until the fence is added back to the pool, the caller has ownership of it.
// If there are no fences available, a new one is allocated.
func (w *FencePool) Next() Fence {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(w.available) == 0 {
		return NewFence(w.device, w.name, false)
	}
	fence := w.available[len(w.available)-1]
	w.available = w.available[:len(w.available)-1]
	return fence
}

// Watch adds a fence to the pool and sets a callback to run when the fence is done.
// The callback is run in the same goroutine as the Poll call.
func (w *FencePool) Watch(fence Fence, callback func()) {
	if callback == nil {
		panic("callback cant be nil")
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.waiting = append(w.waiting, fence)
	w.callbacks = append(w.callbacks, callback)
}

// Destroy cleans up all fences in the pool.
// It waits (1 sec) for all fences to be signaled and runs their callbacks.
func (w *FencePool) Destroy() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// wait for any pending fences
	if len(w.waiting) > 0 {
		w.device.Ptr().WaitForFences(true, time.Second, lo.Map(w.waiting, func(f Fence, _ int) core1_0.Fence { return f.Ptr() }))
	}

	// run all pending callbacks
	// this may not be a great idea
	for _, callback := range w.callbacks {
		callback()
	}

	// clean up
	for _, fence := range w.available {
		fence.Destroy()
	}
	for _, fence := range w.waiting {
		fence.Destroy()
	}

	w.available = nil
	w.waiting = nil
	w.callbacks = nil
}
