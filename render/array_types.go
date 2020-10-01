package render

import (
	"unsafe"
)

type UInt32Buffer []uint32

func (a UInt32Buffer) Elements() int {
	return len(a)
}

func (a UInt32Buffer) Size() int {
	return 4
}

func (a UInt32Buffer) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&a[0])
}

type Int32Buffer []int32

func (a Int32Buffer) Elements() int {
	return len(a)
}

func (a Int32Buffer) Size() int {
	return 4
}

func (a Int32Buffer) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&a[0])
}

// FloatBuffer is a simple implementation of the VertexData interface for buffering arrays of 32-bit floats
type FloatBuffer []float32

// Elements returns the number of vertex elements in the buffer
func (vtx FloatBuffer) Elements() int {
	return len(vtx)
}

// Size returns the byte size of a buffer element
func (vtx FloatBuffer) Size() int {
	return 4
}

func (vtx FloatBuffer) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&vtx[0])
}
