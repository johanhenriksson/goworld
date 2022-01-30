package engine

import (
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_shader"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/light"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/shader"
)

// LightPass draws the deferred lighting pass
type LightPass struct {
	GBuffer        framebuffer.Geometry
	Output         framebuffer.Color
	Shadows        *ShadowPass
	Ambient        color.T
	ShadowStrength float32
	ShadowBias     float32
	SSAOAmount     float32

	quad   *Quad
	shader shader.T
	mat    material.T
}

// NewLightPass creates a new deferred lighting pass
func NewLightPass(input framebuffer.Geometry) *LightPass {
	// child passes
	shadowPass := NewShadowPass()

	// instantiate light pass shader
	shader := gl_shader.CompileShader(
		"light_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/light.fs")

	// add gbuffer, shadow and ssao pass inputs
	mat := material.New("light_pass", shader)
	mat.Texture("tex_diffuse", input.Diffuse())
	mat.Texture("tex_normal", input.Normal())
	mat.Texture("tex_depth", input.Depth())
	mat.Texture("tex_shadow", shadowPass.Output)

	p := &LightPass{
		GBuffer:        input,
		Output:         gl_framebuffer.NewColor(input.Width(), input.Height()),
		Shadows:        shadowPass,
		Ambient:        color.RGB(0.25, 0.25, 0.25),
		ShadowStrength: 0.3,
		ShadowBias:     0.0001,
		SSAOAmount:     0.5,

		quad:   NewQuad(shader),
		shader: shader,
		mat:    mat,
	}

	// set up static uniforms
	shader.Use()
	shader.Float("shadow_bias", p.ShadowBias)
	shader.Float("shadow_strength", p.ShadowStrength)
	// shader.Float("ssao_amount", p.SSAOAmount)

	return p
}

func (p *LightPass) setLightUniforms(light *light.T) {
	// compute world to lightspace (light view projection) matrix
	// note: this is only for directional lights
	lp := light.Projection
	lv := mat4.LookAt(light.Position, vec3.Zero)
	lvp := lp.Mul(&lv)
	p.shader.Mat4("light_vp", lvp)

	/* set light uniform attributes */
	p.shader.Vec3("light.Position", light.Position)
	p.shader.RGB("light.Color", light.Color)
	p.shader.Int32("light.Type", int(light.Type))
	p.shader.Float("light.Range", light.Range)
	p.shader.Float("light.Intensity", light.Intensity)
	p.shader.Float("light.attenuation.Constant", light.Attenuation.Constant)
	p.shader.Float("light.attenuation.Linear", light.Attenuation.Linear)
	p.shader.Float("light.attenuation.Quadratic", light.Attenuation.Quadratic)
}

// Draw executes the deferred lighting pass.
func (p *LightPass) Draw(args render.Args, scene scene.T) {
	// compute camera view projection inverse
	vInv := scene.Camera().ViewInv()
	vpInv := scene.Camera().ViewProjInv()

	// clear output buffer
	p.Output.Bind()
	defer p.Output.Unbind()
	p.Output.Resize(args.Viewport.FrameWidth, args.Viewport.FrameHeight)
	render.ClearWith(scene.Camera().ClearColor())

	// enable back face culling
	render.CullFace(render.CullBack)

	// enable blending
	render.Blend(false)

	p.mat.Use()
	p.shader.Mat4("cameraInverse", vpInv)
	p.shader.Mat4("viewInverse", vInv)

	render.DepthOutput(true)

	// ambient light pass
	ambient := light.T{
		Type:      light.Ambient,
		Color:     p.Ambient,
		Intensity: 1.3,
	}
	p.setLightUniforms(&ambient)
	p.quad.Draw()

	// set blending mode to additive
	render.BlendAdditive()

	render.DepthOutput(false)

	// draw lights one by one
	for _, light := range scene.Lights() {
		// draw shadow pass for this light into shadow map
		p.Shadows.DrawLight(&light)

		// first light pass we want the shader to restore the depth buffer
		// then, disable depth masking so that multiple lights can be drawn

		p.Output.Bind()

		// use light shader again
		p.mat.Use()
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
