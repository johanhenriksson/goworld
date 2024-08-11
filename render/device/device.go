package device

import (
	"log"
	"slices"

	"github.com/johanhenriksson/goworld/render/instance"
	"github.com/samber/lo"

	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
	"github.com/vkngwrapper/extensions/v2/ext_debug_utils"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Resource[T any] interface {
	Destroy()
	Ptr() T
}

type Device struct {
	physical core1_0.PhysicalDevice
	ptr      core1_0.Device
	limits   *core1_0.PhysicalDeviceLimits
	debug    ext_debug_utils.Extension

	memtypes map[memtype]int
	queue    Queue
}

func New(instance *instance.Instance, physDevice core1_0.PhysicalDevice) (*Device, error) {
	log.Println("creating device with extensions", deviceExtensions)

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

	indexingFeatures := ext_descriptor_indexing.PhysicalDeviceDescriptorIndexingFeatures{
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
		NextOptions:           common.NextOptions{Next: indexingFeatures},
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

			MultiDrawIndirect:         true,
			DrawIndirectFirstInstance: true,
		},
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

	return &Device{
		ptr:      dev,
		debug:    debug,
		physical: physDevice,
		limits:   properties.Limits,
		memtypes: make(map[memtype]int),
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

func (d *Device) GetMemoryTypeIndex(typeBits uint32, flags core1_0.MemoryPropertyFlags) int {
	mtype := memtype{typeBits, flags}
	if t, ok := d.memtypes[mtype]; ok {
		return t
	}

	props := d.physical.MemoryProperties()
	for i, kind := range props.MemoryTypes {
		if typeBits&1 == 1 {
			if kind.PropertyFlags&flags == flags {
				d.memtypes[mtype] = i
				return i
			}
		}
		typeBits >>= 1
	}

	d.memtypes[mtype] = 0
	return 0
}

func (d *Device) GetLimits() *core1_0.PhysicalDeviceLimits {
	return d.limits
}

func (d *Device) Allocate(key string, req core1_0.MemoryRequirements, flags core1_0.MemoryPropertyFlags) Memory {
	if req.Size == 0 {
		panic("allocating 0 bytes of memory")
	}
	return alloc(d, key, req, flags)
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
