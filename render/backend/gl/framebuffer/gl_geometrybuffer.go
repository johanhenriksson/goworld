package framebuffer

import (
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// glgeombuf is a frame buffer for defered shading
type glgeombuf struct {
	framebuffer.T

	diffuse  texture.T
	normal   texture.T
	position texture.T
	depth    texture.T
}

// NewGeometry creates a frame buffer suitable for storing geometry data in defered shading
func NewGeometry(width, height int) framebuffer.Geometry {
	fbo := New(width, height)

	g := &glgeombuf{
		T: fbo,

		depth:    fbo.NewBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT), // depth
		diffuse:  fbo.NewBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE),                  // diffuse (rgb)
		normal:   fbo.NewBuffer(gl.COLOR_ATTACHMENT1, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE),                  // view normal (rgb)
		position: fbo.NewBuffer(gl.COLOR_ATTACHMENT2, gl.RGB32F, gl.RGB, gl.FLOAT),                       // world position (rgb)
		// todo: specular & smoothness buffer maybe
	}

	// bind color buffer outputs
	// thought: is this only for the gbuffer? do we need to do this elsewhere?
	fbo.DrawBuffers()

	return g
}

func (g *glgeombuf) Diffuse() texture.T  { return g.diffuse }
func (g *glgeombuf) Normal() texture.T   { return g.normal }
func (g *glgeombuf) Position() texture.T { return g.position }
func (g *glgeombuf) Depth() texture.T    { return g.depth }

// SampleNormal samples the view space normal at the given pixel location
func (g *glgeombuf) SampleNormal(p vec2.T) (vec3.T, bool) {
	g.Bind()
	x, y := int(p.X), int(p.Y)
	// sample normal buffer (COLOR_ATTACHMENT1)
	normalEncoded, _ := g.Sample(gl.COLOR_ATTACHMENT1, vec2.NewI(x, g.normal.Height()-y-1))
	if normalEncoded.R == 0 && normalEncoded.G == 0 && normalEncoded.B == 0 {
		return vec3.Zero, false // normal does not exist
	}

	// unpack view normal
	viewNormal := normalEncoded.Vec3().Scaled(2).Sub(vec3.One).Normalized()
	return viewNormal, true
}

// SampleDepth samples the depth at a given point
func (g *glgeombuf) SampleDepth(p vec2.T) (float32, bool) {
	g.Bind()
	x, y := int(p.X), int(p.Y)
	depth, _ := g.T.SampleDepth(vec2.NewI(x, g.depth.Height()-y-1))
	return depth, depth != 0.0
}
