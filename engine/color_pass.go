package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
)

// ColorPass represents a color correction pass and its settings.
type ColorPass struct {
	Input    *render.Texture
	Output   *render.Texture
	Lut      *render.Texture
	Gamma    float32
	fbo      *render.FrameBuffer
	shader   *render.Shader
	textures *render.TextureMap
	quad     *Quad
}

// NewColorPass instantiates a new color correction pass.
func NewColorPass(input *render.Texture, filter string) *ColorPass {
	fbo := render.CreateFrameBuffer(input.Width, input.Height)
	output := fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE)

	// load lookup table
	lutName := fmt.Sprintf("textures/color_grading/%s.png", filter)
	lut := assets.GetTexture(lutName)

	shader := render.CompileShader(
		"color_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/color.fs")
	tx := render.NewTextureMap(shader)

	tx.Add("tex_input", input)
	tx.Add("tex_lut", lut)

	return &ColorPass{
		Input:  input,
		Output: output,
		Lut:    lut,
		Gamma:  1.7,

		fbo:      fbo,
		quad:     NewQuad(shader),
		textures: tx,
		shader:   shader,
	}
}

// DrawPass applies color correction to the scene
func (p *ColorPass) DrawPass(scene *Scene) {
	p.fbo.Bind()
	defer p.fbo.Unbind()

	// pass shader settings
	p.shader.Use()
	p.textures.Use()
	p.shader.Float("gamma", p.Gamma)

	render.Clear()
	p.quad.Draw()
}
