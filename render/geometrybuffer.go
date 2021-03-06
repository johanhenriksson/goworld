package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
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
func CreateGeometryBuffer(width, height int) *GeometryBuffer {
	fbo := CreateFrameBuffer(width, height)

	g := &GeometryBuffer{
		FrameBuffer: fbo,

		Depth:    fbo.NewBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT), // depth
		Diffuse:  fbo.NewBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE),                  // diffuse (rgb)
		Normal:   fbo.NewBuffer(gl.COLOR_ATTACHMENT1, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE),                  // world normal (rgb)
		Position: fbo.NewBuffer(gl.COLOR_ATTACHMENT2, gl.RGB32F, gl.RGB, gl.FLOAT),                       // world position (rgb)
		// todo: specular & smoothness buffer maybe
	}

	// bind color buffer outputs
	// thought: is this only for the gbuffer? do we need to do this elsewhere?
	fbo.DrawBuffers()

	return g
}

// SampleNormal samples the view space normal at the given pixel location
func (g *GeometryBuffer) SampleNormal(p vec2.T) (vec3.T, bool) {
	g.Bind()
	x, y := int(p.X), int(p.Y)
	// sample normal buffer (COLOR_ATTACHMENT1)
	normalEncoded := g.FrameBuffer.Sample(gl.COLOR_ATTACHMENT1, x, int(g.Normal.Height)-y-1)
	if normalEncoded.R == 0 && normalEncoded.G == 0 && normalEncoded.B == 0 {
		return vec3.Zero, false // normal does not exist
	}

	// unpack view normal
	viewNormal := normalEncoded.Vec3().Scaled(2).Sub(vec3.One).Normalized()
	return viewNormal, true
}

// SampleDepth samples the depth at a given point
func (g *GeometryBuffer) SampleDepth(p vec2.T) (float32, bool) {
	g.Bind()
	x, y := int(p.X), int(p.Y)
	depth := g.FrameBuffer.SampleDepth(x, int(g.Depth.Height)-y-1)
	return depth, depth != 0.0
}
