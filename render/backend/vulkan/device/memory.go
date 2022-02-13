package device

import (
	"fmt"
	"reflect"
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

type Memory interface {
	Resource[vk.DeviceMemory]
	Read(data any, offset int)
	Write(data any, offset int)
	IsHostVisible() bool
}

type memtype struct {
	TypeBits uint32
	Flags    vk.MemoryPropertyFlags
}

type memory struct {
	ptr     vk.DeviceMemory
	device  T
	size    int
	flags   vk.MemoryPropertyFlags
	hostptr unsafe.Pointer
}

func alloc(device T, req vk.MemoryRequirements, flags vk.MemoryPropertyFlags) Memory {
	typeIdx := device.GetMemoryTypeIndex(req.MemoryTypeBits, flags)

	alloc := vk.MemoryAllocateInfo{
		SType:           vk.StructureTypeMemoryAllocateInfo,
		AllocationSize:  req.Size,
		MemoryTypeIndex: uint32(typeIdx),
	}
	var ptr vk.DeviceMemory
	vk.AllocateMemory(device.Ptr(), &alloc, nil, &ptr)

	m := &memory{
		device: device,
		ptr:    ptr,
		flags:  flags,
		size:   int(req.Size),
	}

	return m
}

func (m *memory) IsHostVisible() bool {
	bit := vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit)
	return m.flags&bit == bit
}

func (m *memory) Ptr() vk.DeviceMemory {
	return m.ptr
}

func (m *memory) Destroy() {
	vk.FreeMemory(m.device.Ptr(), m.ptr, nil)
	m.ptr = nil
}

func (m *memory) Write(data any, offset int) {
	if !m.IsHostVisible() {
		panic("memory is not visible to host")
	}

	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Slice {
		panic(fmt.Errorf("buffered data must be a slice"))
	}

	v := reflect.ValueOf(data)

	// calculate copy size
	count := v.Len()
	sizeof := int(t.Elem().Size())
	size := count * sizeof

	// get a pointer to the beginning of the array
	src := unsafe.Pointer(v.Pointer())

	if size+offset > m.size {
		panic("out of bounds")
	}

	// map shared memory
	var dst unsafe.Pointer
	vk.MapMemory(
		m.device.Ptr(),
		m.ptr,
		vk.DeviceSize(offset),
		vk.DeviceSize(size),
		vk.MemoryMapFlags(0),
		&dst)
	// copy from host
	memcpy(dst, src, size)

	// unmap shared memory
	vk.UnmapMemory(m.device.Ptr(), m.Ptr())
}

func (m *memory) Read(target any, offset int) {
	if !m.IsHostVisible() {
		panic("memory is not visible to host")
	}

	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Slice {
		panic(fmt.Errorf("output buffer must be a slice"))
	}

	v := reflect.ValueOf(target)

	// calculate copy size
	count := v.Len()
	sizeof := int(t.Elem().Size())
	size := count * sizeof

	// get a pointer to the beginning of the array
	dst := unsafe.Pointer(v.Pointer())

	if size+offset > m.size {
		panic("out of bounds")
	}

	// map shared memory
	var src unsafe.Pointer
	vk.MapMemory(
		m.device.Ptr(),
		m.ptr,
		vk.DeviceSize(offset),
		vk.DeviceSize(size),
		vk.MemoryMapFlags(0),
		&src)

	// copy to host
	memcpy(dst, src, size)

	// unmap shared memory
	vk.UnmapMemory(m.device.Ptr(), m.Ptr())
}
