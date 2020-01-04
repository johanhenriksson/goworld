package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render"
)

type OutputPass struct {
	Input *render.Texture
	quad  *render.Quad
	mat   *render.Material
}

func NewOutputPass(input *render.Texture) *OutputPass {
	mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/screen_quad"))
	mat.AddDescriptors(render.F32_XYZUV)
	mat.AddTexture("tex_input", input)

	/* create a render quad */
	quad := render.NewQuad(mat)

	return &OutputPass{
		Input: input,
		quad:  quad,
		mat:   mat,
	}
}

func (p *OutputPass) DrawPass(scene *Scene) {
	camera := scene.Camera

	// camera settings
	camera.Use()

	// draw
	gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)
	p.quad.Draw()
}
