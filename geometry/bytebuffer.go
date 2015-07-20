package geometry

import (
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type ByteBuffer []uint8

func (bytes ByteBuffer) Elements() int {
    return len(bytes)
}

func (bytes ByteBuffer) Size() int {
    return 1
}

func (bytes ByteBuffer) GLPtr() unsafe.Pointer {
    return gl.Ptr(bytes)
}

func (bytes ByteBuffer) GLType() uint32 {
    return gl.UNSIGNED_BYTE
}
