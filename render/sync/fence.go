package sync

import (
	"time"

	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Fence interface {
	device.Resource[core1_0.Fence]

	Reset()
	Wait()
	Done() bool
}

type fence struct {
	device device.T
	ptr    core1_0.Fence
}

func NewFence(device device.T, signaled bool) Fence {
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

	return &fence{
		device: device,
		ptr:    ptr,
	}
}

func (f *fence) Ptr() core1_0.Fence {
	return f.ptr
}

func (f *fence) Reset() {
	f.device.Ptr().ResetFences([]core1_0.Fence{f.ptr})
}

func (f *fence) Destroy() {
	f.ptr.Destroy(nil)
	f.ptr = nil
}

func (f *fence) Wait() {
	f.device.Ptr().WaitForFences(true, time.Hour, []core1_0.Fence{f.ptr})
}

func (f *fence) Done() bool {
	r, err := f.ptr.Status()
	if err != nil {
		panic(err)
	}
	return r == core1_0.VKSuccess
}
