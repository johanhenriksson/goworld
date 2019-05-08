package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/render"
)

type LightPass struct {
	Material *render.Material
	quad     *render.RenderQuad
	SSAO     *SSAOPass
	Shadows  *ShadowPass
	Ambient  render.Color
}

func NewLightPass(input *render.GeometryBuffer) *LightPass {
	ssao := SSAOSettings{
		Samples: 64,
		Radius:  1.5,
		Bias:    0.1,
		Power:   2.5,
	}

	ssaoPass := NewSSAOPass(input, &ssao)
	shadowPass := NewShadowPass(input)

	/* use a virtual material to help with vertex attributes and textures */
	mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/ssao_light_pass"))

	/* we're going to render a simple quad, so we input
	 * position and texture coordinates */
	//mat.AddDescriptor("position", gl.FLOAT, 3, 20, 0, false, false)
	//mat.AddDescriptor("texcoord", gl.FLOAT, 2, 20, 12, false, false)
	mat.AddDescriptors(render.F32_XYZUV)

	/* the shader uses 3 textures from the geometry frame buffer.
	 * they are previously rendered in the geometry pass. */
	mat.AddTexture("tex_diffuse", input.Diffuse)
	mat.AddTexture("tex_normal", input.Normal)
	mat.AddTexture("tex_depth", input.Depth)
	mat.AddTexture("tex_shadow", shadowPass.Output)
	mat.AddTexture("tex_occlusion", ssaoPass.Gaussian.Output)

	/* create a render quad */
	quad := render.NewRenderQuad(mat)

	p := &LightPass{
		Material: mat,
		quad:     quad,
		Shadows:  shadowPass,
		SSAO:     ssaoPass,
		Ambient:  render.Color4(1, 1, 1, 0.1),
	}
	return p
}

func (p *LightPass) DrawPass(scene *Scene) {
	// ssao pass
	p.SSAO.DrawPass(scene)

	/* use light pass shader */
	shader := p.Material.Shader

	/* compute camera view projection inverse */
	vp := scene.Camera.Projection.Mul4(scene.Camera.View)
	vpInv := vp.Inv()
	vInv := scene.Camera.View.Inv()

	shader.Use()
	shader.Matrix4f("cameraInverse", &vpInv[0])
	shader.Matrix4f("viewInverse", &vInv[0])
	shader.RGBA("ambient", p.Ambient)

	/* clear */
	clr := scene.Camera.Clear
	gl.ClearColor(clr.R, clr.G, clr.B, clr.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	/* set blending mode to additive */

	/* draw lights */
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.BLEND)

	for i, light := range scene.Lights {
		/* draw shadow pass for this light into shadow map */
		p.Shadows.DrawPass(scene, &light)

		if i == 1 {
			/* first light pass we want the shader to restore the depth buffer
			 * then, disable depth masking so that multiple lights can be drawn */
			gl.Enable(gl.BLEND)
			gl.BlendFunc(gl.ONE, gl.ONE)
			gl.DepthMask(true)
		} else {
			gl.DepthMask(false)
		}

		/* use light pass shader */
		shader.Use()

		/* compute world to lightspace (light view projection) matrix */
		lp := light.Projection
		lv := mgl.LookAtV(light.Position, mgl.Vec3{}, mgl.Vec3{0, 1, 0}) // only for directional light
		lvp := lp.Mul4(lv)
		shader.Matrix4f("light_vp", &lvp[0])

		/* set light uniform attributes */
		shader.Vec3("light.Position", &light.Position)
		shader.Vec3("light.Color", &light.Color)
		shader.Int32("light.Type", int32(light.Type))
		shader.Float("light.Range", light.Range)
		shader.Float("light.attenuation.Constant", light.Attenuation.Constant)
		shader.Float("light.attenuation.Linear", light.Attenuation.Linear)
		shader.Float("light.attenuation.Quadratic", light.Attenuation.Quadratic)

		/* render light */
		gl.Viewport(0, 0, int32(scene.Camera.Width), int32(scene.Camera.Height))
		p.quad.Draw()
	}

	/* reset GL state */
	gl.DepthMask(true)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
