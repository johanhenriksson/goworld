package buffer

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Buffer]

	Size() int
	Read(data any, offset int)
	Write(data any, offset int)
}

type buffer struct {
	ptr    vk.Buffer
	device device.T
	memory device.Memory
	size   int
}

func New(device device.T, size int, usage vk.BufferUsageFlags, properties vk.MemoryPropertyFlags, sharing vk.SharingMode) T {
	queueIdx := device.GetQueueFamilyIndex(vk.QueueFlags(vk.QueueGraphicsBit))
	info := vk.BufferCreateInfo{
		SType:                 vk.StructureTypeBufferCreateInfo,
		Flags:                 vk.BufferCreateFlags(0),
		Size:                  vk.DeviceSize(size),
		Usage:                 usage,
		SharingMode:           sharing,
		QueueFamilyIndexCount: 1,
		PQueueFamilyIndices:   []uint32{uint32(queueIdx)},
	}

	var ptr vk.Buffer
	vk.CreateBuffer(device.Ptr(), &info, nil, &ptr)

	var memreq vk.MemoryRequirements
	vk.GetBufferMemoryRequirements(device.Ptr(), ptr, &memreq)
	memreq.Deref()

	mem := device.Allocate(memreq, properties)

	vk.BindBufferMemory(device.Ptr(), ptr, mem.Ptr(), 0)

	return &buffer{
		ptr:    ptr,
		device: device,
		memory: mem,
		size:   int(memreq.Size),
	}
}

func NewShared(device device.T, size int) T {
	return New(
		device, size,
		vk.BufferUsageFlags(vk.BufferUsageTransferSrcBit),
		vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit|vk.MemoryPropertyHostCoherentBit),
		vk.SharingModeExclusive)
}

func NewRemote(device device.T, size int, flags vk.BufferUsageFlags) T {
	return New(
		device, size,
		vk.BufferUsageFlags(vk.BufferUsageTransferDstBit)|flags,
		vk.MemoryPropertyFlags(vk.MemoryPropertyDeviceLocalBit),
		vk.SharingModeExclusive)
}

func (b *buffer) Ptr() vk.Buffer {
	return b.ptr
}

func (b *buffer) Size() int {
	return b.size
}

func (b *buffer) Destroy() {
	b.memory.Destroy()
	vk.DestroyBuffer(b.device.Ptr(), b.ptr, nil)
	b.ptr = nil
}

func (b *buffer) Write(data any, offset int) {
	b.memory.Write(data, offset)
}

func (b *buffer) Read(data any, offset int) {
	b.memory.Read(data, offset)
}
