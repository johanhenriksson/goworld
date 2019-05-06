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
	Diffuse *Texture
	Normal  *Texture
	Depth   *Texture
}

/** Geometry buffer constructor */
func CreateGeometryBuffer(width, height int32) *GeometryBuffer {
	/* create frame buffer object */
	f := CreateFrameBuffer(width, height)

	g := &GeometryBuffer{
		FrameBuffer: f,

		Diffuse: f.AddBuffer(gl.COLOR_ATTACHMENT0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE), // diffuse (rgb)
		Normal:  f.AddBuffer(gl.COLOR_ATTACHMENT1, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE), // world normal (rgb)
		// todo: specular & smoothness buffer

		Depth: f.AddBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT), // depth
	}

	// bind color buffer outputs
	buff := []uint32{}
	for _, buffer := range f.Buffers {
		buff = append(buff, buffer.Target)
	}
	gl.DrawBuffers(int32(len(buff)), &buff[0])
	return g
}
