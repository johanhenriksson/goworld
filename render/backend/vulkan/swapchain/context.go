package swapchain

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
)

type Context struct {
	Index       int
	Color       image.T
	Depth       image.T
	ColorView   image.View
	DepthView   image.View
	Framebuffer framebuffer.T
	Workers     command.Workers
	Output      pipeline.Pass
	Width       int
	Height      int
}

func (c Context) Destroy() {
	if c.ColorView != nil {

		c.ColorView.Destroy()
	}
	if c.Framebuffer != nil {

		c.Framebuffer.Destroy()
	}
	if c.Workers != nil {
		for _, worker := range c.Workers {
			worker.Destroy()
		}
	}
}
