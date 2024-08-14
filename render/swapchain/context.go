package swapchain

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"
)

type Context struct {
	Index          int
	ImageAvailable *sync.Semaphore
	RenderComplete *sync.Semaphore
}

func NewContext(dev *device.Device, index int) *Context {
	return &Context{
		Index:          index,
		ImageAvailable: sync.NewSemaphore(dev, fmt.Sprintf("ImageAvailable:%d", index)),
		RenderComplete: sync.NewSemaphore(dev, fmt.Sprintf("RenderComplete:%d", index)),
	}
}

func (c *Context) Destroy() {
	c.ImageAvailable.Destroy()
	c.RenderComplete.Destroy()
}
