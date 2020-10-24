package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render"
)

// GaussianPass represents a gaussian blur pass.
type GaussianPass struct {
	fbo      *render.FrameBuffer
	material *render.Material
	quad     *Quad
	Output   *render.Texture
}

// NewGaussianPass creates a new Gaussian Blur pass.
func NewGaussianPass(input *render.Texture) *GaussianPass {
	fbo := render.CreateFrameBuffer(input.Width, input.Height)
	texture := fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RED, gl.RGB, gl.FLOAT)

	mat := render.CreateMaterial("gaussian_pass", render.CompileShader(
		"gaussian_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/gaussian.fs"))
	mat.AddTexture("tex_input", input)

	quad := NewQuad(mat)

	return &GaussianPass{
		fbo:      fbo,
		material: mat,
		quad:     quad,
		Output:   texture,
	}
}

// DrawPass draws the gaussian blurred output to the frame buffer.
func (p *GaussianPass) DrawPass(scene *Scene) {
	render.Blend(false)
	render.DepthOutput(false)

	p.fbo.Bind()
	defer p.fbo.Unbind()

	render.ClearWith(render.White)
	p.quad.Draw()
}
