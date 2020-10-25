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

	quad     *Quad
	fbo      *render.FrameBuffer
	shader   *render.Shader
	textures *render.TextureMap
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
	tx.Add("tex_occlusion", ssaoPass.Gaussian.Output)

	p := &LightPass{
		Output:         output,
		Shadows:        shadowPass,
		SSAO:           ssaoPass,
		Ambient:        render.Color4(0.25, 0.25, 0.25, 1),
		ShadowStrength: 0.3,
		ShadowBias:     0.0001,
		SSAOAmount:     0.5,

		quad:     NewQuad(shader),
		shader:   shader,
		textures: tx,
		fbo:      fbo,
	}

	// set up static uniforms
	shader.Use()
	shader.Float("shadow_bias", p.ShadowBias)
	shader.Float("shadow_strength", p.ShadowStrength)
	shader.Float("ssao_amount", p.SSAOAmount)

	return p
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

// DrawPass executes the deferred lighting pass.
func (p *LightPass) DrawPass(scene *Scene) {
	// ssao pass
	p.SSAO.DrawPass(scene)

	// compute camera view projection inverse
	vp := scene.Camera.Projection.Mul(&scene.Camera.View)
	vpInv := vp.Invert()
	vInv := scene.Camera.View.Invert()
	p.shader.Use()
	p.textures.Use()
	p.shader.Mat4("cameraInverse", &vpInv)
	p.shader.Mat4("viewInverse", &vInv)

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
		p.Shadows.DrawPass(scene, &light)

		// first light pass we want the shader to restore the depth buffer
		// then, disable depth masking so that multiple lights can be drawn

		// use light shader again
		p.shader.Use()
		p.textures.Use()
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
