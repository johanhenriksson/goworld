package device

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/johanhenriksson/goworld/util"
	vk "github.com/vulkan-go/vulkan"
)

type Memory interface {
	Resource[vk.DeviceMemory]
	Read(offset int, data any)
	Write(offset int, data any)
	Flush()
	Invalidate()
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

	align := int(device.GetLimits().NonCoherentAtomSize)
	size := util.Align(int(req.Size), align)

	alloc := vk.MemoryAllocateInfo{
		SType:           vk.StructureTypeMemoryAllocateInfo,
		AllocationSize:  vk.DeviceSize(size),
		MemoryTypeIndex: uint32(typeIdx),
	}
	var ptr vk.DeviceMemory
	r := vk.AllocateMemory(device.Ptr(), &alloc, nil, &ptr)
	if r != vk.Success {
		panic(fmt.Sprintf("failed to allocate %d bytes of memory", req.Size))
	}

	m := &memory{
		device: device,
		ptr:    ptr,
		flags:  flags,
		size:   size,
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

type eface struct {
	rtype, ptr unsafe.Pointer
}

func (m *memory) Write(offset int, data any) {
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

	if size+offset > m.size {
		panic("out of bounds")
	}

	// map shared memory
	var dst unsafe.Pointer
	vk.MapMemory(
		m.device.Ptr(),
		m.ptr,
		vk.DeviceSize(0),
		vk.DeviceSize(vk.WholeSize),
		vk.MemoryMapFlags(0),
		&dst)

	// create pointer at offset
	offsetDst := unsafe.Pointer(uintptr(dst) + uintptr(offset))

	// copy from host
	memcpy(offsetDst, src, size)

	// flush region
	// todo: optimize to the smallest possible region
	m.Flush()

	// unmap shared memory
	vk.UnmapMemory(m.device.Ptr(), m.Ptr())
}

func (m *memory) Read(offset int, target any) {
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

func (m *memory) Flush() {
	vk.FlushMappedMemoryRanges(m.device.Ptr(), 1, []vk.MappedMemoryRange{
		{
			SType:  vk.StructureTypeMappedMemoryRange,
			Memory: m.ptr,
			Offset: 0,
			Size:   vk.DeviceSize(vk.WholeSize),
		},
	})
}

func (m *memory) Invalidate() {
	vk.InvalidateMappedMemoryRanges(m.device.Ptr(), 1, []vk.MappedMemoryRange{
		{
			SType:  vk.StructureTypeMappedMemoryRange,
			Memory: m.ptr,
			Offset: 0,
			Size:   vk.DeviceSize(vk.WholeSize),
		},
	})
}
