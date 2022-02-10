package device

import (
	"fmt"

	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Resource interface {
	Destroy()
}

type T interface {
	Resource
	Ptr() vk.Device

	GetQueue(queueFamily, queueIndex int) vk.Queue
	GetSurfaceFormats(vk.Surface) []vk.SurfaceFormat
	WaitIdle()
}

type device struct {
	physical vk.PhysicalDevice
	ptr      vk.Device
}

func New(physDevice vk.PhysicalDevice) (T, error) {
	queueInfo := vk.DeviceQueueCreateInfo{
		SType:            vk.StructureTypeDeviceQueueCreateInfo,
		QueueFamilyIndex: 0,
		QueueCount:       1,
		PQueuePriorities: []float32{1},
	}

	var deviceExtensions = util.CStrings([]string{
		"VK_KHR_swapchain",
		"VK_KHR_portability_subset",
	})

	var dev vk.Device
	deviceInfo := vk.DeviceCreateInfo{
		SType:                   vk.StructureTypeDeviceCreateInfo,
		EnabledExtensionCount:   uint32(len(deviceExtensions)),
		PpEnabledExtensionNames: deviceExtensions,
		PQueueCreateInfos:       []vk.DeviceQueueCreateInfo{queueInfo},
		QueueCreateInfoCount:    1,
	}

	r := vk.CreateDevice(physDevice, &deviceInfo, nil, &dev)
	if r != vk.Success {
		return nil, fmt.Errorf("failed to create logical device")
	}

	return &device{
		physical: physDevice,
		ptr:      dev,
	}, nil
}

func (d *device) Ptr() vk.Device {
	return d.ptr
}

func (d *device) GetQueue(queueFamily, queueIndex int) vk.Queue {
	var queue vk.Queue
	vk.GetDeviceQueue(d.ptr, uint32(queueFamily), uint32(queueIndex), &queue)
	return queue
}

func (d *device) GetSurfaceFormats(surface vk.Surface) []vk.SurfaceFormat {
	surfaceFormatCount := uint32(0)
	vk.GetPhysicalDeviceSurfaceFormats(d.physical, surface, &surfaceFormatCount, nil)
	surfaceFormats := make([]vk.SurfaceFormat, surfaceFormatCount)
	vk.GetPhysicalDeviceSurfaceFormats(d.physical, surface, &surfaceFormatCount, surfaceFormats)
	for i, format := range surfaceFormats {
		format.Deref()
		surfaceFormats[i] = format
	}
	return surfaceFormats
}

func (d *device) Destroy() {
	vk.DestroyDevice(d.ptr, nil)
	d.ptr = nil
}

func (d *device) WaitIdle() {
	vk.DeviceWaitIdle(d.ptr)
}
