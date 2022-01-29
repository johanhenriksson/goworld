package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// gldepthbuf is a frame buffer containing only a depth texture
type gldepthbuf struct {
	framebuffer.T

	depth texture.T
}

// NewShadowBuffer creates a new shadow buffer
func NewDepth(width, height int) framebuffer.Depth {
	f := New(width, height)
	return &gldepthbuf{
		T:     f,
		depth: f.NewBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT),
	}
}

func (d *gldepthbuf) Depth() texture.T {
	return d.depth
}
