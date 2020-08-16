package render

import "unsafe"

// VertexData is an interface for data types that can be stored in a vertex buffer object
type VertexData interface {
	Elements() int /* Number of items, usually len(slice) */
	Size() int     /* Size of each individual element */
	Pointer() unsafe.Pointer
}
