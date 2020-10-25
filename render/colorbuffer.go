package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// ColorBuffer is a framebuffer with a color component
type ColorBuffer struct {
	*FrameBuffer
	Texture *Texture
}

// NewColorBuffer creates a frame buffer suitable for storing color data
func NewColorBuffer(width, height int) *ColorBuffer {
	fbo := CreateFrameBuffer(width, height)
	return &ColorBuffer{
		FrameBuffer: fbo,
		Texture:     fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE),
	}
}
