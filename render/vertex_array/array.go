package vertex_array

import (
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/shader"
)

type T interface {
	Bind()
	Unbind()
	Delete()
	SetIndexType(t types.Type)
	Indexed() bool
	Draw()
	Buffer(name string, data interface{})
	BufferTo(ptrs shader.Pointers, data interface{})
}
