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
		Scale:   1,
	})

	// create output frame buffer
	fbo := render.CreateFrameBuffer(input.Width, input.Height)
	output := fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE)

	// instantiate light pass shader
	mat := render.CreateMaterial(render.CompileShader("/assets/shaders/light_pass"))
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
	// compute world to lightspace (light view projection) matrix
	// note: this is only for directional lights
	lp := light.Projection
	lv := mat4.LookAt(light.Position, vec3.Zero)
	lvp := lp.Mul(&lv)
	p.mat.Mat4("light_vp", &lvp)

	/* set light uniform attributes */
	p.mat.Vec3("light.Position", &light.Position)
	p.mat.Vec3("light.Color", &light.Color)
	p.mat.Int32("light.Type", int32(light.Type))
	p.mat.Float("light.Range", light.Range)
	p.mat.Float("light.Intensity", light.Intensity)
	p.mat.Float("light.attenuation.Constant", light.Attenuation.Constant)
	p.mat.Float("light.attenuation.Linear", light.Attenuation.Linear)
	p.mat.Float("light.attenuation.Quadratic", light.Attenuation.Quadratic)
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
	p.mat.Mat4("cameraInverse", &vpInv)
	p.mat.Mat4("viewInverse", &vInv)

	// clear output buffer
	p.fbo.Bind()
	render.ClearWith(render.Black)

	// enable back face culling
	render.CullFace(render.CullBack)

	// enable blending
	render.Blend(false)

	// ambient light pass
	ambient := Light{
		Color:     p.Ambient.Vec3(),
		Intensity: 1,
	}
	p.setLightUniforms(&ambient)
	p.quad.Draw()

	// set blending mode to additive
	render.BlendAdditive()

	// draw lights one by one
	for _, light := range scene.Lights {
		// draw shadow pass for this light into shadow map
		p.Shadows.DrawPass(scene, &light)

		// first light pass we want the shader to restore the depth buffer
		// then, disable depth masking so that multiple lights can be drawn
		// render.DepthOutput(i == 0)

		// use light shader again
		p.mat.Use()
		p.setLightUniforms(&light)

		// render light
		// todo: draw light volumes instead of a fullscreen quad
		p.fbo.Bind()
		p.quad.Draw()
	}

	// reset GL state
	render.DepthOutput(true)
	render.Blend(false)
}
