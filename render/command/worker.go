package command

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type CommandFn func(*Buffer)

// Workers manage a command pool thread
type Worker interface {
	Submit(SubmitInfo)
	Destroy()
	Flush()
	Invoke(func())
}

type Workers []Worker

type worker struct {
	device *device.Device
	queue  device.Queue
	name   string
	pool   *Pool
	work   *Channel

	fences *sync.FencePool
}

func NewWorker(device *device.Device, name string, queue device.Queue) Worker {
	pool := NewPool(device, core1_0.CommandPoolCreateTransient, queue.FamilyIndex())

	name = fmt.Sprintf("%s:%d:%x", name, queue.FamilyIndex(), queue.Ptr().Handle())
	device.SetDebugObjectName(driver.VulkanHandle(queue.Ptr().Handle()), core1_0.ObjectTypeQueue, name)

	w := &worker{
		device: device,
		name:   name,
		queue:  queue,
		pool:   pool,
		work:   NewChannel(100),
		fences: sync.NewFencePool(device, name),
	}
	go w.workloop()
	return w
}

func (w *worker) workloop() {
	// lock the worker to its current thread
	runtime.LockOSThread()

	// work loop
	checkFences := time.NewTicker(time.Millisecond)
	for {
		select {
		case work, more := <-w.work.Recv():
			if !more {
				return
			}
			work()

		case <-checkFences.C:
			w.fences.Poll()
		}
	}
}

// Invoke schedules a callback to be called from the worker thread
func (w *worker) Invoke(callback func()) {
	w.work.Invoke(callback)
}

type SubmitInfo struct {
	Commands Recorder
	Marker   string
	Wait     []Wait
	Signal   []sync.Semaphore
	Callback func()
}

type Wait struct {
	Semaphore sync.Semaphore
	Mask      core1_0.PipelineStageFlags
}

// Submit the current batch of command buffers
// Blocks until the queue submission is confirmed
func (w *worker) Submit(submit SubmitInfo) {
	w.work.Invoke(func() {
		w.submit(submit)
	})
}

var noBuffers = []core1_0.CommandBuffer{}

func (w *worker) submit(submit SubmitInfo) {
	if submit.Commands == nil {
		panic("no commands submit. marker: " + submit.Marker)
	}

	debug.SetPanicOnFault(true)

	buffers := noBuffers
	if submit.Commands != Empty {
		// allocate & record command buffers
		buffer := w.pool.Allocate(core1_0.CommandBufferLevelPrimary)
		buffer.Begin()
		submit.Commands.Apply(buffer)
		buffer.End()

		// set debug name
		w.device.SetDebugObjectName(driver.VulkanHandle(buffer.Ptr().Handle()), core1_0.ObjectTypeCommandBuffer, submit.Marker)

		buffers = []core1_0.CommandBuffer{buffer.Ptr()}
	}

	// submit buffers to the given queue
	fence := w.fences.Next()
	w.queue.Ptr().Submit(fence.Ptr(), []core1_0.SubmitInfo{
		{
			CommandBuffers:   buffers,
			SignalSemaphores: util.Map(submit.Signal, func(sem sync.Semaphore) core1_0.Semaphore { return sem.Ptr() }),
			WaitSemaphores:   util.Map(submit.Wait, func(w Wait) core1_0.Semaphore { return w.Semaphore.Ptr() }),
			WaitDstStageMask: util.Map(submit.Wait, func(w Wait) core1_0.PipelineStageFlags { return w.Mask }),
		},
	})

	// clean up command buffers and run the callback once the fence is signaled
	w.fences.Watch(fence, func() {
		// free buffers
		if len(buffers) > 0 {
			w.device.Ptr().FreeCommandBuffers(buffers)
		}

		// run callback (on the worker thead)
		if submit.Callback != nil {
			submit.Callback()
		}
	})
}

func (w *worker) Destroy() {
	w.work.Stop()
	w.fences.Destroy()
	w.pool.Destroy()
}

func (w *worker) Flush() {
	w.work.Flush()
}
