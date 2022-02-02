package framebuffer

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Target uint32

type T interface {
	Width() int
	Height() int
	Size() vec2.T
	Bind()
	Unbind()
	Delete()
	Sample(target Target, pos vec2.T) (color.T, bool)
	SampleDepth(pos vec2.T) (float32, bool)
	NewBuffer(target Target, internalFormat, format texture.Format, datatype types.Type) texture.T
	AttachBuffer(target Target, tex texture.T)
	DrawBuffers()
	Resize(int, int)
}

// Buffer holds a target texture of a frame buffer object
type Buffer struct {
	Target  Target
	Texture texture.T
}
