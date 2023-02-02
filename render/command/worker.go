package command

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type CommandFn func(Buffer)

// Workers manage a command pool thread
type Worker interface {
	Ptr() core1_0.Queue
	Queue(CommandFn)
	Submit(SubmitInfo)
	Destroy()
	Flush()
	Invoke(func())
}

type Workers []Worker

type worker struct {
	device    device.T
	name      string
	queue     core1_0.Queue
	pool      Pool
	work      chan func()
	batch     []Buffer
	callbacks map[sync.Fence]func()
}

func NewWorker(device device.T, queueFlags core1_0.QueueFlags, queueIndex int) Worker {
	pool := NewPool(device, core1_0.CommandPoolCreateTransient, queueIndex)
	queue := device.GetQueue(queueIndex, queueFlags)

	name := fmt.Sprintf("Worker:%d", queueIndex)
	device.SetDebugObjectName(driver.VulkanHandle(queue.Handle()), core1_0.ObjectTypeQueue, name)

	w := &worker{
		device:    device,
		name:      name,
		queue:     queue,
		pool:      pool,
		work:      make(chan func(), 100),
		batch:     make([]Buffer, 0, 10),
		callbacks: map[sync.Fence]func(){},
	}

	go w.run()

	return w
}

func (w *worker) Ptr() core1_0.Queue {
	return w.queue
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
		case work, more := <-w.work:
			if !more {
				running = false
				break
			}
			work()
		default:
			time.Sleep(100 * time.Microsecond)
		}
	}

	// dealloc
	w.pool.Destroy()

	// close command channels
	w.work = nil

	// return the thread
	runtime.UnlockOSThread()
	log.Println(w.name, "exited")
}

func (w *worker) Invoke(callback func()) {
	w.work <- callback
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

func (w *worker) Destroy() {
	// run all pending cleanups
	for _, callback := range w.callbacks {
		callback()
	}
	w.callbacks = nil

	w.work <- func() {
		close(w.work)
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
