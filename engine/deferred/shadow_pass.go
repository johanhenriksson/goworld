package deferred

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_shader"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// ShadowPass renders shadow maps for lights.
type ShadowPass struct {
	Output texture.T
	Width  int
	Height int

	mat       material.T
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

	shader := gl_shader.CompileShader(
		"shadow_pass",
		"assets/shaders/pass/shadow.vs",
		"assets/shaders/util/empty.fs")

	mat := material.New("shadow_pass", shader)

	p := &ShadowPass{
		Output: fbo.Depth(),
		Width:  size,
		Height: size,

		mat:       mat,
		shadowmap: fbo,
	}
	return p
}

// DrawLight draws a shadow pass for the given light.
func (p *ShadowPass) DrawLight(objects []ShadowDrawable, lit *light.Descriptor) {
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
	render.CullFace(render.CullFront)

	args := render.Args{
		Projection: lit.Projection,
		View:       lit.View,
		VP:         lit.ViewProj,
		MVP:        lit.ViewProj,
		Transform:  mat4.Ident(),
	}

	// todo: select only objects that cast shadows
	// todo: view frustum culling based on the lights view projection

	p.mat.Use()
	for _, drawable := range objects {
		objArgs := args.Apply(drawable.Object().Transform().World())
		if err := p.mat.Mat4("mvp", objArgs.MVP); err != nil {
			panic("failed to set shadow projection")
		}
		drawable.DrawShadow(objArgs)
	}

	render.CullFace(render.CullBack)
}
