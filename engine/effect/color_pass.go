package effect

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/screen_quad"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_shader"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
)

// ColorPass represents a color correction pass and its settings.
type ColorPass struct {
	Input     framebuffer.Color
	Output    framebuffer.Color
	Lut       texture.T
	Occlusion texture.T
	Gamma     float32

	shader shader.T
	mat    material.T
	quad   screen_quad.T
}

// NewColorPass instantiates a new color correction pass.
// this pass also mixes in occlusion output to save an additional full screen pass. a little bit gross
func NewColorPass(input framebuffer.Color, filter string, occlusion texture.T) *ColorPass {
	// load lookup table
	lutName := fmt.Sprintf("textures/color_grading/%s.png", filter)
	lut := assets.GetTexture(lutName)

	shader := gl_shader.CompileShader(
		"color_pass",
		"assets/shaders/pass/postprocess.vs",
		"assets/shaders/pass/color.fs")

	mat := material.New("color_pass", shader)
	mat.Texture("tex_input", input.Texture())
	mat.Texture("tex_ssao", occlusion)
	mat.Texture("tex_lut", lut)

	return &ColorPass{
		Input:     input,
		Output:    gl_framebuffer.NewColor(input.Width(), input.Height()),
		Lut:       lut,
		Occlusion: occlusion,
		Gamma:     1.8,

		quad:   screen_quad.New(shader),
		mat:    mat,
		shader: shader,
	}
}

// DrawPass applies color correction to the scene
func (p *ColorPass) Draw(args render.Args, scene object.T) {
	p.Output.Bind()
	defer p.Output.Unbind()
	p.Output.Resize(p.Input.Width(), p.Input.Height())

	// pass shader settings
	p.mat.Use()
	p.mat.Float("gamma", p.Gamma)

	p.quad.Draw()
}

func (p *ColorPass) Resize(width, height int) {
	p.Output.Resize(width, height)
}
