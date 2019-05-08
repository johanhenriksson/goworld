package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render"
)

type OutputPass struct {
	Input *render.Texture
	quad  *render.RenderQuad
	mat   *render.Material
}

func NewOutputPass(input *render.Texture) *OutputPass {
	mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/screen_quad"))
	mat.AddDescriptors(render.F32_XYZUV)
	mat.AddTexture("tex_input", input)

	/* create a render quad */
	quad := render.NewRenderQuad(mat)

	return &OutputPass{
		Input: input,
		quad:  quad,
		mat:   mat,
	}
}

func (p *OutputPass) DrawPass(scene *Scene) {
	camera := scene.Camera

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// camera settings
	gl.Viewport(0, 0, int32(camera.Width), int32(camera.Height))
	gl.ClearColor(camera.Clear.R, camera.Clear.G, camera.Clear.B, camera.Clear.A)

	// draw
	gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)
	p.quad.Draw()
}
