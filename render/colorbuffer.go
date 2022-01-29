package render

import (
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// ColorBuffer is a framebuffer with a color component
type ColorBuffer struct {
	*FrameBuffer
	Texture texture.T
}

// NewColorBuffer creates a frame buffer suitable for storing color data
func NewColorBuffer(width, height int) *ColorBuffer {
	fbo := CreateFrameBuffer(width, height)
	return &ColorBuffer{
		FrameBuffer: fbo,
		Texture:     fbo.NewBuffer(gl.COLOR_ATTACHMENT0, texture.RGB, texture.RGB, types.UInt8),
	}
}
