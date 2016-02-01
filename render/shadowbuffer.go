package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShadowBuffer struct {
    *FrameBuffer
    Depth   *Texture
}

/** Shadow buffer constructor */
func NewShadowBuffer(width, height int32) *ShadowBuffer {
    /* create frame buffer object */
    f := CreateFrameBuffer(width, height)
    s := &ShadowBuffer {
        FrameBuffer: f,

        // add a depth buffer
        Depth: f.AddBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT),
    }
    return s
}
