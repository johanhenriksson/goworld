package device

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

var Nil = &device{
	ptr: vk.Device(vk.NullHandle),
}

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
	GetLimits() *vk.PhysicalDeviceLimits
	WaitIdle()
}

type device struct {
	physical vk.PhysicalDevice
	ptr      vk.Device
	limits   vk.PhysicalDeviceLimits

	memtypes map[memtype]int
	queues   map[vk.QueueFlags]int
}

func New(physDevice vk.PhysicalDevice) (T, error) {
	log.Println("creating device with extensions", deviceExtensions)

	// VK_EXT_descriptor_indexing settings
	indexingFeatures := vk.PhysicalDeviceDescriptorIndexingFeatures{
		SType: vk.StructureTypePhysicalDeviceDescriptorIndexingFeatures,
		ShaderSampledImageArrayNonUniformIndexing:          vk.True,
		RuntimeDescriptorArray:                             vk.True,
		DescriptorBindingPartiallyBound:                    vk.True,
		DescriptorBindingVariableDescriptorCount:           vk.True,
		DescriptorBindingUpdateUnusedWhilePending:          vk.True,
		DescriptorBindingUniformBufferUpdateAfterBind:      vk.True,
		DescriptorBindingSampledImageUpdateAfterBind:       vk.True,
		DescriptorBindingStorageBufferUpdateAfterBind:      vk.True,
		DescriptorBindingStorageTexelBufferUpdateAfterBind: vk.True,
	}

	var dev vk.Device
	deviceInfo := vk.DeviceCreateInfo{
		SType:                   vk.StructureTypeDeviceCreateInfo,
		PNext:                   unsafe.Pointer(&indexingFeatures),
		EnabledExtensionCount:   uint32(len(deviceExtensions)),
		PpEnabledExtensionNames: util.CStrings(deviceExtensions),
		PQueueCreateInfos: []vk.DeviceQueueCreateInfo{
			{
				SType:            vk.StructureTypeDeviceQueueCreateInfo,
				QueueFamilyIndex: 0,
				QueueCount:       1,
				PQueuePriorities: []float32{1},
			},
			{
				SType:            vk.StructureTypeDeviceQueueCreateInfo,
				QueueFamilyIndex: 1,
				QueueCount:       1,
				PQueuePriorities: []float32{1},
			},
			{
				SType:            vk.StructureTypeDeviceQueueCreateInfo,
				QueueFamilyIndex: 2,
				QueueCount:       1,
				PQueuePriorities: []float32{1},
			},
			{
				SType:            vk.StructureTypeDeviceQueueCreateInfo,
				QueueFamilyIndex: 3,
				QueueCount:       1,
				PQueuePriorities: []float32{1},
			},
		},
		QueueCreateInfoCount: 4,
	}

	r := vk.CreateDevice(physDevice, &deviceInfo, nil, &dev)
	if r != vk.Success {
		return nil, fmt.Errorf("failed to create logical device")
	}

	var properties vk.PhysicalDeviceProperties
	vk.GetPhysicalDeviceProperties(physDevice, &properties)
	properties.Deref()
	properties.Limits.Deref()

	log.Println("minimum uniform buffer alignment:", properties.Limits.MinUniformBufferOffsetAlignment)
	log.Println("minimum storage buffer alignment:", properties.Limits.MinStorageBufferOffsetAlignment)

	return &device{
		physical: physDevice,
		ptr:      dev,
		limits:   properties.Limits,
		memtypes: make(map[memtype]int),
		queues:   make(map[vk.QueueFlags]int),
	}, nil
}

func (d *device) Ptr() vk.Device {
	return d.ptr
}

func (d *device) GetQueue(queueIndex int, flags vk.QueueFlags) vk.Queue {
	// familyIndex := d.GetQueueFamilyIndex(flags)
	var queue vk.Queue
	vk.GetDeviceQueue(d.ptr, uint32(queueIndex), 0, &queue)
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

	panic("no such queue available")
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

func (d *device) GetLimits() *vk.PhysicalDeviceLimits {
	return &d.limits
}

func (d *device) Allocate(req vk.MemoryRequirements, flags vk.MemoryPropertyFlags) Memory {
	return alloc(d, req, flags)
}

func (d *device) Destroy() {
	vk.DestroyDevice(d.ptr, nil)
	d.ptr = vk.Device(vk.NullHandle)
}

func (d *device) WaitIdle() {
	vk.DeviceWaitIdle(d.ptr)
}
