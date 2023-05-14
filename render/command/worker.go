package command

import (
	"fmt"
	"runtime/debug"

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
	device device.T
	name   string
	queue  core1_0.Queue
	pool   Pool
	batch  []Buffer
	work   *ThreadWorker
}

func NewWorker(device device.T, queueFlags core1_0.QueueFlags, queueIndex int) Worker {
	pool := NewPool(device, core1_0.CommandPoolCreateTransient, 0)
	queue := device.GetQueue(queueIndex, queueFlags)

	name := fmt.Sprintf("Worker:%d", queueIndex)
	device.SetDebugObjectName(driver.VulkanHandle(queue.Handle()), core1_0.ObjectTypeQueue, name)

	return &worker{
		device: device,
		name:   name,
		queue:  queue,
		pool:   pool,
		batch:  make([]Buffer, 0, 128),
		work:   NewThreadWorker(name, 100, true),
	}
}

func (w *worker) Ptr() core1_0.Queue {
	return w.queue
}

// Invoke schedules a callback to be called from the worker thread
func (w *worker) Invoke(callback func()) {
	w.work.Invoke(callback)
}

func (w *worker) Queue(batch CommandFn) {
	w.work.Invoke(func() {
		w.enqueue(batch)
	})
}

func (w *worker) enqueue(batch CommandFn) {
	// todo: dont make a command buffer for each call to Queue() !!
	//       instead, allocate and record everything that we've batched just prior to submission

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

func (w *worker) submit(submit SubmitInfo) {
	debug.SetPanicOnFault(true)
	buffers := util.Map(w.batch, func(buf Buffer) core1_0.CommandBuffer { return buf.Ptr() })

	// create a cleanup callback
	// todo: reuse fences
	fence := sync.NewFence(w.device, submit.Marker, false)

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

	// fire up a cleanup goroutine that will execute when the work fence is signalled
	go func() {
		fence.Wait()
		fence.Destroy()

		w.work.Invoke(func() {
			// free buffers
			if len(buffers) > 0 {
				w.device.Ptr().FreeCommandBuffers(buffers)
			}

			// run callback (on the worker thead)
			if submit.Callback != nil {
				submit.Callback()
			}
		})
	}()
}

func (w *worker) Destroy() {
	w.work.Stop()
	w.pool.Destroy()
}

func (w *worker) Flush() {
	w.work.Flush()
}
