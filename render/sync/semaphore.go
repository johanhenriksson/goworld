package sync

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Semaphore interface {
	device.Resource[core1_0.Semaphore]
}

type semaphore struct {
	device device.T
	ptr    core1_0.Semaphore
}

func NewSemaphore(dev device.T, name string) Semaphore {
	ptr, _, err := dev.Ptr().CreateSemaphore(nil, core1_0.SemaphoreCreateInfo{})
	if err != nil {
		panic(err)
	}
	dev.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()), core1_0.ObjectTypeSemaphore, name)

	return &semaphore{
		device: dev,
		ptr:    ptr,
	}
}

func (s semaphore) Ptr() core1_0.Semaphore {
	return s.ptr
}

func (s *semaphore) Destroy() {
	s.ptr.Destroy(nil)
	s.ptr = nil
}

func NewSemaphoreArray(dev device.T, name string, count int) []Semaphore {
	arr := make([]Semaphore, count)
	for i := range arr {
		arr[i] = NewSemaphore(dev, fmt.Sprintf("%s:%d", name, count))
	}
	return arr
}
