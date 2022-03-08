package device

import (
	"fmt"

	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Resource[T any] interface {
	Destroy()
	Ptr() T
}

type T interface {
	Resource[vk.Device]

	Allocate(vk.MemoryRequirements, vk.MemoryPropertyFlags) Memory
	GetQueue(queueIndex int, flags vk.QueueFlags) vk.Queue
	GetQueueFamilyIndex(flags vk.QueueFlags) int
	GetSurfaceFormats(vk.Surface) []vk.SurfaceFormat
	GetSurfaceCapabilities(surface vk.Surface) *vk.SurfaceCapabilities
	GetDepthFormat() vk.Format
	GetMemoryTypeIndex(uint32, vk.MemoryPropertyFlags) int
	WaitIdle()
}

type device struct {
	physical vk.PhysicalDevice
	ptr      vk.Device

	memtypes map[memtype]int
	queues   map[vk.QueueFlags]int
}

func New(physDevice vk.PhysicalDevice) (T, error) {
	queueInfo := vk.DeviceQueueCreateInfo{
		SType:            vk.StructureTypeDeviceQueueCreateInfo,
		QueueFamilyIndex: 0,
		QueueCount:       1,
		PQueuePriorities: []float32{1},
	}

	var dev vk.Device
	deviceInfo := vk.DeviceCreateInfo{
		SType:                   vk.StructureTypeDeviceCreateInfo,
		EnabledExtensionCount:   uint32(len(deviceExtensions)),
		PpEnabledExtensionNames: util.CStrings(deviceExtensions),
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
		memtypes: make(map[memtype]int),
		queues:   make(map[vk.QueueFlags]int),
	}, nil
}

func (d *device) Ptr() vk.Device {
	return d.ptr
}

func (d *device) GetQueue(queueIndex int, flags vk.QueueFlags) vk.Queue {
	familyIndex := d.GetQueueFamilyIndex(flags)
	var queue vk.Queue
	vk.GetDeviceQueue(d.ptr, uint32(familyIndex), uint32(queueIndex), &queue)
	return queue
}

func (d *device) GetQueueFamilyIndex(flags vk.QueueFlags) int {
	if q, ok := d.queues[flags]; ok {
		return q
	}

	var familyCount uint32
	vk.GetPhysicalDeviceQueueFamilyProperties(d.physical, &familyCount, nil)
	families := make([]vk.QueueFamilyProperties, uint32(familyCount))
	vk.GetPhysicalDeviceQueueFamilyProperties(d.physical, &familyCount, families)

	for index, family := range families {
		family.Deref()
		if family.QueueFlags&flags == flags {
			d.queues[flags] = index
			return index
		}
	}

	d.queues[flags] = 0
	return 0
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

func (d *device) GetSurfaceCapabilities(surface vk.Surface) *vk.SurfaceCapabilities {
	var caps vk.SurfaceCapabilities
	vk.GetPhysicalDeviceSurfaceCapabilities(d.physical, surface, &caps)
	caps.Deref()
	return &caps
}

func (d *device) GetDepthFormat() vk.Format {
	depthFormats := []vk.Format{
		vk.FormatD32SfloatS8Uint,
		vk.FormatD32Sfloat,
		vk.FormatD24UnormS8Uint,
		vk.FormatD16UnormS8Uint,
		vk.FormatD16Unorm,
	}
	for _, format := range depthFormats {
		var properties vk.FormatProperties
		vk.GetPhysicalDeviceFormatProperties(d.physical, format, &properties)
		properties.Deref()

		if properties.OptimalTilingFeatures&vk.FormatFeatureFlags(vk.FormatFeatureDepthStencilAttachmentBit) == vk.FormatFeatureFlags(vk.FormatFeatureDepthStencilAttachmentBit) {
			return format
		}
	}

	return depthFormats[0]
}

func (d *device) GetMemoryTypeIndex(typeBits uint32, flags vk.MemoryPropertyFlags) int {
	mtype := memtype{typeBits, flags}
	if t, ok := d.memtypes[mtype]; ok {
		return t
	}

	var props vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(d.physical, &props)
	props.Deref()

	for i := 0; i < int(props.MemoryTypeCount); i++ {
		if typeBits&1 == 1 {
			props.MemoryTypes[i].Deref()
			if props.MemoryTypes[i].PropertyFlags&flags == flags {
				d.memtypes[mtype] = i
				return i
			}
		}
		typeBits >>= 1
	}

	d.memtypes[mtype] = 0
	return 0
}

func (d *device) Allocate(req vk.MemoryRequirements, flags vk.MemoryPropertyFlags) Memory {
	return alloc(d, req, flags)
}

func (d *device) Destroy() {
	if d.ptr != nil {
		vk.DestroyDevice(d.ptr, nil)
		d.ptr = nil
	}
}

func (d *device) WaitIdle() {
	vk.DeviceWaitIdle(d.ptr)
}
