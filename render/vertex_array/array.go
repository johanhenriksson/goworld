package vertex_array

import (
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	Bind() error
	Unbind()
	Delete()
	Indexed() bool
	Draw() error
	Buffer(name string, data interface{})
	BufferRaw(name string, elements int, data []byte)
	SetPointers(vertex.Pointers)
}
