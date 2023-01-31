package command

import (
	"runtime"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type CommandFn func(Buffer)

// Workers manage a command pool thread
type Worker interface {
	Queue(CommandFn)
	Submit(SubmitInfo)
	Destroy()
	Wait()
	Present(swap swapchain.T, ctx swapchain.Context)
}

type Workers []Worker

type worker struct {
	device    device.T
	queue     vk.Queue
	pool      Pool
	work      chan func()
	stop      chan bool
	batch     []Buffer
	callbacks map[sync.Fence]func()
}

func NewWorker(device device.T, queueFlags vk.QueueFlags, queueIndex int) Worker {
	pool := NewPool(device, vk.CommandPoolCreateFlags(vk.CommandPoolCreateTransientBit), queueIndex)
	queue := device.GetQueue(queueIndex, queueFlags)

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
	w.device.WaitIdle()
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
	buf := w.pool.Allocate(vk.CommandBufferLevelPrimary)

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
	Mask      vk.PipelineStageFlagBits
}

// Submit the current batch of command buffers
// Blocks until the queue submission is confirmed
func (w *worker) Submit(submit SubmitInfo) {
	w.work <- func() {
		w.submit(submit)
	}
}

func (w *worker) submit(submit SubmitInfo) {
	buffers := util.Map(w.batch, func(buf Buffer) vk.CommandBuffer { return buf.Ptr() })

	// create a cleanup callback
	// todo: reuse fences
	fence := sync.NewFence(w.device, false)

	w.callbacks[fence] = func() {
		// free buffers
		if len(buffers) > 0 {
			vk.FreeCommandBuffers(w.device.Ptr(), w.pool.Ptr(), uint32(len(buffers)), buffers)
		}

		// free fence
		fence.Destroy()

		// run callback if provided
		if submit.Then != nil {
			submit.Then()
		}
	}

	// submit buffers to the given queue
	info := []vk.SubmitInfo{
		{
			SType:                vk.StructureTypeSubmitInfo,
			CommandBufferCount:   uint32(len(buffers)),
			WaitSemaphoreCount:   uint32(len(submit.Wait)),
			SignalSemaphoreCount: uint32(len(submit.Signal)),
			PCommandBuffers:      buffers,
			PSignalSemaphores:    util.Map(submit.Signal, func(sem sync.Semaphore) vk.Semaphore { return sem.Ptr() }),
			PWaitSemaphores:      util.Map(submit.Wait, func(w Wait) vk.Semaphore { return w.Semaphore.Ptr() }),
			PWaitDstStageMask:    util.Map(submit.Wait, func(w Wait) vk.PipelineStageFlags { return vk.PipelineStageFlags(w.Mask) }),
		},
	}
	vk.QueueSubmit(w.queue, 1, info, fence.Ptr())

	// clear batch slice but keep memory
	w.batch = w.batch[:0]
}

func (w *worker) Present(swap swapchain.T, ctx swapchain.Context) {
	var waits []vk.Semaphore
	if ctx.RenderComplete != nil {
		waits = []vk.Semaphore{ctx.RenderComplete.Ptr()}
	}
	w.work <- func() {
		presentInfo := vk.PresentInfo{
			SType:              vk.StructureTypePresentInfo,
			WaitSemaphoreCount: uint32(len(waits)),
			PWaitSemaphores:    waits,
			SwapchainCount:     1,
			PSwapchains:        []vk.Swapchain{swap.Ptr()},
			PImageIndices:      []uint32{uint32(ctx.Index)},
		}
		vk.QueuePresent(w.queue, &presentInfo)
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

func (w *worker) Wait() {
	done := make(chan struct{})
	w.work <- func() {
		done <- struct{}{}
	}
	<-done
	close(done)
}
