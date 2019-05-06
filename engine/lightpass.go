package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/render"
)

type LightPass struct {
	Material *render.Material
	quad     *render.RenderQuad
	Shadows  *ShadowPass
}

func NewLightPass(input *render.GeometryBuffer) *LightPass {

	shadowPass := NewShadowPass(input)

	/* use a virtual material to help with vertex attributes and textures */
	mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/voxel_light_pass"))

	/* we're going to render a simple quad, so we input
	 * position and texture coordinates */
	mat.AddDescriptor("position", gl.FLOAT, 3, 20, 0, false, false)
	mat.AddDescriptor("texcoord", gl.FLOAT, 2, 20, 12, false, false)

	/* the shader uses 3 textures from the geometry frame buffer.
	 * they are previously rendered in the geometry pass. */
	mat.AddTexture("tex_diffuse", input.Diffuse)
	mat.AddTexture("tex_normal", input.Normal)
	mat.AddTexture("tex_depth", input.Depth)
	mat.AddTexture("tex_shadow", shadowPass.Output)

	/* create a render quad */
	quad := render.NewRenderQuad()
	/* set up vertex attribute pointers */
	mat.SetupVertexPointers()

	p := &LightPass{
		Material: mat,
		quad:     quad,
		Shadows:  shadowPass,
	}
	return p
}

func (p *LightPass) DrawPass(scene *Scene) {
	/* use light pass shader */
	p.Material.Use()
	shader := p.Material.Shader

	/* compute camera view projection inverse */
	vp := scene.Camera.Projection.Mul4(scene.Camera.View)
	vp_inv := vp.Inv()
	shader.Matrix4f("cameraInverse", &vp_inv[0])

	/* clear */
	gl.ClearColor(0.12, 0.12, 0.12, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	/* set blending mode to additive */

	/* draw lights */
	gl.BlendFunc(gl.ONE, gl.ONE)

	for i, light := range scene.Lights {
		/* draw shadow pass for this light into shadow map */
		p.Shadows.DrawPass(scene, &light)

		if i == 1 {
			/* first light pass we want the shader to restore the depth buffer
			 * then, disable depth masking so that multiple lights can be drawn */
			gl.DepthMask(true)
		} else {
			gl.DepthMask(false)
		}


		/* use light pass shader */
		p.Material.Use()

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
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
