package device

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/johanhenriksson/goworld/util"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Memory interface {
	Resource[core1_0.DeviceMemory]
	Read(offset int, data any) int
	Write(offset int, data any) int
	Flush()
	Invalidate()
	IsHostVisible() bool
}

type memtype struct {
	TypeBits uint32
	Flags    core1_0.MemoryPropertyFlags
}

type memory struct {
	ptr    core1_0.DeviceMemory
	device T
	size   int
	flags  core1_0.MemoryPropertyFlags
}

func alloc(device T, req core1_0.MemoryRequirements, flags core1_0.MemoryPropertyFlags) Memory {
	typeIdx := device.GetMemoryTypeIndex(req.MemoryTypeBits, flags)

	align := int(device.GetLimits().NonCoherentAtomSize)
	size := util.Align(int(req.Size), align)

	ptr, _, err := device.Ptr().AllocateMemory(nil, core1_0.MemoryAllocateInfo{
		AllocationSize:  size,
		MemoryTypeIndex: typeIdx,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to allocate %d bytes of memory: %s", req.Size, err))
	}

	return &memory{
		device: device,
		ptr:    ptr,
		flags:  flags,
		size:   size,
	}
}

func (m *memory) IsHostVisible() bool {
	bit := core1_0.MemoryPropertyHostVisible
	return m.flags&bit == bit
}

func (m *memory) Ptr() core1_0.DeviceMemory {
	return m.ptr
}

func (m *memory) Destroy() {
	m.ptr.Free(nil)
	m.ptr = nil
}

func (m *memory) Write(offset int, data any) int {
	if m.ptr == nil {
		panic("write to freed memory block")
	}
	if !m.IsHostVisible() {
		panic("memory is not visible to host")
	}

	size := 0
	var src unsafe.Pointer

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Slice {
		// calculate copy size
		count := v.Len()
		sizeof := int(t.Elem().Size())
		size = count * sizeof

		// get a pointer to the beginning of the array
		src = unsafe.Pointer(v.Pointer())
	} else if t.Kind() == reflect.Pointer {
		src = v.UnsafePointer()
		size = int(v.Elem().Type().Size())
	} else {
		panic(fmt.Errorf("buffered data must be a slice, struct or a pointer"))
	}

	if offset < 0 || offset+size > m.size {
		panic("out of bounds")
	}

	// map shared memory
	var dst unsafe.Pointer
	dst, _, err := m.ptr.Map(0, -1, 0)
	if err != nil {
		panic(err)
	}

	// create pointer at offset
	offsetDst := unsafe.Pointer(uintptr(dst) + uintptr(offset))

	// copy from host
	Memcpy(offsetDst, src, size)

	// flush region
	// todo: optimize to the smallest possible region
	m.Flush()

	// unmap shared memory
	m.ptr.Unmap()

	return size
}

func (m *memory) Read(offset int, target any) int {
	if m.ptr == nil {
		panic("read from freed memory block")
	}
	if !m.IsHostVisible() {
		panic("memory is not visible to host")
	}

	size := 0
	var dst unsafe.Pointer

	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	if t.Kind() == reflect.Slice {
		// calculate copy size
		count := v.Len()
		sizeof := int(t.Elem().Size())
		size = count * sizeof

		// get a pointer to the beginning of the array
		dst = unsafe.Pointer(v.Pointer())
	} else if t.Kind() == reflect.Pointer {
		dst = v.UnsafePointer()
		size = int(v.Elem().Type().Size())
	} else {
		panic(fmt.Errorf("buffered data must be a slice, struct or a pointer"))
	}

	if size+offset > m.size {
		panic("out of bounds")
	}

	// map shared memory
	var src unsafe.Pointer
	src, _, err := m.ptr.Map(offset, -1, 0)
	if err != nil {
		panic(err)
	}

	// copy to host
	Memcpy(dst, src, size)

	// unmap shared memory
	m.ptr.Unmap()

	return size
}

func (m *memory) Flush() {
	m.ptr.FlushAll()
}

func (m *memory) Invalidate() {
	m.ptr.InvalidateAll()
}
