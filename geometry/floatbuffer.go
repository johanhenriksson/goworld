package geometry

import (
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type FloatBuffer []float32

func (floats FloatBuffer) Elements() int {
    return len(floats)
}

func (floats FloatBuffer) Size() int {
    return 4
}

func (floats FloatBuffer) GLPtr() unsafe.Pointer {
    return gl.Ptr(floats)
}
