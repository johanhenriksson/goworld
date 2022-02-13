package pipeline

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	vk "github.com/vulkan-go/vulkan"
)

type Pass interface {
	device.Resource[vk.RenderPass]
}

type pass struct {
	device device.T
	ptr    vk.RenderPass
}

func NewPass(device device.T, args *vk.RenderPassCreateInfo) Pass {
	var ptr vk.RenderPass
	vk.CreateRenderPass(device.Ptr(), args, nil, &ptr)

	return &pass{
		ptr:    ptr,
		device: device,
	}
}

func (p *pass) Ptr() vk.RenderPass {
	return p.ptr
}

func (p *pass) Destroy() {
	vk.DestroyRenderPass(p.device.Ptr(), p.ptr, nil)
	p.ptr = nil
}
