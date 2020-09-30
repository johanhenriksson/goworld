package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// LightPass draws the deferred lighting pass
type LightPass struct {
	Output         *render.Texture
	SSAO           *SSAOPass
	Shadows        *ShadowPass
	Ambient        render.Color
	ShadowStrength float32
	ShadowBias     float32
	SSAOAmount     float32

	fbo  *render.FrameBuffer
	mat  *render.Material
	quad *render.Quad
}

// NewLightPass creates a new deferred lighting pass
func NewLightPass(input *render.GeometryBuffer) *LightPass {
	// child passes
	shadowPass := NewShadowPass(input)
	ssaoPass := NewSSAOPass(input, &SSAOSettings{
		Samples: 16,
		Radius:  0.5,
		Bias:    0.03,
		Power:   2.0,
		Scale:   2,
	})

	// create output frame buffer
	fbo := render.CreateFrameBuffer(input.Width, input.Height)
	output := fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE)

	// instantiate light pass shader
	mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/light_pass"))
	mat.AddDescriptors(render.F32_XYZUV)

	// create full screen render quad
	quad := render.NewQuad(mat)

	// add gbuffer, shadow and ssao pass inputs
	mat.AddTexture("tex_diffuse", input.Diffuse)
	mat.AddTexture("tex_normal", input.Normal)
	mat.AddTexture("tex_depth", input.Depth)
	mat.AddTexture("tex_shadow", shadowPass.Output)
	mat.AddTexture("tex_occlusion", ssaoPass.Gaussian.Output)

	p := &LightPass{
		fbo:            fbo,
		Output:         output,
		mat:            mat,
		quad:           quad,
		Shadows:        shadowPass,
		SSAO:           ssaoPass,
		Ambient:        render.Color4(0.25, 0.25, 0.25, 1),
		ShadowStrength: 0.3,
		ShadowBias:     0.0001,
		SSAOAmount:     0.5,
	}

	// set up static uniforms
	mat.Use()
	mat.Float("shadow_bias", p.ShadowBias)
	mat.Float("shadow_strength", p.ShadowStrength)
	mat.Float("ssao_amount", p.SSAOAmount)

	return p
}

func (p *LightPass) setLightUniforms(light *Light) {
	shader := p.mat.ShaderProgram

	/* compute world to lightspace (light view projection) matrix */
	lp := light.Projection
	lv := mat4.LookAt(light.Position, vec3.Zero) // only for directional light
	lvp := lp.Mul(&lv)
	shader.Mat4f("light_vp", &lvp)

	/* set light uniform attributes */
	shader.Vec3("light.Position", &light.Position)
	shader.Vec3("light.Color", &light.Color)
	shader.Int32("light.Type", int32(light.Type))
	shader.Float("light.Range", light.Range)
	shader.Float("light.Intensity", light.Intensity)
	shader.Float("light.attenuation.Constant", light.Attenuation.Constant)
	shader.Float("light.attenuation.Linear", light.Attenuation.Linear)
	shader.Float("light.attenuation.Quadratic", light.Attenuation.Quadratic)
}

// DrawPass executes the deferred lighting pass.
func (p *LightPass) DrawPass(scene *Scene) {
	// ssao pass
	p.SSAO.DrawPass(scene)

	// compute camera view projection inverse
	vp := scene.Camera.Projection.Mul(&scene.Camera.View)
	vpInv := vp.Invert()
	vInv := scene.Camera.View.Invert()
	p.mat.Use()
	p.mat.Mat4f("cameraInverse", &vpInv)
	p.mat.Mat4f("viewInverse", &vInv)

	// clear output buffer
	p.fbo.Bind()
	p.fbo.ClearColor = scene.Camera.Clear
	p.fbo.Clear()

	// enable back face culling
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	// enable blending
	gl.Enable(gl.BLEND)

	// ambient light pass
	ambient := Light{
		Color:     p.Ambient.Vec3(),
		Intensity: 1,
	}
	gl.BlendFunc(gl.SRC_ALPHA, gl.ZERO)
	p.setLightUniforms(&ambient)
	p.quad.Draw()

	// set blending mode to additive
	gl.BlendFunc(gl.ONE, gl.ONE)

	// draw lights one by one
	for i, light := range scene.Lights {
		gl.DepthMask(i == 0)
		// draw shadow pass for this light into shadow map
		p.Shadows.DrawPass(scene, &light)

		// first light pass we want the shader to restore the depth buffer
		// then, disable depth masking so that multiple lights can be drawn

		// use light shader again
		p.mat.Use()
		p.setLightUniforms(&light)

		// render light
		// todo: draw light volumes instead of a fullscreen quad
		p.fbo.Bind()
		p.quad.Draw()
	}

	// reset GL state
	gl.DepthMask(true)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
}
