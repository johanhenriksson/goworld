package vertex_buffer

import "unsafe"

type T interface {
	Bind()
	Unbind()
	Delete()
	Buffer(data interface{}) int
	BufferFrom(elements, size int, ptr unsafe.Pointer)
}
