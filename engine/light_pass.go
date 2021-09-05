package engine

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// LightPass draws the deferred lighting pass
type LightPass struct {
	GBuffer        *render.GeometryBuffer
	Output         *render.ColorBuffer
	Shadows        *ShadowPass
	Ambient        render.Color
	ShadowStrength float32
	ShadowBias     float32
	SSAOAmount     float32

	quad     *Quad
	shader   *render.Shader
	textures *render.TextureMap
}

// NewLightPass creates a new deferred lighting pass
func NewLightPass(input *render.GeometryBuffer) *LightPass {
	// child passes
	shadowPass := NewShadowPass(input)

	// instantiate light pass shader
	shader := render.CompileShader(
		"light_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/light.fs")

	// add gbuffer, shadow and ssao pass inputs
	tx := render.NewTextureMap(shader)
	tx.Add("tex_diffuse", input.Diffuse)
	tx.Add("tex_normal", input.Normal)
	tx.Add("tex_depth", input.Depth)
	tx.Add("tex_shadow", shadowPass.Output)
	// tx.Add("tex_occlusion", ssaoPass.Gaussian.Output)

	p := &LightPass{
		GBuffer:        input,
		Output:         render.NewColorBuffer(input.Width, input.Height),
		Shadows:        shadowPass,
		Ambient:        render.Color4(0.25, 0.25, 0.25, 1),
		ShadowStrength: 0.3,
		ShadowBias:     0.0001,
		SSAOAmount:     0.5,

		quad:     NewQuad(shader),
		shader:   shader,
		textures: tx,
	}

	// set up static uniforms
	shader.Use()
	shader.Float("shadow_bias", p.ShadowBias)
	shader.Float("shadow_strength", p.ShadowStrength)
	// shader.Float("ssao_amount", p.SSAOAmount)

	return p
}

func (p *LightPass) Type() render.Pass {
	return render.Lights
}

// Resize is called on window resize. Should update any window size-dependent buffers
func (p *LightPass) Resize(width, height int) {
	// p.SSAO.Resize(width, height)
	p.Output.Resize(width, height)
}

func (p *LightPass) setLightUniforms(light *Light) {
	// compute world to lightspace (light view projection) matrix
	// note: this is only for directional lights
	lp := light.Projection
	lv := mat4.LookAt(light.Position, vec3.Zero)
	lvp := lp.Mul(&lv)
	p.shader.Mat4("light_vp", &lvp)

	/* set light uniform attributes */
	p.shader.Vec3("light.Position", &light.Position)
	p.shader.Vec3("light.Color", &light.Color)
	p.shader.Int32("light.Type", int(light.Type))
	p.shader.Float("light.Range", light.Range)
	p.shader.Float("light.Intensity", light.Intensity)
	p.shader.Float("light.attenuation.Constant", light.Attenuation.Constant)
	p.shader.Float("light.attenuation.Linear", light.Attenuation.Linear)
	p.shader.Float("light.attenuation.Quadratic", light.Attenuation.Quadratic)
}

// Draw executes the deferred lighting pass.
func (p *LightPass) Draw(scene *Scene) {
	// compute camera view projection inverse
	vInv := scene.Camera.ViewInv()
	vpInv := scene.Camera.ViewProjInv()

	// clear output buffer
	p.Output.Bind()
	render.ClearWith(render.Black)

	// enable back face culling
	render.CullFace(render.CullBack)

	// enable blending
	render.Blend(false)

	p.shader.Use()
	p.textures.Use()
	p.shader.Mat4("cameraInverse", &vpInv)
	p.shader.Mat4("viewInverse", &vInv)

	render.DepthOutput(true)

	// ambient light pass
	ambient := Light{
		Color:     p.Ambient.Vec3(),
		Intensity: 1.3,
	}
	p.setLightUniforms(&ambient)
	p.quad.Draw()

	// set blending mode to additive
	render.BlendAdditive()

	render.DepthOutput(false)

	// draw lights one by one
	for _, light := range scene.Lights {
		// draw shadow pass for this light into shadow map
		p.Shadows.DrawLight(scene, &light)

		// first light pass we want the shader to restore the depth buffer
		// then, disable depth masking so that multiple lights can be drawn

		p.Output.Bind()

		// use light shader again
		p.shader.Use()
		p.textures.Use()
		p.setLightUniforms(&light)

		render.DepthOutput(true)
		// render light
		// todo: draw light volumes instead of a fullscreen quad
		p.quad.Draw()
	}

	// reset GL state
	render.DepthOutput(true)
	render.Blend(false)
}
