package vertex_buffer

import "unsafe"

type T interface {
	Bind()
	Unbind()
	Delete()
	Buffer(data interface{}) (int, int)
	BufferFrom(ptr unsafe.Pointer, size int)
}
