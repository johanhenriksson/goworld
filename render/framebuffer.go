package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type DrawBuffer struct {
    Target uint32 // GL attachment
    Texture *Texture
}

/** Represents an OpenGL frame buffer object */
type FrameBuffer struct {
    Buffers []DrawBuffer
    Width   int32
    Height  int32
    id      uint32
    mipLvl  int32
}

/* TODO: Rename to CreateBuffer */
/** 
 * Create a new frame buffer texture and attach it to the given target.
 * Returns a pointer to the created texture object 
 */
func (f *FrameBuffer) AddBuffer(target, internal_fmt, format, datatype uint32) *Texture {
    // Create texture object 
    texture := CreateTexture(f.Width, f.Height)
    texture.Format = format
    texture.InternalFormat = internal_fmt
    texture.DataType = datatype
    texture.Clear()

    // Set texture as frame buffer target
    texture.FrameBufferTarget(target)

    if target != gl.DEPTH_ATTACHMENT {
        // Attach to frame buffer
        f.Buffers = append(f.Buffers, DrawBuffer {
            Target: target,
            Texture: texture,
        })
    }

    return texture
}

func CreateFrameBuffer(width, height int32) *FrameBuffer {
    f := &FrameBuffer {
        Width: width,
        Height: height,
        Buffers: []DrawBuffer { },
    }
    gl.GenFramebuffers(1, &f.id)
    gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
    return f
}

func (f *FrameBuffer) Bind() {
    gl.BindTexture(gl.TEXTURE_2D, 0) // why?

    // bind this frame buffer
    gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)

    // set viewport size equal to buffer size
    gl.Viewport(0, 0, f.Width, f.Height)
}

func (f *FrameBuffer) Unbind() {
    // finish drawing
    gl.Flush()

    // unbind
    gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

/* Clear the frame buffer. Make sure its bound first */
func (f *FrameBuffer) Clear() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
}

/* Delete frame buffer object */
func (f *FrameBuffer) Delete() {
    gl.DeleteFramebuffers(1, &f.id)
    f.id = 0
}
