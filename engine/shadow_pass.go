package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// ShadowPass renders shadow maps for lights.
type ShadowPass struct {
	Output    *render.Texture
	Width     int
	Height    int
	shadowmap *render.FrameBuffer
}

// NewShadowPass creates a new shadow pass
func NewShadowPass(input *render.GeometryBuffer) *ShadowPass {
	size := 4096
	fbo := render.CreateFrameBuffer(size, size)
	texture := fbo.AttachBuffer(gl.DEPTH_ATTACHMENT, gl.DEPTH_COMPONENT24, gl.DEPTH_COMPONENT, gl.FLOAT)

	// set the shadow buffer texture to clamp to a white border so that samples
	// outside the map do not fall in shadow.
	border := []float32{1, 1, 1, 1}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &border[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	p := &ShadowPass{
		shadowmap: fbo,
		Output:    texture,
		Width:     size,
		Height:    size,
	}
	return p
}

// DrawPass draws a shadow pass for the given light.
func (sp *ShadowPass) DrawPass(scene *Scene, light *Light) {
	if !light.Shadows {
		return
	}
	if light.Type != DirectionalLight {
		// only directional lights support shadows atm
		return
	}

	// bind shadow map depth render target
	// todo: each light needs its own shadow buffer?
	sp.shadowmap.Bind()
	defer sp.shadowmap.Unbind()

	render.DepthOutput(true)
	render.ClearWith(render.White)
	render.ClearDepth()

	// compute world to lightspace (light's view projection) matrix
	// todo: move to light object
	p := light.Projection
	v := mat4.LookAt(light.Position, vec3.Zero)
	vp := p.Mul(&v)

	// draw shadow casters
	scene.Draw(DrawArgs{
		Projection: p,
		View:       v,
		VP:         vp,
		MVP:        vp,
		Transform:  mat4.Ident(),
		Pass:       DrawGeometry,
	})

	render.DepthOutput(false)
}
