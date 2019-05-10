package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// VertexFormat describes vertex buffer pointers
type VertexFormat []BufferDescriptor

// F32_XYZUV contains position & texture coords as 5 32-bit floats: X, Y, Z, U, V
var F32_XYZUV = VertexFormat{
	BufferDescriptor{
		Name:      "position",
		Type:      gl.FLOAT,
		Elements:  3,
		Stride:    20,
		Offset:    0,
		Integer:   false,
		Normalize: false,
	},
	BufferDescriptor{
		Name:      "texcoord",
		Type:      gl.FLOAT,
		Elements:  2,
		Stride:    20,
		Offset:    12,
		Integer:   false,
		Normalize: false,
	},
}
