package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

// OutputPass is the final pass that writes to a camera frame buffer.
type OutputPass struct {
	Input    *render.ColorBuffer
	Geometry *render.GeometryBuffer
	shader   *render.Shader
	textures *render.TextureMap
	quad     *Quad
}

// NewOutputPass creates a new output pass for the given input texture.
func NewOutputPass(input *render.ColorBuffer, gbuffer *render.GeometryBuffer) *OutputPass {
	shader := render.CompileShader(
		"output_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/output.fs")

	tx := render.NewTextureMap(shader)
	tx.Add("tex_input", input.Texture)
	tx.Add("tex_depth", gbuffer.Depth)

	return &OutputPass{
		Input:    input,
		Geometry: gbuffer,
		shader:   shader,
		textures: tx,
		quad:     NewQuad(shader),
	}
}

func (p *OutputPass) Type() render.Pass {
	return render.Postprocess
}

// DrawPass draws the input texture to the scene camera buffer.
func (p *OutputPass) Draw(scene *Scene) {
	camera := scene.Camera

	// camera settings
	camera.Use()
	render.ClearWith(scene.Camera.Clear)

	// draw
	p.shader.Use()
	p.textures.Use()
	p.quad.Draw()
}

func (p *OutputPass) Visible(c Component, args DrawArgs) bool {
	return false
}

func (p *OutputPass) Queue(c Component, args DrawArgs) {}

func (p *OutputPass) Resize(width, height int) {}
