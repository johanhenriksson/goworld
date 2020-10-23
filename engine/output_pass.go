package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

// OutputPass is the final pass that writes to a camera frame buffer.
type OutputPass struct {
	Input *render.Texture
	quad  *Quad
	mat   *render.Material
}

// NewOutputPass creates a new output pass for the given input texture.
func NewOutputPass(input, depth *render.Texture) *OutputPass {
	mat := render.CreateMaterial("output_pass", render.CompileShader(
		"output_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/output.fs"))
	mat.AddDescriptors(render.F32_XYZUV)
	mat.AddTexture("tex_input", input)
	mat.AddTexture("tex_depth", depth)

	/* create a render quad */
	quad := NewQuad(mat)

	return &OutputPass{
		Input: input,
		quad:  quad,
		mat:   mat,
	}
}

// DrawPass draws the input texture to the scene camera buffer.
func (p *OutputPass) DrawPass(scene *Scene) {
	camera := scene.Camera

	// camera settings
	camera.Use()
	render.ClearWith(scene.Camera.Clear)

	// draw
	p.quad.Draw()
}
