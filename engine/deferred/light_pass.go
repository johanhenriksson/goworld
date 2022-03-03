package deferred

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/screen_quad"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl"
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
	ShadowSoft     bool

	quad   screen_quad.T
	shader LightShader
}

// NewLightPass creates a new deferred lighting pass
func NewLightPass(input framebuffer.Geometry, meshes cache.Meshes) *LightPass {
	shadowsize := 4096

	// child passes
	shadowPass := NewShadowPass(shadowsize, meshes)

	// instantiate light pass shader
	shader := NewLightShader(input)

	p := &LightPass{
		GBuffer:        input,
		Output:         gl_framebuffer.NewColor(input.Width(), input.Height()),
		Shadows:        shadowPass,
		Ambient:        color.RGB(0.15, 0.15, 0.15),
		ShadowStrength: 0.8,
		ShadowBias:     0.001,
		ShadowSize:     shadowsize,
		ShadowSoft:     true,

		quad:   screen_quad.New(shader),
		shader: shader,
	}

	// set up static uniforms
	shader.Use()
	shader.SetShadowMap(shadowPass.Output)
	shader.SetShadowStrength(p.ShadowStrength)
	shader.SetShadowBias(p.ShadowBias)
	shader.SetShadowSoft(p.ShadowSoft)

	return p
}

// Draw executes the deferred lighting pass.
func (p *LightPass) Draw(args render.Args, scene object.T) {
	// clear output buffer
	p.Output.Bind()
	defer p.Output.Unbind()
	p.Output.Resize(args.Viewport.Width, args.Viewport.Height)
	render.ClearWith(args.Clear)

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

	// collect all shadow casting objects
	drawables := query.New[mesh.T]().Where(castsShadows).Collect(scene)

	lights := query.New[light.T]().Collect(scene)
	for _, lit := range lights {
		desc := lit.LightDescriptor()

		// fit light projection matrix to the current camera frustum
		p.FitLightToCamera(args.VP.Invert(), &desc)

		// draw shadow pass for this light into the shadow map
		p.Shadows.DrawLight(drawables, &desc)

		// accumulate light from this source
		p.drawLight(desc)
	}

	if err := gl.GetError(); err != nil {
		fmt.Println("uncaught GL error during light pass:", err)
	}
}

func (p *LightPass) drawLight(desc light.Descriptor) {
	p.Output.Bind()
	p.shader.Use()
	p.shader.SetLightDescriptor(desc)

	// todo: draw light volumes instead of a fullscreen quad
	p.quad.Draw()
}

var maxExtent = float32(0)

func (p *LightPass) FitLightToCamera(vpi mat4.T, desc *light.Descriptor) {
	if desc.Type != light.Directional {
		return
	}

	// ideally this should be moved into the directional light
	// however, it does depend on several awkward things:
	//  - camera view projection
	//  - shadow map dimensions
	//  - scene bounds/AABB (later)
	//
	// might not be worth fixing this until we have proper cascading shadow maps

	// projection matrix
	// create a frustum in light view space
	cameraClipToLightView := desc.View.Mul(&vpi)
	lfst := camera.NewFrustum(cameraClipToLightView)

	// calculate the max extent in X/Y, so that we can create a square projection
	// additionally, this number should be stable over time, use the historical maximum
	// this avoids shimmering
	extent := math.Max(lfst.Max.X-lfst.Min.X, lfst.Max.Y-lfst.Min.Y) * 0.5
	maxExtent = math.Max(maxExtent, extent)

	// snap the orthogonal bounds to the shadowmap texel size avoid shadow shimmering
	snap := 2 * maxExtent / float32(p.ShadowSize)

	// todo: we can further optimize the Z bounds
	// project the scenes AABB into light space and use the min/max Z values
	// lsft.Min.Z / lfst.Max.Z should work but doesnt?
	zmin := float32(-100)
	zmax := float32(100)

	desc.Projection = mat4.Orthographic(
		math.Snap(lfst.Center.X-maxExtent, snap), math.Snap(lfst.Center.X+maxExtent, snap),
		math.Snap(lfst.Center.Y-maxExtent, snap), math.Snap(lfst.Center.Y+maxExtent, snap),
		zmin, zmax)

	// finally, update light view projection
	desc.ViewProj = desc.Projection.Mul(&desc.View)
}

func castsShadows(m mesh.T) bool {
	return m.CastShadows()
}
