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

type GeometryBuffer struct {
    *FrameBuffer
    Diffuse     *Texture
    Specular    *Texture
    Normal      *Texture
    Depth       *Texture
}

func (f *FrameBuffer) AddBuffer(target, internal_fmt, format, datatype uint32) *Texture {
    if target == gl.DEPTH_ATTACHMENT {
        /* Set up a depth render buffer? */
    }
    texture := CreateTexture(f.Width, f.Height)
    texture.Format = format
    texture.InternalFormat = internal_fmt
    texture.DataType = datatype
    texture.Clear()
    texture.FrameBufferTarget(target)
    f.Buffers = append(f.Buffers, DrawBuffer {
        Target: target,
        Texture: texture,
    })
    return texture
}

/** Sets up a geometry buffer for defered shading */
func CreateGeometryBuffer(width, height int32) *GeometryBuffer {
    f := CreateFrameBuffer(width, height)
    g := &GeometryBuffer {
        FrameBuffer: f,
        Diffuse: f.AddBuffer(gl.COLOR_ATTACHMENT0, gl.RGB,  gl.RGB,  gl.UNSIGNED_BYTE), // diffuse (rgb)
        Specular: f.AddBuffer(gl.COLOR_ATTACHMENT1, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE), // specular (rgb) + smoothness (a)
        Normal: f.AddBuffer(gl.COLOR_ATTACHMENT2, gl.RGB,  gl.RGB,  gl.UNSIGNED_BYTE), // world normal (rgb)
        Depth: f.AddBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT),
    }
    buff := []uint32 { }
    for _, buffer := range f.Buffers {
        buff = append(buff, buffer.Target)
    }
    gl.DrawBuffers(int32(len(buff)), &buff[0])
    return g
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

func CreateRenderTexture() *FrameBuffer {
    f := &FrameBuffer { }
    /* TODO */
    return f
}

func (f *FrameBuffer) Bind() {
    gl.BindTexture(gl.TEXTURE_2D, 0)
    gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
    gl.Viewport(0, 0, f.Width, f.Height)
}

func (f *FrameBuffer) Unbind() {
    gl.Flush()
    gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (f *FrameBuffer) Delete() {
    gl.DeleteFramebuffers(1, &f.id)
    f.id = 0
}
