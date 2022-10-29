package swapchain

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
)

type Context struct {
	Index          int
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
}
