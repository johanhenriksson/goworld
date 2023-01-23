package swapchain

import (
	gosync "sync"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"
)

type Context struct {
	Index          int
	ImageAvailable sync.Semaphore
	RenderComplete sync.Semaphore
	InFlight       *gosync.Mutex
}

func newContext(dev device.T, index int) Context {
	return Context{
		Index:          index,
		ImageAvailable: sync.NewSemaphore(dev),
		RenderComplete: sync.NewSemaphore(dev),
		InFlight:       &gosync.Mutex{},
	}
}

func (c *Context) Destroy() {
	if c.ImageAvailable != nil {
		c.ImageAvailable.Destroy()
		c.ImageAvailable = nil
	}
	if c.RenderComplete != nil {
		c.RenderComplete.Destroy()
		c.RenderComplete = nil
	}
}
