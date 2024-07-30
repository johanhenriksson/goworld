package widget

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
)

type Quad struct {
	Min      vec2.T
	Max      vec2.T
	MinUV    vec2.T
	MaxUV    vec2.T
	Color    [4]color.T
	ZIndex   float32
	Radius   float32
	Softness float32
	Border   float32
	Texture  uint32
	_padding [3]uint32
}

func init() {
	size := unsafe.Sizeof(Quad{})
	if size != 128 {
		panic(fmt.Sprintf("Quad size is not 64 bytes, was %d", size))
	}
	align := unsafe.Alignof(Quad{})
	if align != 4 {
		panic(fmt.Sprintf("Quad alignment is not 4 bytes, was %d", align))
	}
}

type QuadBuffer struct {
	Data []Quad
}

func NewQuadBuffer(capacity int) *QuadBuffer {
	return &QuadBuffer{
		Data: make([]Quad, 0, capacity),
	}
}

func (qb *QuadBuffer) Push(quad Quad) {
	qb.Data = append(qb.Data, quad)
}

func (qb *QuadBuffer) Reset() {
	qb.Data = qb.Data[:0]
}
