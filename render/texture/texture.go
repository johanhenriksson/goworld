package texture

import (
	"image"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/backend/types"
)

type ID uint32
type Slot uint32

type T interface {
	Bind()
	Use(slot Slot)
	FrameBufferTarget(attachment uint32)
	Resize(width, heigt int)
	Clear()

	Size() vec2.T
	Width() int
	Height() int
	Bounds() image.Rectangle

	Format() Format
	SetFormat(Format)
	InternalFormat() Format
	SetInternalFormat(Format)
	WrapMode() WrapMode
	SetWrapMode(WrapMode)
	Filter() Filter
	SetFilter(Filter)
	DataType() types.Type
	SetDataType(types.Type)

	BufferImage(*image.RGBA)
	BufferFloats([]float32)
	ToImage() *image.RGBA
	Save(path string) error
}
