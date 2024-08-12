package command

type Channel struct {
	buffer int
	work   chan InvokeFunc
}

func NewChannel(buffer int) *Channel {
	return &Channel{
		buffer: buffer,
		work:   make(chan InvokeFunc, buffer),
	}
}

func (ch *Channel) Recv() <-chan InvokeFunc {
	return ch.work
}

// Invoke schedules a callback to be called from the worker thread
func (ch *Channel) Invoke(callback InvokeFunc) {
	ch.work <- callback
}

// InvokeSync schedules a callback to be called on the worker thread,
// and blocks until the callback is finished.
func (ch *Channel) InvokeSync(callback InvokeFunc) {
	done := make(chan struct{})
	ch.Invoke(func() {
		callback()
		close(done)
	})
	<-done
}

// Aborts the worker, cancelling any pending work.
func (ch *Channel) Close() {
	close(ch.work)
}

// Stop the worker and release any resources. Blocks until all work in completed.
func (ch *Channel) Stop() {
	ch.InvokeSync(func() {
		close(ch.work)
	})
}

// Flush blocks the caller until all pending work is completed
func (ch *Channel) Flush() {
	ch.InvokeSync(func() {})
}
