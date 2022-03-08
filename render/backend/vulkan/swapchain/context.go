package swapchain

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
)

type Context struct {
	Index          int
	Workers        command.Workers
	Width          int
	Height         int
	ImageAvailable sync.Semaphore
	RenderComplete sync.Semaphore
}

func (c Context) Destroy() {
	if c.ImageAvailable != nil {
		c.ImageAvailable.Destroy()
	}
	if c.RenderComplete != nil {
		c.RenderComplete.Destroy()
	}
	if c.Workers != nil {
		for _, worker := range c.Workers {
			worker.Destroy()
		}
	}
}
