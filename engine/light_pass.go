package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/render"
)

type LightPass struct {
	Material       *render.Material
	quad           *render.Quad
	Output         *render.Texture
	SSAO           *SSAOPass
	Shadows        *ShadowPass
	Ambient        render.Color
	ShadowStrength float32
	ShadowBias     float32

	fbo *render.FrameBuffer
}

func NewLightPass(input *render.GeometryBuffer) *LightPass {
	ssao := SSAOSettings{
		Samples: 32,
		Radius:  0.8,
		Bias:    0.025,
		Power:   1.5,
	}

	ssaoPass := NewSSAOPass(input, &ssao)
	shadowPass := NewShadowPass(input)

	fbo := render.CreateFrameBuffer(input.Width, input.Height)
	output := fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.FLOAT)

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
	quad := render.NewQuad(mat)

	p := &LightPass{
		fbo:            fbo,
		Output:         output,
		Material:       mat,
		quad:           quad,
		Shadows:        shadowPass,
		SSAO:           ssaoPass,
		Ambient:        render.Color4(1, 1, 1, 0.1),
		ShadowStrength: 0.3,
		ShadowBias:     0.0001,
	}
	return p
}

func (p *LightPass) DrawPass(scene *Scene) {
	// enable back face culling
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	// ssao pass
	p.SSAO.DrawPass(scene)

	p.fbo.Bind()

	/* use light pass shader */
	shader := p.Material.ShaderProgram

	/* compute camera view projection inverse */
	vp := scene.Camera.Projection.Mul4(scene.Camera.View)
	vpInv := vp.Inv()
	vInv := scene.Camera.View.Inv()

	shader.Use()
	shader.Mat4f("cameraInverse", vpInv)
	shader.Mat4f("viewInverse", vInv)
	shader.RGBA("ambient", p.Ambient)
	shader.Float("shadow_bias", p.ShadowBias)
	shader.Float("shadow_strength", p.ShadowStrength)

	// clear
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// set blending mode to additive
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)

	// draw lights
	for i, light := range scene.Lights {
		/* draw shadow pass for this light into shadow map */
		p.Shadows.DrawPass(scene, &light)

		if i == 0 {
			/* first light pass we want the shader to restore the depth buffer
			 * then, disable depth masking so that multiple lights can be drawn */
			gl.DepthMask(true)
		} else {
			gl.DepthMask(false)
		}

		p.fbo.Bind()

		/* use light pass shader */
		shader.Use()

		/* compute world to lightspace (light view projection) matrix */
		lp := light.Projection
		lv := mgl.LookAtV(light.Position, mgl.Vec3{}, mgl.Vec3{0, 1, 0}) // only for directional light
		lvp := lp.Mul4(lv)
		shader.Mat4f("light_vp", lvp)

		/* set light uniform attributes */
		shader.Vec3("light.Position", &light.Position)
		shader.Vec3("light.Color", &light.Color)
		shader.Int32("light.Type", int32(light.Type))
		shader.Float("light.Range", light.Range)
		shader.Float("light.Intensity", light.Intensity)
		shader.Float("light.attenuation.Constant", light.Attenuation.Constant)
		shader.Float("light.attenuation.Linear", light.Attenuation.Linear)
		shader.Float("light.attenuation.Quadratic", light.Attenuation.Quadratic)

		/* render light */
		// todo: draw light volumes instead of a fullscreen quad
		p.quad.Draw()

		p.fbo.Unbind()
	}

	/* reset GL state */
	gl.DepthMask(true)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
}
