package command

import (
	"runtime"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type CommandFn func(Buffer)

// Workers manage a command pool thread
type Worker interface {
	Queue(CommandFn)
	Submit(SubmitInfo)
	Wait()
	Destroy()
}

type Workers []Worker

type worker struct {
	device   device.T
	queue    vk.Queue
	pool     Pool
	input    chan CommandFn
	signal   chan SubmitInfo
	complete chan bool
	stop     chan bool
	batch    []Buffer
	destroy  []vk.CommandBuffer
	fence    sync.Fence
}

func NewWorker(device device.T, queueFlags vk.QueueFlags) Worker {
	pool := NewPool(device, vk.CommandPoolCreateFlags(vk.CommandPoolCreateTransientBit), queueFlags)
	queue := device.GetQueue(0, queueFlags)

	w := &worker{
		device:   device,
		queue:    queue,
		pool:     pool,
		input:    make(chan CommandFn),
		signal:   make(chan SubmitInfo),
		complete: make(chan bool),
		stop:     make(chan bool),
		batch:    make([]Buffer, 0, 10),
		destroy:  make([]vk.CommandBuffer, 0, 10),
		fence:    sync.NewFence(device, true),
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
		select {
		case batch := <-w.input:
			w.enqueue(batch)
		case info := <-w.signal:
			w.submit(info)
			w.complete <- true
		case <-w.stop:
			running = false
		}
	}

	// dealloc
	w.device.WaitIdle()
	w.fence.Destroy()
	w.pool.Destroy()

	// close command channels
	w.complete <- true
	close(w.input)
	close(w.signal)
	close(w.stop)
	close(w.complete)

	// return the thread
	runtime.UnlockOSThread()
}

func (w *worker) Queue(batch CommandFn) { w.input <- batch }

func (w *worker) enqueue(batch CommandFn) {
	// allocate a new buffer
	buf := w.pool.Allocate(vk.CommandBufferLevelPrimary)

	// record commands
	buf.Begin()
	batch(buf)
	buf.End()

	// append to the next batch
	w.batch = append(w.batch, buf)
}

type SubmitInfo struct {
	Wait     []sync.Semaphore
	Signal   []sync.Semaphore
	WaitMask []vk.PipelineStageFlags
}

// Submit the current batch of command buffers
// Blocks until the queue submission is confirmed
func (w *worker) Submit(submit SubmitInfo) {
	w.signal <- submit
	<-w.complete
}

func (w *worker) submit(submit SubmitInfo) {
	// make sure we have something to submit
	if len(w.batch) == 0 {
		return
	}

	// wait for any ongoing gpu execution to complete
	w.fence.Wait()
	w.fence.Reset()

	// delete the previous batch buffers
	// we can free all of them with a single call instead of Destroy()
	if len(w.destroy) > 0 {
		vk.FreeCommandBuffers(w.device.Ptr(), w.pool.Ptr(), uint32(len(w.destroy)), w.destroy)
		w.destroy = w.destroy[:0]
	}

	// submit buffers to the given queue
	buffers := util.Map(w.batch, func(i int, buf Buffer) vk.CommandBuffer { return buf.Ptr() })
	info := []vk.SubmitInfo{
		{
			SType:                vk.StructureTypeSubmitInfo,
			CommandBufferCount:   uint32(len(buffers)),
			WaitSemaphoreCount:   uint32(len(submit.Wait)),
			SignalSemaphoreCount: uint32(len(submit.Signal)),
			PCommandBuffers:      buffers,
			PWaitSemaphores:      util.Map(submit.Wait, func(i int, sem sync.Semaphore) vk.Semaphore { return sem.Ptr() }),
			PSignalSemaphores:    util.Map(submit.Signal, func(i int, sem sync.Semaphore) vk.Semaphore { return sem.Ptr() }),
			PWaitDstStageMask:    submit.WaitMask,
		},
	}
	vk.QueueSubmit(w.queue, uint32(len(w.batch)), info, w.fence.Ptr())

	// add submitted buffers to destroy list
	w.destroy = append(w.destroy, buffers...)

	// clear batch slice but keep memory
	w.batch = w.batch[:0]
}

func (w *worker) Wait() {
	w.fence.Wait()
}

func (w *worker) Destroy() {
	w.stop <- true
	<-w.complete
}
