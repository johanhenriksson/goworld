package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// ShadowPass renders shadow maps for lights.
type ShadowPass struct {
	Output *render.Texture
	Width  int
	Height int

	shadowmap *render.FrameBuffer
	queue     *DrawQueue
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
		Output: texture,
		Width:  size,
		Height: size,

		shadowmap: fbo,
		queue:     NewDrawQueue(),
	}
	return p
}

func (p *ShadowPass) Type() render.Pass {
	return render.Geometry
}

// Resize is called on window resize. Should update any window size-dependent buffers
func (p *ShadowPass) Resize(width, height int) {}

func (p *ShadowPass) Draw(scene *Scene) {}

// DrawLight draws a shadow pass for the given light.
func (p *ShadowPass) DrawLight(scene *Scene, light *Light) {
	if !light.Shadows {
		return
	}
	if light.Type != DirectionalLight {
		// only directional lights support shadows atm
		return
	}

	// bind shadow map depth render target
	// todo: each light needs its own shadow buffer?
	p.shadowmap.Bind()
	defer p.shadowmap.Unbind()

	render.DepthOutput(true)
	render.ClearWith(render.White)
	render.ClearDepth()

	// compute world to lightspace (light's view projection) matrix
	// todo: move to light object
	lp := light.Projection
	lv := mat4.LookAt(light.Position, vec3.Zero)
	lvp := lp.Mul(&lv)

	// draw shadow casters
	scene.CollectWithArgs(p, DrawArgs{
		Projection: lp,
		View:       lv,
		VP:         lvp,
		MVP:        lvp,
		Transform:  mat4.Ident(),
		Pass:       render.Geometry,
	})

	for _, cmd := range p.queue.items {
		drawable := cmd.Component.(DeferredDrawable)
		drawable.DrawDeferred(cmd.Args)
	}

	render.DepthOutput(false)
}

func (p *ShadowPass) Visible(c Component, args DrawArgs) bool {
	_, ok := c.(DeferredDrawable)
	return ok
}

func (p *ShadowPass) Queue(c Component, args DrawArgs) {
	p.queue.Add(c, args)
}
