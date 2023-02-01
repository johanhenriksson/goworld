package sync

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Semaphore interface {
	device.Resource[core1_0.Semaphore]
}

type semaphore struct {
	device device.T
	ptr    core1_0.Semaphore
}

func NewSemaphore(dev device.T) Semaphore {
	ptr, _, err := dev.Ptr().CreateSemaphore(nil, core1_0.SemaphoreCreateInfo{})
	if err != nil {
		panic(err)
	}

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

func NewSemaphoreArray(dev device.T, count int) []Semaphore {
	arr := make([]Semaphore, count)
	for i := range arr {
		arr[i] = NewSemaphore(dev)
	}
	return arr
}
