package command

import (
	"runtime"
)

type InvokeFunc func()

type ThreadWorker struct {
	name   string
	buffer int
	work   chan InvokeFunc
}

func NewThreadWorker(name string, buffer int, locked bool) *ThreadWorker {
	w := &ThreadWorker{
		name:   name,
		buffer: buffer,
		work:   make(chan InvokeFunc, buffer),
	}
	go w.workloop(locked)
	return w
}

func (tw *ThreadWorker) workloop(locked bool) {
	// lock the worker to its current thread
	if locked {
		runtime.LockOSThread()
	}

	// work loop
	for {
		work, more := <-tw.work
		if !more {
			break
		}
		work()
	}
}

// Invoke schedules a callback to be called from the worker thread
func (tw *ThreadWorker) Invoke(callback InvokeFunc) {
	// debug.PrintStack()
	tw.work <- callback
}

// InvokeSync schedules a callback to be called on the worker thread,
// and blocks until the callback is finished.
func (tw *ThreadWorker) InvokeSync(callback InvokeFunc) {
	done := make(chan struct{})
	tw.work <- func() {
		callback()
		done <- struct{}{}
	}
	<-done
}

// Aborts the worker, cancelling any pending work.
func (tw *ThreadWorker) Abort() {
	close(tw.work)
}

// Stop the worker and release any resources. Blocks until all work in completed.
func (tw *ThreadWorker) Stop() {
	tw.InvokeSync(func() {
		close(tw.work)
	})
}

// Flush blocks the caller until all pending work is completed
func (tw *ThreadWorker) Flush() {
	tw.InvokeSync(func() {})
}
