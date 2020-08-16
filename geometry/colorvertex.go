package geometry

import (
	"github.com/johanhenriksson/goworld/render"
	"unsafe"
)

// ColorVertex is used to represent vertices in solid-color elements
type ColorVertex struct {
	X, Y, Z      float32 // 12 bytes
	render.Color         // 16 bytes
} // 28 bytes

// ColorVertices is a GPU bufferable slice of ColorVertex objects
type ColorVertices []ColorVertex

// Elements returns the number of verticies
func (buffer ColorVertices) Elements() int {
	return len(buffer)
}

// Size returns the element size in bytes
func (buffer ColorVertices) Size() int {
	return 28
}

func (buffer ColorVertices) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&buffer[0])
}
