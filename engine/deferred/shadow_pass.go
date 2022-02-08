package deferred

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// ShadowPass renders shadow maps for lights.
type ShadowPass struct {
	Output texture.T
	Width  int
	Height int

	shadowmap framebuffer.Depth
}

// NewShadowPass creates a new shadow pass
func NewShadowPass(size int) *ShadowPass {
	fbo := gl_framebuffer.NewDepth(size, size)

	// set the shadow buffer texture to clamp to a white border so that samples
	// outside the map do not fall in shadow.
	fbo.Depth().Bind()
	border := []float32{1, 1, 1, 1}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &border[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	p := &ShadowPass{
		Output: fbo.Depth(),
		Width:  size,
		Height: size,

		shadowmap: fbo,
	}
	return p
}

// DrawLight draws a shadow pass for the given light.
func (p *ShadowPass) DrawLight(scene object.T, lit *light.Descriptor) {
	if !lit.Shadows {
		return
	}
	if lit.Type != light.Directional {
		// only directional lights support shadows atm
		return
	}

	// bind shadow map depth render target
	p.shadowmap.Bind()
	defer p.shadowmap.Unbind()

	render.DepthOutput(true)
	render.ClearDepth()

	// use front-face culling while rendering shadows to mitigate panning
	// but it seems to cause problems??
	//render.CullFace(render.CullFront)

	args := render.Args{
		Projection: lit.Projection,
		View:       lit.View,
		VP:         lit.ViewProj,
		MVP:        lit.ViewProj,
		Transform:  mat4.Ident(),
	}

	// todo: select only objects that cast shadows
	// todo: view frustum culling based on the lights view projection

	objects := query.Any().
		Where(query.Is[DeferredDrawable]).
		Collect(scene)

	// todo: draw objects with a simplified shader that only outputs depth information

	for _, component := range objects {
		drawable := component.(DeferredDrawable)
		drawable.DrawDeferred(args.Apply(component.Object().Transform().World()))
	}

	render.CullFace(render.CullBack)
}
