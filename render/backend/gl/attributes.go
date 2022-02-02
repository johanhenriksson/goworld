package gl

import (
	"strings"

	"github.com/johanhenriksson/goworld/render/shader"

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
		Name:  buffer[:length],
		Index: int(loc),
		Size:  int(size),
		Type:  Type(gltype).Cast(),
	}
}

func GetActiveAttributeCount(id shader.ShaderID) int {
	var attributes int32
	gl.GetProgramiv(uint32(id), gl.ACTIVE_ATTRIBUTES, &attributes)
	return int(attributes)
}
