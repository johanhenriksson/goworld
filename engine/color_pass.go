package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
	glframebuf "github.com/johanhenriksson/goworld/render/backend/gl/framebuffer"
	glshader "github.com/johanhenriksson/goworld/render/backend/gl/shader"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
)

// ColorPass represents a color correction pass and its settings.
type ColorPass struct {
	Input  framebuffer.Color
	Output framebuffer.Color
	AO     texture.T
	Lut    texture.T
	Gamma  float32
	shader shader.T
	mat    material.T
	quad   *Quad
}

// NewColorPass instantiates a new color correction pass.
func NewColorPass(input framebuffer.Color, filter string, ssao texture.T) *ColorPass {
	// load lookup table
	lutName := fmt.Sprintf("textures/color_grading/%s.png", filter)
	lut := assets.GetTexture(lutName)

	shader := glshader.CompileShader(
		"color_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/color.fs")

	mat := material.New("color_pass", shader)
	mat.Texture("tex_input", input.Texture())
	mat.Texture("tex_ssao", ssao)
	mat.Texture("tex_lut", lut)

	return &ColorPass{
		Input:  input,
		Output: glframebuf.NewColor(input.Width(), input.Height()),
		Lut:    lut,
		Gamma:  1.8,

		quad:   NewQuad(shader),
		mat:    mat,
		shader: shader,
	}
}

// DrawPass applies color correction to the scene
func (p *ColorPass) Draw(args render.Args, scene scene.T) {
	p.Output.Bind()
	defer p.Output.Unbind()
	p.Output.Resize(args.Viewport.FrameWidth, args.Viewport.FrameHeight)

	// pass shader settings
	p.mat.Use()
	p.mat.Float("gamma", p.Gamma)

	render.Clear()
	render.Blend(true)
	render.BlendMultiply()
	p.quad.Draw()
}

func (p *ColorPass) Resize(width, height int) {
	p.Output.Resize(width, height)
}
