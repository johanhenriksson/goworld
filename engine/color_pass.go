package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
)

// ColorPass represents a color correction pass and its settings.
type ColorPass struct {
	Input  *render.Texture
	Output *render.Texture
	Lut    *render.Texture
	Gamma  float32
	fbo    *render.FrameBuffer
	mat    *render.Material
	quad   *render.Quad
}

// NewColorPass instantiates a new color correction pass.
func NewColorPass(input *render.Texture, filter string) *ColorPass {
	fbo := render.CreateFrameBuffer(input.Width, input.Height)
	output := fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE)

	// load lookup table
	lutName := fmt.Sprintf("textures/color_grading/%s.png", filter)
	lut := assets.GetTexture(lutName)

	// create virtual material
	mat := render.CreateMaterial("color_pass", render.CompileShader("/assets/shaders/color_pass"))
	mat.AddDescriptors(render.F32_XYZUV)
	mat.AddTexture("tex_input", input)
	mat.AddTexture("tex_lut", lut)

	quad := render.NewQuad(mat)

	return &ColorPass{
		Input:  input,
		Output: output,
		Lut:    lut,
		Gamma:  1.7,
		fbo:    fbo,
		mat:    mat,
		quad:   quad,
	}
}

// DrawPass applies color correction to the scene
func (p *ColorPass) DrawPass(scene *Scene) {
	p.fbo.Bind()
	defer p.fbo.Unbind()

	// pass shader settings
	p.mat.Use()
	p.mat.Float("gamma", p.Gamma)

	render.Clear()
	p.quad.Draw()
}
