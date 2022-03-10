package gl

import (
	"strings"

	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func GetActiveAttribute(id shader.ShaderID, index int) shader.AttributeDesc {
	var gltype uint32
	var length, size int32
	buffer := strings.Repeat("\x00", 64)
	bufferPtr := gl.Str(buffer)
	gl.GetActiveAttrib(uint32(id), uint32(index), int32(len(buffer))-1, &length, &size, &gltype, bufferPtr)
	loc := gl.GetAttribLocation(uint32(id), bufferPtr)

	return shader.AttributeDesc{
		Name: buffer[:length],
		Loc:  int(loc),
		Size: int(size),
		Type: Type(gltype).Cast(),
	}
}

func GetActiveAttributeCount(id shader.ShaderID) int {
	var attributes int32
	gl.GetProgramiv(uint32(id), gl.ACTIVE_ATTRIBUTES, &attributes)
	return int(attributes)
}

func EnablePointers(ptrs vertex.Pointers) {
	for _, p := range ptrs {
		if p.Binding < 0 {
			// destination type 0 implies that the pointer is unbound
			continue
		}
		gl.EnableVertexAttribArray(uint32(p.Binding))
		if p.Destination.Integer() {
			gl.VertexAttribIPointer(
				uint32(p.Binding),
				int32(p.Elements),
				uint32(TypeCast(p.Source)),
				int32(p.Stride),
				gl.PtrOffset(int(p.Offset)))
		} else {
			gl.VertexAttribPointer(
				uint32(p.Binding),
				int32(p.Elements),
				uint32(TypeCast(p.Source)),
				p.Normalize,
				int32(p.Stride),
				gl.PtrOffset(int(p.Offset)))
		}
	}
}

func DisablePointers(ptrs vertex.Pointers) {
	for _, p := range ptrs {
		if p.Binding < 0 {
			// destination type 0 implies that the pointer is unbound
			continue
		}
		gl.DisableVertexAttribArray(uint32(p.Binding))
	}
}
