package sync

import (
	"time"

	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Fence struct {
	device *device.Device
	ptr    core1_0.Fence
}

func NewFence(device *device.Device, name string, signaled bool) *Fence {
	var flags core1_0.FenceCreateFlags
	if signaled {
		flags = core1_0.FenceCreateSignaled
	}

	ptr, _, err := device.Ptr().CreateFence(nil, core1_0.FenceCreateInfo{
		Flags: flags,
	})
	if err != nil {
		panic(err)
	}
	device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()), core1_0.ObjectTypeFence, name)

	return &Fence{
		device: device,
		ptr:    ptr,
	}
}

func (f *Fence) Ptr() core1_0.Fence {
	return f.ptr
}

func (f *Fence) Reset() {
	f.device.Ptr().ResetFences([]core1_0.Fence{f.ptr})
}

func (f *Fence) Destroy() {
	f.ptr.Destroy(nil)
	f.ptr = nil
}

func (f *Fence) Wait() {
	f.device.Ptr().WaitForFences(true, time.Hour, []core1_0.Fence{f.ptr})
}

func (f *Fence) Done() bool {
	r, err := f.ptr.Status()
	if err != nil {
		panic(err)
	}
	return r == core1_0.VKSuccess
}
