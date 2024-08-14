package sync

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Semaphore struct {
	device *device.Device
	ptr    core1_0.Semaphore
	name   string
}

func NewSemaphore(dev *device.Device, name string) *Semaphore {
	ptr, _, err := dev.Ptr().CreateSemaphore(nil, core1_0.SemaphoreCreateInfo{})
	if err != nil {
		panic(err)
	}
	dev.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()), core1_0.ObjectTypeSemaphore, name)

	return &Semaphore{
		device: dev,
		ptr:    ptr,
		name:   name,
	}
}

func (s Semaphore) Ptr() core1_0.Semaphore {
	return s.ptr
}

func (s *Semaphore) Name() string {
	return s.name
}

func (s *Semaphore) Destroy() {
	if s.ptr != nil {
		s.ptr.Destroy(nil)
		s.ptr = nil
	}
}

func NewSemaphoreArray(dev *device.Device, name string, count int) []*Semaphore {
	arr := make([]*Semaphore, count)
	for i := range arr {
		arr[i] = NewSemaphore(dev, fmt.Sprintf("%s:%d", name, i))
	}
	return arr
}
