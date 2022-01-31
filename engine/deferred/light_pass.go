package deferred

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/engine/screen_quad"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
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
	ShadowSize     int

	quad   screen_quad.T
	shader LightShader
}

// NewLightPass creates a new deferred lighting pass
func NewLightPass(input framebuffer.Geometry) *LightPass {
	shadowsize := 2048

	// child passes
	shadowPass := NewShadowPass(shadowsize)

	// instantiate light pass shader
	shader := NewLightShader(input)

	p := &LightPass{
		GBuffer:        input,
		Output:         gl_framebuffer.NewColor(input.Width(), input.Height()),
		Shadows:        shadowPass,
		Ambient:        color.RGB(0.25, 0.25, 0.25),
		ShadowStrength: 0.8,
		ShadowBias:     0.0001,
		ShadowSize:     shadowsize,

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
	// fmt.Println(cmin, cmax)

	// clear output buffer
	p.Output.Bind()
	defer p.Output.Unbind()
	p.Output.Resize(args.Viewport.FrameWidth, args.Viewport.FrameHeight)
	render.ClearWith(scene.Camera().ClearColor())

	// enable back face culling
	render.CullFace(render.CullBack)

	p.shader.Use()
	//p.shader.SetCamera(scene.Camera())
	p.shader.Mat4("cameraInverse", args.VP.Invert())
	p.shader.Mat4("viewInverse", args.View.Invert())

	// disable blending for the first light
	// we are drawing on a non-black background (camera clear color)
	// so we dont want to add to it. perhaps the clear color should be added later
	// this only works when the first light is the ambient light pass, since it lights everything
	render.Blend(false)

	// ambient light pass
	p.drawLight(light.Descriptor{
		Type:      light.Ambient,
		Color:     p.Ambient,
		Intensity: 1.3,
	})

	// accumulate the light from the non-ambient light sources
	render.BlendAdditive()

	lights := object.NewQuery().
		Where(IsLight).
		Collect(scene)

	for _, component := range lights {
		light := component.(light.T)
		desc := light.LightDescriptor()

		// fit light projection matrix to the current camera frustum
		p.FitLightToCamera(scene.Camera(), &desc)

		// draw shadow pass for this light into the shadow map
		p.Shadows.DrawLight(scene, &desc)

		// accumulate light from this source
		p.drawLight(desc)
	}
}

func (p *LightPass) drawLight(desc light.Descriptor) {
	p.Output.Bind()
	p.shader.Use()
	p.shader.SetLightDescriptor(desc)

	// todo: draw light volumes instead of a fullscreen quad
	p.quad.Draw()
}

func IsLight(c object.Component) bool {
	_, ok := c.(light.T)
	return ok
}

var maxExtent = float32(0)

func (p *LightPass) FitLightToCamera(cam camera.T, desc *light.Descriptor) {
	if desc.Type != light.Directional {
		return
	}

	fst := cam.Frustum()

	// view matrix
	target := fst.NTL.Add(fst.NTR).Add(fst.NBL).Add(fst.NBR).Add(fst.FTL).Add(fst.FTR).Add(fst.FBL).Add(fst.FBR).Scaled(1 / 8.0)

	// projection

	fst.NTL = desc.View.TransformPoint(fst.NTL)
	fst.NTR = desc.View.TransformPoint(fst.NTR)
	fst.NBL = desc.View.TransformPoint(fst.NBL)
	fst.NBR = desc.View.TransformPoint(fst.NBR)
	fst.FTL = desc.View.TransformPoint(fst.FTL)
	fst.FTR = desc.View.TransformPoint(fst.FTR)
	fst.FBL = desc.View.TransformPoint(fst.FBL)
	fst.FBR = desc.View.TransformPoint(fst.FBR)
	lmin, lmax := fst.Bounds()

	extent := math.Max(lmax.X-lmin.X, lmax.Y-lmin.Y) * 0.5
	maxExtent = math.Max(maxExtent, extent)

	center := lmax.Add(lmin).Scaled(0.5)

	// snap to avoid shadow swimming
	// does not seem to work
	snap := 2.0 * maxExtent / float32(p.ShadowSize)
	center.X = math.Snap(center.X, snap)
	center.Y = math.Snap(center.Y, snap)
	center.Z = math.Snap(center.Z, snap)

	desc.Projection = mat4.Orthographic(
		center.X-maxExtent, center.X+maxExtent,
		center.Y-maxExtent, center.Y+maxExtent,
		lmin.Z, lmax.Z)

	desc.View = mat4.LookAt(target.Add(desc.Position.Normalized()), target)

	desc.ViewProj = desc.Projection.Mul(&desc.View)
}
