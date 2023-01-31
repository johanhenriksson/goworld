package buffer

import (
	"github.com/johanhenriksson/goworld/render/device"

	vk "github.com/vulkan-go/vulkan"
)

var Nil T = &buffer{
	ptr:    vk.NullBuffer,
	device: device.Nil,
	memory: device.NilMemory,
}

type T interface {
	device.Resource[vk.Buffer]

	// Size returns the total allocation size of the buffer in bytes
	Size() int

	// Read directly from the buffer at the given offset
	Read(offset int, data any) int

	// Write directly to the buffer at the given offset
	Write(offset int, data any) int

	// Memory returns a handle to the underlying memory block
	Memory() device.Memory
}

type Args struct {
	Size   int
	Usage  vk.BufferUsageFlagBits
	Memory vk.MemoryPropertyFlagBits
}

type buffer struct {
	ptr    vk.Buffer
	device device.T
	memory device.Memory
	size   int
}

func New(device device.T, args Args) T {
	queueIdx := device.GetQueueFamilyIndex(vk.QueueFlags(vk.QueueGraphicsBit))
	info := vk.BufferCreateInfo{
		SType:       vk.StructureTypeBufferCreateInfo,
		Flags:       vk.BufferCreateFlags(0),
		Size:        vk.DeviceSize(args.Size),
		Usage:       vk.BufferUsageFlags(args.Usage),
		SharingMode: vk.SharingModeExclusive,

		QueueFamilyIndexCount: 1,
		PQueueFamilyIndices:   []uint32{uint32(queueIdx)},
	}

	var ptr vk.Buffer
	r := vk.CreateBuffer(device.Ptr(), &info, nil, &ptr)
	if r != vk.Success {
		panic("failed to create buffer")
	}

	var memreq vk.MemoryRequirements
	vk.GetBufferMemoryRequirements(device.Ptr(), ptr, &memreq)
	memreq.Deref()

	mem := device.Allocate(memreq, vk.MemoryPropertyFlags(args.Memory))

	vk.BindBufferMemory(device.Ptr(), ptr, mem.Ptr(), 0)

	return &buffer{
		ptr:    ptr,
		device: device,
		memory: mem,
		size:   int(memreq.Size),
	}
}

func NewUniform(device device.T, size int) T {
	return New(device, Args{
		Size:   size,
		Usage:  vk.BufferUsageTransferDstBit | vk.BufferUsageUniformBufferBit,
		Memory: vk.MemoryPropertyHostVisibleBit,
	})
}

func NewStorage(device device.T, size int) T {
	return New(device, Args{
		Size:   size,
		Usage:  vk.BufferUsageTransferDstBit | vk.BufferUsageStorageBufferBit,
		Memory: vk.MemoryPropertyDeviceLocalBit | vk.MemoryPropertyHostVisibleBit,
	})
}

func NewShared(device device.T, size int) T {
	return New(device, Args{
		Size:   size,
		Usage:  vk.BufferUsageTransferSrcBit | vk.BufferUsageTransferDstBit,
		Memory: vk.MemoryPropertyHostVisibleBit | vk.MemoryPropertyHostCoherentBit,
	})
}

func NewRemote(device device.T, size int, flags vk.BufferUsageFlagBits) T {
	return New(device, Args{
		Size:   size,
		Usage:  vk.BufferUsageTransferDstBit | flags,
		Memory: vk.MemoryPropertyDeviceLocalBit,
	})
}

func (b *buffer) Ptr() vk.Buffer {
	return b.ptr
}

func (b *buffer) Size() int {
	return b.size
}

func (b *buffer) Memory() device.Memory {
	return b.memory
}

func (b *buffer) Destroy() {
	vk.DestroyBuffer(b.device.Ptr(), b.ptr, nil)
	b.memory.Destroy()
	*b = *Nil.(*buffer)
}

func (b *buffer) Write(offset int, data any) int {
	return b.memory.Write(offset, data)
}

func (b *buffer) Read(offset int, data any) int {
	return b.memory.Read(offset, data)
}
