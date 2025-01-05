package pipeline

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Compute struct {
	ptr    core1_0.Pipeline
	module shader.Module
}

func NewCompute(dev *device.Device, layout *Layout, path string) *Compute {
	// todo: we wanna use the shader cache probably, or do this somewhere else
	module := shader.NewModule(dev, assets.FS, fmt.Sprintf("shaders/%s.cs.glsl", path), shader.StageCompute)

	info := core1_0.ComputePipelineCreateInfo{
		Layout: layout.Ptr(),
		Stage: core1_0.PipelineShaderStageCreateInfo{
			Name:   "main",
			Module: module.Ptr(),
			Stage:  core1_0.ShaderStageFlags(module.Stage()),
		},
	}

	ptrs, _, err := dev.Ptr().CreateComputePipelines(nil, nil, []core1_0.ComputePipelineCreateInfo{info})
	if err != nil {
		panic(err)
	}

	dev.SetDebugObjectName(driver.VulkanHandle(ptrs[0].Handle()), core1_0.ObjectTypePipeline, path)

	return &Compute{
		ptr:    ptrs[0],
		module: module,
	}
}

func (c *Compute) Ptr() core1_0.Pipeline {
	return c.ptr
}

func (c *Compute) Destroy() {
	c.ptr.Destroy(nil)
	c.ptr = nil
	c.module.Destroy()
	c.module = nil
}
