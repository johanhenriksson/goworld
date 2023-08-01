package device

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/johanhenriksson/goworld/util"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Memory interface {
	Resource[core1_0.DeviceMemory]
	Read(offset int, data any) int
	Write(offset int, data any) int
	Flush()
	Invalidate()
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
	mapPtr unsafe.Pointer
}

func alloc(device T, key string, req core1_0.MemoryRequirements, flags core1_0.MemoryPropertyFlags) Memory {
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

	if key != "" {
		device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()),
			core1_0.ObjectTypeDeviceMemory, key)
	}

	return &memory{
		device: device,
		ptr:    ptr,
		flags:  flags,
		size:   size,
	}
}

func (m *memory) isHostVisible() bool {
	bit := core1_0.MemoryPropertyHostVisible
	return m.flags&bit == bit
}

func (m *memory) isCoherent() bool {
	bit := core1_0.MemoryPropertyHostCoherent
	return m.flags&bit == bit
}

func (m *memory) Ptr() core1_0.DeviceMemory {
	return m.ptr
}

func (m *memory) Destroy() {
	m.unmap()
	m.ptr.Free(nil)
	m.ptr = nil
}

func (m *memory) mmap() {
	var nullPtr unsafe.Pointer
	if m.mapPtr != nullPtr {
		// already mapped
		return
	}
	var dst unsafe.Pointer
	dst, _, err := m.ptr.Map(0, -1, 0)
	if err != nil {
		panic(err)
	}
	m.mapPtr = dst
}

func (m *memory) unmap() {
	var nullPtr unsafe.Pointer
	if m.mapPtr == nullPtr {
		// already unmapped
		return
	}
	m.ptr.Unmap()
	m.mapPtr = nullPtr
}

func (m *memory) Write(offset int, data any) int {
	if m.ptr == nil {
		panic("write to freed memory block")
	}
	if !m.isHostVisible() {
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
	m.mmap()

	// create pointer at offset
	offsetDst := unsafe.Pointer(uintptr(m.mapPtr) + uintptr(offset))

	// copy from host
	Memcpy(offsetDst, src, size)

	// flush region
	// todo: optimize to the smallest possible region
	// m.Flush()

	// unmap shared memory
	// m.ptr.Unmap()

	return size
}

func (m *memory) Read(offset int, target any) int {
	if m.ptr == nil {
		panic("read from freed memory block")
	}
	if !m.isHostVisible() {
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
	m.mmap()

	// copy to host
	offsetPtr := unsafe.Pointer(uintptr(m.mapPtr) + uintptr(offset))
	Memcpy(dst, offsetPtr, size)

	// unmap shared memory
	// m.ptr.Unmap()

	return size
}

func (m *memory) Flush() {
	if !m.isCoherent() {
		m.ptr.FlushAll()
	}
}

func (m *memory) Invalidate() {
	m.ptr.InvalidateAll()
}
