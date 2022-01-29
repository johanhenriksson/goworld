package shader

import (
	"github.com/johanhenriksson/goworld/render/backend/gl"

	ogl "github.com/go-gl/gl/v4.1-core/gl"
)

// Pointer describes a vertex pointer into a buffer
type Pointer struct {
	Name        string
	Index       int
	Source      gl.Type
	Destination gl.Type
	Elements    int
	Stride      int
	Offset      int
	Normalize   bool
}

func (p Pointer) String() string {
	return p.Name
}

func (p Pointer) Enable() {
	ogl.EnableVertexAttribArray(uint32(p.Index))
	if p.Destination.Integer() {
		ogl.VertexAttribIPointer(
			uint32(p.Index),
			int32(p.Elements),
			uint32(p.Source),
			int32(p.Stride),
			ogl.PtrOffset(int(p.Offset)))
	} else {
		ogl.VertexAttribPointer(
			uint32(p.Index),
			int32(p.Elements),
			uint32(p.Source),
			p.Normalize,
			int32(p.Stride),
			ogl.PtrOffset(int(p.Offset)))
	}
}

func (p Pointer) Disable() {
	ogl.DisableVertexAttribArray(uint32(p.Index))
}
