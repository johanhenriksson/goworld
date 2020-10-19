package geometry

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// ColorVertex is used to represent vertices in solid-color elements
type ColorVertex struct {
	Position vec3.T       // 12 bytes
	Normal   vec3.T       // 12 bytes
	Color    render.Color // 16 bytes
} // 40 bytes

// ColorVertices is a GPU bufferable slice of ColorVertex objects
type ColorVertices []ColorVertex

// Elements returns the number of verticies
func (buffer ColorVertices) Elements() int {
	return len(buffer)
}

// Size returns the element size in bytes
func (buffer ColorVertices) Size() int {
	return 40
}

// Pointer returns an unsafe pointer to the first element
func (buffer ColorVertices) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&buffer[0])
}
