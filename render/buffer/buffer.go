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
	Memory() *device.Memory
}

type Args struct {
	Key    string
	Size   int
	Usage  core1_0.BufferUsageFlags
	Memory device.MemoryType
}

type Buffer struct {
	ptr    core1_0.Buffer
	device *device.Device
	memory *device.Memory
	size   int
}

func New(device *device.Device, args Args) *Buffer {
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

	mem := device.Allocate(args.Key, args.Size, args.Memory)
	ptr.BindBufferMemory(mem.Ptr(), 0)

	return &Buffer{
		ptr:    ptr,
		device: device,
		memory: mem,
		size:   int(args.Size),
	}
}

func NewCpuLocal(dev *device.Device, key string, size int) *Buffer {
	return New(dev, Args{
		Key:    key,
		Size:   size,
		Usage:  core1_0.BufferUsageTransferSrc,
		Memory: device.MemoryTypeCPU,
	})
}

func NewShared(dev *device.Device, key string, size int) *Buffer {
	return New(dev, Args{
		Key:    key,
		Size:   size,
		Usage:  core1_0.BufferUsageTransferSrc | core1_0.BufferUsageTransferDst,
		Memory: device.MemoryTypeShared,
	})
}

func NewGpuLocal(dev *device.Device, key string, size int, flags core1_0.BufferUsageFlags) *Buffer {
	return New(dev, Args{
		Key:    key,
		Size:   size,
		Usage:  core1_0.BufferUsageTransferDst | flags,
		Memory: device.MemoryTypeGPU,
	})
}

func (b *Buffer) Ptr() core1_0.Buffer {
	return b.ptr
}

func (b *Buffer) Size() int {
	return b.size
}

func (b *Buffer) Memory() *device.Memory {
	return b.memory
}

func (b *Buffer) Destroy() {
	b.ptr.Destroy(nil)
	b.memory.Destroy()
	b.ptr = nil
	b.memory = nil
	b.device = nil
}

func (b *Buffer) Write(offset int, data any) int {
	return b.memory.Write(offset, data)
}

func (b *Buffer) Read(offset int, data any) int {
	return b.memory.Read(offset, data)
}

func (b *Buffer) Flush() {
	b.memory.Flush()
}
