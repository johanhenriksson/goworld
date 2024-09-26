package device

import (
	"fmt"
	"log"
	"slices"

	"github.com/johanhenriksson/goworld/render/instance"
	"github.com/johanhenriksson/goworld/util"
	"github.com/samber/lo"

	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/core1_1"
	"github.com/vkngwrapper/core/v2/core1_2"
	"github.com/vkngwrapper/core/v2/driver"
	"github.com/vkngwrapper/extensions/v2/ext_debug_utils"
	"github.com/vkngwrapper/extensions/v2/khr_buffer_device_address"
)

type Resource[T any] interface {
	Destroy()
	Ptr() T
}

type Device struct {
	physical core1_0.PhysicalDevice
	ptr      core1_0.Device
	limits   core1_0.PhysicalDeviceLimits

	debug   ext_debug_utils.Extension
	address khr_buffer_device_address.Extension

	queue    Queue
	memtypes [5]memtype
}

func New(instance *instance.Instance, physDevice core1_0.PhysicalDevice) (*Device, error) {
	log.Println("creating device with extensions", deviceExtensions)

	//
	// find a suitable queue
	//

	queues := []Queue{}
	families := physDevice.QueueFamilyProperties()
	log.Println("Queue families:", len(families))
	for index, family := range families {
		log.Printf("  %d: %dx %s\n", index, family.QueueCount, family.QueueFlags)
		for i := 0; i < int(family.QueueCount); i++ {
			queues = append(queues, Queue{
				flags:  family.QueueFlags,
				family: index,
				index:  i,
			})
		}
	}

	mostSpecificQueue := func(flags core1_0.QueueFlags, avoid ...Queue) Queue {
		options := lo.Filter(queues, func(q Queue, _ int) bool { return q.Matches(flags) })

		// try to avoid certain families
		optimal := lo.Filter(options, func(q Queue, _ int) bool { return !slices.Contains(avoid, q) })
		if len(optimal) > 0 {
			options = optimal
		}

		return lo.MinBy(options, func(a Queue, b Queue) bool { return int(a.flags) < int(b.flags) })
	}

	queue := mostSpecificQueue(core1_0.QueueGraphics | core1_0.QueueTransfer)
	log.Println("worker queue:", queue)

	core12Features := core1_2.PhysicalDeviceVulkan12Features{
		BufferDeviceAddress: true,

		ScalarBlockLayout: true,

		ShaderInt8:                        true,
		StorageBuffer8BitAccess:           true,
		UniformAndStorageBuffer8BitAccess: true,
		StoragePushConstant8:              true,

		DescriptorIndexing:                                 true,
		ShaderSampledImageArrayNonUniformIndexing:          true,
		RuntimeDescriptorArray:                             true,
		DescriptorBindingPartiallyBound:                    true,
		DescriptorBindingVariableDescriptorCount:           true,
		DescriptorBindingUpdateUnusedWhilePending:          true,
		DescriptorBindingUniformBufferUpdateAfterBind:      true,
		DescriptorBindingSampledImageUpdateAfterBind:       true,
		DescriptorBindingStorageBufferUpdateAfterBind:      true,
		DescriptorBindingStorageTexelBufferUpdateAfterBind: true,
	}

	dev, _, err := physDevice.CreateDevice(nil, core1_0.DeviceCreateInfo{
		EnabledExtensionNames: deviceExtensions,
		QueueCreateInfos: []core1_0.DeviceQueueCreateInfo{
			{
				QueueFamilyIndex: queue.FamilyIndex(),
				QueuePriorities:  []float32{1},
			},
		},
		EnabledFeatures: &core1_0.PhysicalDeviceFeatures{
			IndependentBlend: true,
			DepthClamp:       true,

			ShaderInt16: true,
			ShaderInt64: true,

			MultiDrawIndirect:         true,
			DrawIndirectFirstInstance: true,
		},
		NextOptions: common.NextOptions{Next: core12Features},
	})
	if err != nil {
		return nil, err
	}

	properties, err := physDevice.Properties()
	if err != nil {
		return nil, err
	}
	log.Println("minimum uniform buffer alignment:", properties.Limits.MinUniformBufferOffsetAlignment)
	log.Println("minimum storage buffer alignment:", properties.Limits.MinStorageBufferOffsetAlignment)

	debug := ext_debug_utils.CreateExtensionFromInstance(instance.Ptr())

	// resolve queue pointers
	queue.ptr = dev.GetQueue(queue.FamilyIndex(), queue.Index())
	debug.SetDebugUtilsObjectName(dev, ext_debug_utils.DebugUtilsObjectNameInfo{
		ObjectName:   "graphics",
		ObjectHandle: driver.VulkanHandle(queue.Ptr().Handle()),
		ObjectType:   core1_0.ObjectTypeQueue,
	})

	address := khr_buffer_device_address.CreateExtensionFromDevice(dev)

	//
	// resolve memory types
	//
	imageMemoryTypeBits, err := getImageMemoryTypeBits(dev)
	if err != nil {
		return nil, fmt.Errorf("failed to get image memory type bits: %w", err)
	}
	log.Println("image memory types:", imageMemoryTypeBits)

	bufferMemoryTypeBits, err := getBufferMemoryTypeBits(dev)
	if err != nil {
		return nil, fmt.Errorf("failed to get buffer memory type bits: %w", err)
	}
	log.Println("buffer memory types:", bufferMemoryTypeBits)

	memtypes := [5]memtype{}
	memtypes[0] = memtype{-1, 0}
	memoryProperties := physDevice.MemoryProperties()
	for i, kind := range memoryProperties.MemoryTypes {
		log.Println("memory type", i, ":", kind.PropertyFlags)
	}

	// gpu local memory
	memtypes[MemoryTypeGPU] = pickPreferredMemoryType(memoryProperties.MemoryTypes, bufferMemoryTypeBits,
		core1_0.MemoryPropertyDeviceLocal, 0)
	if memtypes[MemoryTypeGPU].Index == -1 {
		return nil, fmt.Errorf("failed to find gpu local memory type")
	}
	log.Println("gpu local memory type:",
		memtypes[MemoryTypeGPU],
		memoryProperties.MemoryTypes[memtypes[MemoryTypeGPU].Index].PropertyFlags)

	// shared memory
	memtypes[MemoryTypeShared] = pickPreferredMemoryType(memoryProperties.MemoryTypes, bufferMemoryTypeBits,
		core1_0.MemoryPropertyDeviceLocal|
			core1_0.MemoryPropertyHostVisible, core1_0.MemoryPropertyHostCoherent)
	if memtypes[MemoryTypeShared].Index == -1 {
		return nil, fmt.Errorf("failed to find shared memory type")
	}
	log.Println("shared memory type:",
		memtypes[MemoryTypeShared],
		memoryProperties.MemoryTypes[memtypes[MemoryTypeShared].Index].PropertyFlags)

	// cpu local memory
	memtypes[MemoryTypeCPU] = pickPreferredMemoryType(memoryProperties.MemoryTypes, bufferMemoryTypeBits,
		core1_0.MemoryPropertyHostVisible, 0)
	if memtypes[MemoryTypeCPU].Index == -1 {
		return nil, fmt.Errorf("failed to find cpu local memory type")
	}
	log.Println("cpu local memory type:",
		memtypes[MemoryTypeCPU],
		memoryProperties.MemoryTypes[memtypes[MemoryTypeCPU].Index].PropertyFlags)

	// texture memory
	memtypes[MemoryTypeTexture] = pickPreferredMemoryType(memoryProperties.MemoryTypes, imageMemoryTypeBits,
		core1_0.MemoryPropertyDeviceLocal, 0)
	if memtypes[MemoryTypeTexture].Index == -1 {
		return nil, fmt.Errorf("failed to find texture memory type")
	}
	log.Println("texture memory type:",
		memtypes[MemoryTypeTexture],
		memoryProperties.MemoryTypes[memtypes[MemoryTypeTexture].Index].PropertyFlags)

	return &Device{
		ptr:      dev,
		debug:    debug,
		address:  address,
		physical: physDevice,
		limits:   *properties.Limits,
		memtypes: memtypes,
		queue:    queue,
	}, nil

}

func (d *Device) Ptr() core1_0.Device {
	return d.ptr
}

func (d *Device) Physical() core1_0.PhysicalDevice {
	return d.physical
}

func (d *Device) Queue() Queue {
	return d.queue
}

func (d *Device) GetFormatProperties(format core1_0.Format) *core1_0.FormatProperties {
	return d.physical.FormatProperties(format)
}

func (d *Device) GetDepthFormat() core1_0.Format {
	depthFormats := []core1_0.Format{
		core1_0.FormatD32SignedFloatS8UnsignedInt,
		core1_0.FormatD32SignedFloat,
		core1_0.FormatD24UnsignedNormalizedS8UnsignedInt,
		core1_0.FormatD16UnsignedNormalizedS8UnsignedInt,
		core1_0.FormatD16UnsignedNormalized,
	}
	for _, format := range depthFormats {
		props := d.physical.FormatProperties(format)

		if props.OptimalTilingFeatures&core1_0.FormatFeatureDepthStencilAttachment == core1_0.FormatFeatureDepthStencilAttachment {
			return format
		}
	}
	return depthFormats[0]
}

func (d *Device) GetLimits() core1_0.PhysicalDeviceLimits {
	return d.limits
}

func (d *Device) Allocate(key string, size int, kind MemoryType) *Memory {
	if size == 0 {
		panic("allocating 0 bytes of memory")
	}
	if int(kind) > len(d.memtypes) || int(kind) <= 0 {
		panic(fmt.Sprintf("invalid memory type %d", kind))
	}
	mtype := d.memtypes[kind]
	if mtype.Index == -1 {
		panic(fmt.Sprintf("memory type %d is not resolved", kind))
	}

	align := int(d.GetLimits().NonCoherentAtomSize)
	size = util.Align(int(size), align)

	ptr, _, err := d.ptr.AllocateMemory(nil, core1_0.MemoryAllocateInfo{
		AllocationSize:  size,
		MemoryTypeIndex: mtype.Index,
		NextOptions: common.NextOptions{
			Next: core1_1.MemoryAllocateFlagsInfo{
				Flags: khr_buffer_device_address.MemoryAllocateDeviceAddress,
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to allocate %d bytes of memory: %s", size, err))
	}

	if key != "" {
		d.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()),
			core1_0.ObjectTypeDeviceMemory, key)
	}

	return &Memory{
		device: d,
		ptr:    ptr,
		flags:  mtype.Flags,
		size:   size,
	}
}

func (d *Device) Destroy() {
	d.ptr.Destroy(nil)
	d.ptr = nil
}

func (d *Device) WaitIdle() {
	d.ptr.WaitIdle()
}

func (d *Device) SetDebugObjectName(handle driver.VulkanHandle, objType core1_0.ObjectType, name string) {
	d.debug.SetDebugUtilsObjectName(d.ptr, ext_debug_utils.DebugUtilsObjectNameInfo{
		ObjectName:   name,
		ObjectHandle: handle,
		ObjectType:   objType,
	})
}

func (d *Device) GetBufferAddress(buffer core1_0.Buffer) (Address, error) {
	addr, err := d.address.GetBufferDeviceAddress(d.ptr, khr_buffer_device_address.BufferDeviceAddressInfo{
		Buffer: buffer,
	})
	if err != nil {
		return 0, err
	}
	return Address(addr), nil
}
