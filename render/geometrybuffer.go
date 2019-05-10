package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

// GeometryBuffer is a frame buffer for defered shading
type GeometryBuffer struct {
	*FrameBuffer
	Diffuse  *Texture
	Normal   *Texture
	Position *Texture
	Depth    *Texture
}

// CreateGeometryBuffer creates a frame buffer suitable for storing geometry data in defered shading
func CreateGeometryBuffer(width, height int32) *GeometryBuffer {
	fbo := CreateFrameBuffer(width, height)

	g := &GeometryBuffer{
		FrameBuffer: fbo,

		Depth:    fbo.AddBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT), // depth
		Diffuse:  fbo.AddBuffer(gl.COLOR_ATTACHMENT0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE),                // diffuse (rgb)
		Normal:   fbo.AddBuffer(gl.COLOR_ATTACHMENT1, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE),                // world normal (rgb)
		Position: fbo.AddBuffer(gl.COLOR_ATTACHMENT2, gl.RGB32F, gl.RGBA, gl.FLOAT),                      // world position (rgb)
		// todo: specular & smoothness buffer maybe
	}

	// bind color buffer outputs
	// thought: is this only for the gbuffer? do we need to do this elsewhere?
	fbo.DrawBuffers()

	return g
}

// SampleNormal samples the view space normal at the given pixel location
func (g *GeometryBuffer) SampleNormal(x, y int) (mgl.Vec3, bool) {
	g.Bind()

	// sample normal buffer (COLOR_ATTACHMENT1)
	normalEncoded := g.FrameBuffer.Sample(gl.COLOR_ATTACHMENT1, x, y)
	if normalEncoded.R == 0 && normalEncoded.G == 0 && normalEncoded.B == 0 {
		return mgl.Vec3{}, false // normal does not exist
	}

	// unpack view normal
	viewNormal := normalEncoded.Vec3().Mul(2).Sub(mgl.Vec3{1, 1, 1}).Normalize()
	return viewNormal, true
}
