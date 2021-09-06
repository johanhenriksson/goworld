package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
)

// ColorPass represents a color correction pass and its settings.
type ColorPass struct {
	Input    *render.ColorBuffer
	Output   *render.ColorBuffer
	AO       *render.Texture
	Lut      *render.Texture
	Gamma    float32
	shader   *render.Shader
	textures *render.TextureMap
	quad     *Quad
}

// NewColorPass instantiates a new color correction pass.
func NewColorPass(input *render.ColorBuffer, filter string, ssao *render.Texture) *ColorPass {
	// load lookup table
	lutName := fmt.Sprintf("textures/color_grading/%s.png", filter)
	lut := assets.GetTexture(lutName)

	shader := render.CompileShader(
		"color_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/color.fs")
	tx := render.NewTextureMap(shader)

	tx.Add("tex_input", input.Texture)
	tx.Add("tex_ssao", ssao)
	tx.Add("tex_lut", lut)

	return &ColorPass{
		Input:  input,
		Output: render.NewColorBuffer(input.Width, input.Height),
		Lut:    lut,
		Gamma:  1.8,

		quad:     NewQuad(shader),
		textures: tx,
		shader:   shader,
	}
}

// DrawPass applies color correction to the scene
func (p *ColorPass) Draw(scene scene.T) {
	p.Output.Bind()
	defer p.Output.Unbind()

	// pass shader settings
	p.shader.Use()
	p.textures.Use()
	p.shader.Float("gamma", p.Gamma)

	render.Clear()
	render.Blend(true)
	render.BlendMultiply()
	p.quad.Draw()
}

func (p *ColorPass) Resize(width, height int) {
	p.Output.Resize(width, height)
}
