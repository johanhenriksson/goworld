package swapchain

import (
	"fmt"
	gosync "sync"
	"time"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"
)

type Context struct {
	Index          int
	Start          time.Time
	ImageAvailable sync.Semaphore
	RenderComplete sync.Semaphore

	image    int
	inFlight *gosync.Mutex
}

func newContext(dev device.T, index int) *Context {
	return &Context{
		Index:          index,
		ImageAvailable: sync.NewSemaphore(dev, fmt.Sprintf("ImageAvailable:%d", index)),
		RenderComplete: sync.NewSemaphore(dev, fmt.Sprintf("RenderComplete:%d", index)),
		inFlight:       &gosync.Mutex{},
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

func (c *Context) Aquire() {
	c.inFlight.Lock()
}

func (c *Context) Release() {
	c.inFlight.Unlock()
}
