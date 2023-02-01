package command

import (
	"fmt"
	"runtime"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type CommandFn func(Buffer)

// Workers manage a command pool thread
type Worker interface {
	Queue(CommandFn)
	Submit(SubmitInfo)
	Destroy()
	Flush()
	Present(swap swapchain.T, ctx *swapchain.Context)
}

type Workers []Worker

type worker struct {
	device    device.T
	queue     core1_0.Queue
	pool      Pool
	work      chan func()
	stop      chan bool
	batch     []Buffer
	callbacks map[sync.Fence]func()
}

func NewWorker(device device.T, queueFlags core1_0.QueueFlags, queueIndex int) Worker {
	pool := NewPool(device, core1_0.CommandPoolCreateTransient, queueIndex)
	queue := device.GetQueue(queueIndex, queueFlags)
	device.SetDebugObjectName(driver.VulkanHandle(queue.Handle()), core1_0.ObjectTypeQueue, fmt.Sprintf("Worker:%d", queueIndex))

	w := &worker{
		device:    device,
		queue:     queue,
		pool:      pool,
		work:      make(chan func(), 100),
		stop:      make(chan bool),
		batch:     make([]Buffer, 0, 10),
		callbacks: map[sync.Fence]func(){},
	}

	go w.run()

	return w
}

func (w *worker) run() {
	// claim the current thread
	// all command pool operations and buffer recording will execute on this thread
	runtime.LockOSThread()

	// work loop
	running := true
	for running {
		for fence, callback := range w.callbacks {
			if fence.Done() {
				callback()
				delete(w.callbacks, fence)
			}
		}
		select {
		case work := <-w.work:
			work()
		case <-w.stop:
			running = false
		default:
		}
	}

	// dealloc
	w.pool.Destroy()

	// close command channels
	close(w.stop)
	close(w.work)
	w.stop = nil
	w.work = nil

	// return the thread
	runtime.UnlockOSThread()
}

func (w *worker) Queue(batch CommandFn) {
	w.work <- func() {
		w.enqueue(batch)
	}
}

func (w *worker) enqueue(batch CommandFn) {
	// allocate a new buffer
	buf := w.pool.Allocate(core1_0.CommandBufferLevelPrimary)

	// record commands
	buf.Begin()
	defer buf.End()
	batch(buf)

	// append to the next batch
	w.batch = append(w.batch, buf)
}

type SubmitInfo struct {
	Marker string
	Wait   []Wait
	Signal []sync.Semaphore
	Then   func()
}

type Wait struct {
	Semaphore sync.Semaphore
	Mask      core1_0.PipelineStageFlags
}

// Submit the current batch of command buffers
// Blocks until the queue submission is confirmed
func (w *worker) Submit(submit SubmitInfo) {
	w.work <- func() {
		w.submit(submit)
	}
}

func (w *worker) submit(submit SubmitInfo) {
	buffers := util.Map(w.batch, func(buf Buffer) core1_0.CommandBuffer { return buf.Ptr() })

	// create a cleanup callback
	// todo: reuse fences
	fence := sync.NewFence(w.device, "WorkSubmit", false)

	w.callbacks[fence] = func() {
		// free buffers
		if len(buffers) > 0 {
			w.device.Ptr().FreeCommandBuffers(buffers)
		}

		// free fence
		fence.Destroy()

		// run callback if provided
		if submit.Then != nil {
			submit.Then()
		}
	}

	// submit buffers to the given queue
	w.queue.Submit(fence.Ptr(), []core1_0.SubmitInfo{
		{
			CommandBuffers:   buffers,
			SignalSemaphores: util.Map(submit.Signal, func(sem sync.Semaphore) core1_0.Semaphore { return sem.Ptr() }),
			WaitSemaphores:   util.Map(submit.Wait, func(w Wait) core1_0.Semaphore { return w.Semaphore.Ptr() }),
			WaitDstStageMask: util.Map(submit.Wait, func(w Wait) core1_0.PipelineStageFlags { return w.Mask }),
		},
	})

	// clear batch slice but keep memory
	w.batch = w.batch[:0]
}

func (w *worker) Present(swap swapchain.T, ctx *swapchain.Context) {
	w.work <- func() {
		swap.Present(w.queue, ctx)
	}
}

func (w *worker) Destroy() {
	// run all pending cleanups
	for _, callback := range w.callbacks {
		callback()
	}
	w.callbacks = nil

	if w.stop != nil {
		w.stop <- true
		<-w.stop
	}
}

func (w *worker) Flush() {
	done := make(chan struct{})
	w.work <- func() {
		done <- struct{}{}
	}
	<-done
	close(done)
}
