package effect

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/screen_quad"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_shader"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// GaussianPass represents a gaussian blur pass.
type GaussianPass struct {
	Output texture.T

	input  texture.T
	fbo    framebuffer.T
	shader shader.T
	mat    material.T
	quad   screen_quad.T
}

// NewGaussianPass creates a new Gaussian Blur pass.
func NewGaussianPass(input texture.T) *GaussianPass {
	fbo := gl_framebuffer.New(input.Width(), input.Height())
	output := fbo.NewBuffer(gl.COLOR_ATTACHMENT0, texture.Red, texture.RGB, types.Float)

	shader := gl_shader.CompileShader(
		"gaussian_pass",
		"assets/shaders/pass/postprocess.vs",
		"assets/shaders/pass/gaussian.fs")

	mat := material.New("gaussian_pass", shader)
	mat.Texture("tex_input", input)

	return &GaussianPass{
		Output: output,
		input:  input,
		quad:   screen_quad.New(shader),
		fbo:    fbo,
		shader: shader,
		mat:    mat,
	}
}

// Draw draws the gaussian blurred output to the output frame buffer.
func (p *GaussianPass) Draw(args render.Args, scene object.T) {
	render.Blend(false)
	render.DepthOutput(false)

	p.fbo.Bind()
	defer p.fbo.Unbind()
	// make sure output matches input size
	p.fbo.Resize(p.input.Width(), p.input.Height())

	p.mat.Use()
	p.quad.Draw()
}
