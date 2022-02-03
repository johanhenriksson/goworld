package vertex_array

import (
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/shader"
)

type T interface {
	Bind() error
	Unbind()
	Delete()
	SetIndexType(t types.Type)
	Indexed() bool
	Draw() error
	Buffer(name string, data interface{})
	BufferTo(ptrs shader.Pointers, data interface{})
}
