package deferred

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/engine/screen_quad"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/framebuffer"
)

// LightPass draws the deferred lighting pass
type LightPass struct {
	GBuffer        framebuffer.Geometry
	Output         framebuffer.Color
	Shadows        *ShadowPass
	Ambient        color.T
	ShadowStrength float32
	ShadowBias     float32

	quad   screen_quad.T
	shader LightShader
}

// NewLightPass creates a new deferred lighting pass
func NewLightPass(input framebuffer.Geometry) *LightPass {
	// child passes
	shadowPass := NewShadowPass()

	// instantiate light pass shader
	shader := NewLightShader(input)

	p := &LightPass{
		GBuffer:        input,
		Output:         gl_framebuffer.NewColor(input.Width(), input.Height()),
		Shadows:        shadowPass,
		Ambient:        color.RGB(0.25, 0.25, 0.25),
		ShadowStrength: 0.3,
		ShadowBias:     0.0001,

		quad:   screen_quad.New(shader),
		shader: shader,
	}

	// set up static uniforms
	shader.Use()
	shader.SetShadowMap(shadowPass.Output)
	shader.SetShadowStrength(p.ShadowStrength)
	shader.SetShadowBias(p.ShadowBias)

	return p
}

// Draw executes the deferred lighting pass.
func (p *LightPass) Draw(args render.Args, scene scene.T) {
	// clear output buffer
	p.Output.Bind()
	defer p.Output.Unbind()
	p.Output.Resize(args.Viewport.FrameWidth, args.Viewport.FrameHeight)
	render.ClearWith(scene.Camera().ClearColor())

	// enable back face culling
	render.CullFace(render.CullBack)

	// enable blending
	render.Blend(false)

	p.shader.Use()
	p.shader.SetCamera(scene.Camera())

	render.DepthOutput(true)

	// ambient light pass
	p.drawLight(light.Descriptor{
		Type:      light.Ambient,
		Color:     p.Ambient,
		Intensity: 1.3,
	})

	// set blending mode to additive
	render.BlendAdditive()

	render.DepthOutput(false)

	lights := object.NewQuery().
		Where(IsLight).
		Collect(scene)

	for _, component := range lights {
		light := component.(light.T)
		desc := light.LightDescriptor()

		// draw shadow pass for this light into shadow map
		p.Shadows.DrawLight(&desc)

		// perhaps depth output should be toggled for multiple lights?
		// old code indicates so, but everything seems to work ok

		p.Output.Bind()

		p.drawLight(desc)
	}

	// reset GL state
	render.DepthOutput(true)
	render.Blend(false)
}

func (p *LightPass) drawLight(desc light.Descriptor) {
	p.shader.SetLightDescriptor(desc)

	// todo: draw light volumes instead of a fullscreen quad
	p.quad.Draw()
}

func IsLight(c object.Component) bool {
	_, ok := c.(light.T)
	return ok
}
