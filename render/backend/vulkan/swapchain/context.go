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
	c.ColorView.Destroy()
	c.Framebuffer.Destroy()
	c.Workers[0].Destroy()
}
