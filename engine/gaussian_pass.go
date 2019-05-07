package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render"
)

type GaussianPass struct {
	fbo      *render.FrameBuffer
	material *render.Material
	quad     *render.RenderQuad
	Output   *render.Texture
}

func NewGaussianPass(input *render.Texture) *GaussianPass {
	fbo := render.CreateFrameBuffer(input.Width, input.Height)
	fbo.ClearColor = render.Color4(1, 0, 0, 1)
	texture := fbo.AddBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.FLOAT)

	mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/gaussian"))
	mat.AddDescriptors(render.F32_XYZUV)
	mat.AddTexture("tex_input", input)

	quad := render.NewRenderQuad(mat)

	return &GaussianPass{
		fbo:      fbo,
		material: mat,
		quad:     quad,
		Output:   texture,
	}
}

func (p *GaussianPass) DrawPass(scene *Scene) {
	p.fbo.Bind()
	p.fbo.Clear()

	p.quad.Draw()

	p.fbo.Unbind()
}
