package device

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Memory struct {
	ptr    core1_0.DeviceMemory
	device *Device
	size   int
	flags  core1_0.MemoryPropertyFlags
	mapPtr unsafe.Pointer
}

func (m *Memory) isHostVisible() bool {
	bit := core1_0.MemoryPropertyHostVisible
	return m.flags&bit == bit
}

func (m *Memory) isCoherent() bool {
	bit := core1_0.MemoryPropertyHostCoherent
	return m.flags&bit == bit
}

func (m *Memory) Ptr() core1_0.DeviceMemory {
	return m.ptr
}

func (m *Memory) Destroy() {
	m.unmap()
	m.ptr.Free(nil)
	m.ptr = nil
}

func (m *Memory) mmap() {
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

func (m *Memory) unmap() {
	var nullPtr unsafe.Pointer
	if m.mapPtr == nullPtr {
		// already unmapped
		return
	}
	m.ptr.Unmap()
	m.mapPtr = nullPtr
}

func (m *Memory) Write(offset int, data any) int {
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

func (m *Memory) Read(offset int, target any) int {
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

func (m *Memory) Flush() {
	if !m.isCoherent() {
		m.mmap()
		m.ptr.FlushAll()
	}
}

func (m *Memory) Invalidate() {
	m.ptr.InvalidateAll()
}
