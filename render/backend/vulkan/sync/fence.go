package sync

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Fence interface {
	device.Resource[vk.Fence]

	Reset()
	Wait()
}

type fence struct {
	device device.T
	ptr    vk.Fence
}

func NewFence(device device.T, signaled bool) Fence {
	var flags vk.FenceCreateFlags
	if signaled {
		flags = vk.FenceCreateFlags(vk.FenceCreateSignaledBit)
	}

	info := vk.FenceCreateInfo{
		SType: vk.StructureTypeFenceCreateInfo,
		Flags: flags,
	}

	var fnc vk.Fence
	vk.CreateFence(device.Ptr(), &info, nil, &fnc)

	return &fence{
		device: device,
		ptr:    fnc,
	}
}

func (f *fence) Ptr() vk.Fence {
	return f.ptr
}

func (f *fence) Reset() {
	vk.ResetFences(f.device.Ptr(), 1, []vk.Fence{f.ptr})
}

func (f *fence) Destroy() {
	vk.DestroyFence(f.device.Ptr(), f.ptr, nil)
	f.ptr = nil
}

func (f *fence) Wait() {
	vk.WaitForFences(f.device.Ptr(), 1, []vk.Fence{f.ptr}, vk.Bool32(0), vk.MaxUint64)
}
