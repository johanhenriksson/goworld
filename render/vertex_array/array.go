package vertex_array

import (
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type T interface {
	Bind() error
	Unbind()
	Delete()
	SetIndexType(t types.Type)
	Indexed() bool
	Draw() error
	Buffer(name string, data interface{})
	BufferRaw(name string, elements int, data []byte)
	BufferTo(ptrs vertex.Pointers, data interface{})
	SetPointers(vertex.Pointers)
}
