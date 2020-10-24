package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Pointer describes a vertex pointer into a buffer
type Pointer struct {
	Name        string
	Index       int
	Source      GLType
	Destination GLType
	Elements    int
	Stride      int
	Offset      int
	Normalize   bool
}

func (p *Pointer) Enable() {
	gl.EnableVertexAttribArray(uint32(p.Index))
	if p.Destination.Integer() {
		gl.VertexAttribIPointer(
			uint32(p.Index),
			int32(p.Elements),
			uint32(p.Source),
			int32(p.Stride),
			gl.PtrOffset(int(p.Offset)))
	} else {
		gl.VertexAttribPointer(
			uint32(p.Index),
			int32(p.Elements),
			uint32(p.Source),
			p.Normalize,
			int32(p.Stride),
			gl.PtrOffset(int(p.Offset)))
	}
}

func (p *Pointer) Disable() {
	gl.DisableVertexAttribArray(uint32(p.Index))
}
