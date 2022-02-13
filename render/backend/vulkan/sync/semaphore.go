package sync

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Semaphore interface {
	device.Resource[vk.Semaphore]
}

type semaphore struct {
	device device.T
	ptr    vk.Semaphore
}

func NewSemaphore(dev device.T) Semaphore {
	info := vk.SemaphoreCreateInfo{
		SType: vk.StructureTypeSemaphoreCreateInfo,
	}

	var sem vk.Semaphore
	vk.CreateSemaphore(dev.Ptr(), &info, nil, &sem)

	return &semaphore{
		device: dev,
		ptr:    sem,
	}
}

func (s semaphore) Ptr() vk.Semaphore {
	return s.ptr
}

func (s semaphore) Destroy() {
	vk.DestroySemaphore(s.device.Ptr(), s.ptr, nil)
	s.ptr = nil
}
