package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

// OutputPass is the final pass that writes to a camera frame buffer.
type OutputPass struct {
	Input    *render.Texture
	shader   *render.Shader
	textures *render.TextureMap
	quad     *Quad
}

// NewOutputPass creates a new output pass for the given input texture.
func NewOutputPass(input, depth *render.Texture) *OutputPass {
	shader := render.CompileShader(
		"output_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/output.fs")

	tx := render.NewTextureMap(shader)
	tx.Add("tex_input", input)
	tx.Add("tex_depth", depth)

	return &OutputPass{
		Input:    input,
		shader:   shader,
		textures: tx,
		quad:     NewQuad(shader),
	}
}

// DrawPass draws the input texture to the scene camera buffer.
func (p *OutputPass) DrawPass(scene *Scene) {
	camera := scene.Camera

	// camera settings
	camera.Use()
	render.ClearWith(scene.Camera.Clear)

	// draw
	p.shader.Use()
	p.textures.Use()
	p.quad.Draw()
}
