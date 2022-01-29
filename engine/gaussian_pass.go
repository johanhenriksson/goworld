package engine

import (
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
	glframebuf "github.com/johanhenriksson/goworld/render/backend/gl/framebuffer"
	glshader "github.com/johanhenriksson/goworld/render/backend/gl/shader"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// GaussianPass represents a gaussian blur pass.
type GaussianPass struct {
	Output texture.T

	fbo    framebuffer.T
	shader shader.T
	mat    material.T
	quad   *Quad
}

// NewGaussianPass creates a new Gaussian Blur pass.
func NewGaussianPass(input texture.T) *GaussianPass {
	fbo := glframebuf.New(input.Width(), input.Height())
	output := fbo.NewBuffer(gl.COLOR_ATTACHMENT0, texture.Red, texture.RGB, types.Float)

	shader := glshader.CompileShader(
		"gaussian_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/gaussian.fs")

	mat := material.New("gaussian_pass", shader)
	mat.Texture("tex_input", input)

	return &GaussianPass{
		Output: output,
		quad:   NewQuad(shader),
		fbo:    fbo,
		shader: shader,
		mat:    mat,
	}
}

// DrawPass draws the gaussian blurred output to the frame buffer.
func (p *GaussianPass) DrawPass(args render.Args, scene scene.T) {
	render.Blend(false)
	render.DepthOutput(false)

	p.fbo.Bind()
	defer p.fbo.Unbind()
	p.fbo.Resize(args.Viewport.FrameWidth, args.Viewport.FrameHeight)

	p.mat.Use()

	render.ClearWith(color.White)
	p.quad.Draw()
}
