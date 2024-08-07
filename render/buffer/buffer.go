package buffer

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type T interface {
	device.Resource[core1_0.Buffer]

	// Size returns the total allocation size of the buffer in bytes
	Size() int

	// Read directly from the buffer at the given offset
	Read(offset int, data any) int

	// Write directly to the buffer at the given offset
	Write(offset int, data any) int

	Flush()

	// Memory returns a handle to the underlying memory block
	Memory() device.Memory
}

type Args struct {
	Key    string
	Size   int
	Usage  core1_0.BufferUsageFlags
	Memory core1_0.MemoryPropertyFlags
}

type buffer struct {
	ptr    core1_0.Buffer
	device device.T
	memory device.Memory
	size   int
}

func New(device device.T, args Args) T {
	if args.Size == 0 {
		panic("buffer size cant be 0")
	}

	ptr, _, err := device.Ptr().CreateBuffer(nil, core1_0.BufferCreateInfo{
		Flags:       0,
		Size:        args.Size,
		Usage:       args.Usage,
		SharingMode: core1_0.SharingModeExclusive,
		QueueFamilyIndices: []int{
			device.Queue().FamilyIndex(),
		},
	})
	if err != nil {
		panic(err)
	}

	if args.Key != "" {
		device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()),
			core1_0.ObjectTypeBuffer, args.Key)
	}

	memreq := ptr.MemoryRequirements()

	mem := device.Allocate(args.Key, *memreq, args.Memory)
	ptr.BindBufferMemory(mem.Ptr(), 0)

	return &buffer{
		ptr:    ptr,
		device: device,
		memory: mem,
		size:   int(memreq.Size),
	}
}

func NewShared(device device.T, key string, size int) T {
	return New(device, Args{
		Key:    key,
		Size:   size,
		Usage:  core1_0.BufferUsageTransferSrc | core1_0.BufferUsageTransferDst,
		Memory: core1_0.MemoryPropertyHostVisible | core1_0.MemoryPropertyHostCoherent,
	})
}

func NewRemote(device device.T, key string, size int, flags core1_0.BufferUsageFlags) T {
	return New(device, Args{
		Key:    key,
		Size:   size,
		Usage:  core1_0.BufferUsageTransferDst | flags,
		Memory: core1_0.MemoryPropertyDeviceLocal,
	})
}

func (b *buffer) Ptr() core1_0.Buffer {
	return b.ptr
}

func (b *buffer) Size() int {
	return b.size
}

func (b *buffer) Memory() device.Memory {
	return b.memory
}

func (b *buffer) Destroy() {
	b.ptr.Destroy(nil)
	b.memory.Destroy()
	b.ptr = nil
	b.memory = nil
	b.device = nil
}

func (b *buffer) Write(offset int, data any) int {
	return b.memory.Write(offset, data)
}

func (b *buffer) Read(offset int, data any) int {
	return b.memory.Read(offset, data)
}

func (b *buffer) Flush() {
	b.memory.Flush()
}
