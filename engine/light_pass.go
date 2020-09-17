package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
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
	output := fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.FLOAT)

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
		Ambient:        render.Color4(0.1, 0.1, 0.1, 1),
		ShadowStrength: 0.3,
		ShadowBias:     0.0001,
		SSAOAmount:     0.75,
	}

	// set up static uniforms
	mat.Use()
	mat.RGBA("ambient", p.Ambient)
	mat.Float("shadow_bias", p.ShadowBias)
	mat.Float("shadow_strength", p.ShadowStrength)
	mat.Float("ssao_amount", p.SSAOAmount)

	return p
}

func (p *LightPass) setLightUniforms(light *Light) {
	shader := p.mat.ShaderProgram

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
}

// DrawPass executes the deferred lighting pass.
func (p *LightPass) DrawPass(scene *Scene) {
	// ssao pass
	p.SSAO.DrawPass(scene)

	// compute camera view projection inverse
	vp := scene.Camera.Projection.Mul4(scene.Camera.View)
	vpInv := vp.Inv()
	vInv := scene.Camera.View.Inv()
	p.mat.Use()
	p.mat.Mat4f("cameraInverse", vpInv)
	p.mat.Mat4f("viewInverse", vInv)

	// clear output buffer
	p.fbo.Bind()
	p.fbo.Clear()

	// set blending mode to additive
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)

	// enable back face culling
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	// draw lights
	for i, light := range scene.Lights {
		// draw shadow pass for this light into shadow map
		p.Shadows.DrawPass(scene, &light)

		// first light pass we want the shader to restore the depth buffer
		// then, disable depth masking so that multiple lights can be drawn
		gl.DepthMask(i == 0)

		// use light shader again
		p.mat.Use()
		p.setLightUniforms(&light)

		// render light
		// todo: draw light volumes instead of a fullscreen quad
		p.fbo.Bind()
		p.quad.Draw()
		p.fbo.Unbind()
	}

	/* reset GL state */
	gl.DepthMask(true)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
}
