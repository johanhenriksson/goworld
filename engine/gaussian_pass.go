package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

// GaussianPass represents a gaussian blur pass.
type GaussianPass struct {
	Output texture.T

	fbo      *render.FrameBuffer
	shader   *render.Shader
	textures *render.TextureMap
	quad     *Quad
}

// NewGaussianPass creates a new Gaussian Blur pass.
func NewGaussianPass(input texture.T) *GaussianPass {
	fbo := render.CreateFrameBuffer(input.Width(), input.Height())
	output := fbo.NewBuffer(gl.COLOR_ATTACHMENT0, gl.RED, texture.RGB, gl.FLOAT)

	shader := render.CompileShader(
		"gaussian_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/gaussian.fs")

	tx := render.NewTextureMap(shader)
	tx.Add("tex_input", input)

	return &GaussianPass{
		Output:   output,
		quad:     NewQuad(shader),
		fbo:      fbo,
		shader:   shader,
		textures: tx,
	}
}

// DrawPass draws the gaussian blurred output to the frame buffer.
func (p *GaussianPass) DrawPass(args render.Args, scene scene.T) {
	render.Blend(false)
	render.DepthOutput(false)

	p.fbo.Bind()
	defer p.fbo.Unbind()
	p.fbo.Resize(args.Viewport.FrameWidth, args.Viewport.FrameHeight)

	p.shader.Use()
	p.textures.Use()

	render.ClearWith(color.White)
	p.quad.Draw()
}
