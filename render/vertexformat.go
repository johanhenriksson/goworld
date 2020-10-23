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
		Normalize: false,
	},
	BufferDescriptor{
		Name:      "texcoord",
		Type:      gl.FLOAT,
		Elements:  2,
		Stride:    20,
		Offset:    12,
		Normalize: false,
	},
}

// F32_XYZNUV contains position, uvs and normals as 32 bit floats.
var F32_XYZNUV = VertexFormat{
	BufferDescriptor{
		Buffer:    "geometry",
		Name:      "position",
		Type:      gl.FLOAT,
		Elements:  3,
		Stride:    32,
		Offset:    0,
		Normalize: false,
	},
	BufferDescriptor{
		Buffer:    "geometry",
		Name:      "normal",
		Type:      gl.FLOAT,
		Elements:  3,
		Stride:    32,
		Offset:    12,
		Normalize: false,
	},
	BufferDescriptor{
		Buffer:    "geometry",
		Name:      "texcoord",
		Type:      gl.FLOAT,
		Elements:  2,
		Stride:    32,
		Offset:    24,
		Normalize: false,
	},
}
