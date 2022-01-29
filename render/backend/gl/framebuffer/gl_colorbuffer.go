package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// ColorBuffer is a framebuffer with a color component
type glcolorbuf struct {
	framebuffer.T
	texture texture.T
}

func (c *glcolorbuf) Texture() texture.T { return c.texture }

// NewColorBuffer creates a frame buffer suitable for storing color data
func NewColor(width, height int) framebuffer.Color {
	fbo := New(width, height)
	return &glcolorbuf{
		T:       fbo,
		texture: fbo.NewBuffer(gl.COLOR_ATTACHMENT0, texture.RGB, texture.RGB, gl.UNSIGNED_BYTE),
	}
}
