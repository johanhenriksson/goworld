package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

/**
 * Geometry buffer for deferred shading
 * encapsulates an OpenGL framebuffer object
 */
type GeometryBuffer struct {
    *FrameBuffer

    /* Pointers to frame buffer textures */
    Diffuse     *Texture
    Normal      *Texture
    Depth       *Texture
}

/** Geometry buffer constructor */
func CreateGeometryBuffer(width, height int32) *GeometryBuffer {
    /* create frame buffer object */
    f := CreateFrameBuffer(width, height)

    g := &GeometryBuffer {
        FrameBuffer: f,

        // diffuse - create and attach a texture to color buffer 0 
        Diffuse: f.AddBuffer(gl.COLOR_ATTACHMENT0, gl.RGB,  gl.RGB,  gl.UNSIGNED_BYTE), // diffuse (rgb)

        // normal - create and attach a texture to color buffer 2
        Normal: f.AddBuffer(gl.COLOR_ATTACHMENT1, gl.RGB,  gl.RGB,  gl.UNSIGNED_BYTE), // world normal (rgb)

        // add a depth buffer
        Depth: f.AddBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT),
    }

    // bind color buffer outputs
    buff := []uint32 { }
    for _, buffer := range f.Buffers {
        buff = append(buff, buffer.Target)
    }
    gl.DrawBuffers(int32(len(buff)), &buff[0])
    return g
}
