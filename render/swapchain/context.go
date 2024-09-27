package swapchain

import (
	"fmt"
	stdsync "sync"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"
)

type Context struct {
	Index          int
	ImageAvailable *sync.Semaphore
	RenderComplete *sync.Semaphore
	ready          stdsync.Mutex
}

func NewContext(dev *device.Device, index int) *Context {
	return &Context{
		Index:          index,
		ImageAvailable: sync.NewSemaphore(dev, fmt.Sprintf("ImageAvailable:%d", index)),
		RenderComplete: sync.NewSemaphore(dev, fmt.Sprintf("RenderComplete:%d", index)),
	}
}

func (c *Context) Aquire() {
	c.ready.Lock()
}

func (c *Context) Release() {
	c.ready.Unlock()
}

func (c *Context) Destroy() {
	c.ImageAvailable.Destroy()
	c.RenderComplete.Destroy()
}
