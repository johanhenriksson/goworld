package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
)

// ColorPass represents a color correction pass and its settings.
type ColorPass struct {
	Input    *render.ColorBuffer
	Output   *render.ColorBuffer
	Lut      *render.Texture
	Gamma    float32
	shader   *render.Shader
	textures *render.TextureMap
	quad     *Quad
}

// NewColorPass instantiates a new color correction pass.
func NewColorPass(input *render.ColorBuffer, filter string) *ColorPass {
	// load lookup table
	lutName := fmt.Sprintf("textures/color_grading/%s.png", filter)
	lut := assets.GetTexture(lutName)

	shader := render.CompileShader(
		"color_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/color.fs")
	tx := render.NewTextureMap(shader)

	tx.Add("tex_input", input.Texture)
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

func (p *ColorPass) Type() render.Pass {
	return render.Postprocess
}

// DrawPass applies color correction to the scene
func (p *ColorPass) Draw(scene *Scene) {
	p.Output.Bind()
	defer p.Output.Unbind()

	// pass shader settings
	p.shader.Use()
	p.textures.Use()
	p.shader.Float("gamma", p.Gamma)

	render.Clear()
	p.quad.Draw()
}

func (p *ColorPass) Visible(c Component, args DrawArgs) bool {
	return false
}

func (p *ColorPass) Queue(c Component, args DrawArgs) {}

func (p *ColorPass) Resize(width, height int) {
	p.Output.Resize(width, height)
}
