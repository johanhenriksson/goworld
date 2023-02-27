package widget

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
)

type Quad struct {
	Max     vec2.T
	Min     vec2.T
	Color   color.T
	Texture uint32
	ZIndex  float32
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
