package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_framebuffer"
	"github.com/johanhenriksson/goworld/render/framebuffer"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ForwardDrawable interface {
	DrawForward(render.Args)
}

// ForwardPass holds information required to perform a forward rendering pass.
type ForwardPass struct {
	output  framebuffer.Color
	gbuffer framebuffer.Geometry
	fbo     framebuffer.T
}

// NewForwardPass sets up a forward pass.
func NewForwardPass(gbuffer framebuffer.Geometry, output framebuffer.Color) *ForwardPass {
	// the forward pass renders into the output of the final deferred pass.
	// it reuses the normal, position and depth buffers and writes new data according to what is rendered
	// this ensures that we have complete information in those buffers for later passes
	fbo := gl_framebuffer.NewGeometry(gbuffer.Width(), gbuffer.Height())
	fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, output.Texture())
	fbo.AttachBuffer(gl.COLOR_ATTACHMENT1, gbuffer.Normal())
	fbo.AttachBuffer(gl.COLOR_ATTACHMENT2, gbuffer.Position())
	fbo.AttachBuffer(gl.DEPTH_ATTACHMENT, gbuffer.Depth())

	return &ForwardPass{
		fbo:     fbo,
		output:  output,
		gbuffer: gbuffer,
	}
}

// DrawPass executes the forward pass
func (p *ForwardPass) Draw(args render.Args, scene scene.T) {

	// setup rendering
	render.Blend(true)
	render.BlendMultiply()
	render.CullFace(render.CullBack)

	// todo: depth sorting
	// there is finally a decent way of doing it!!
	// now we just need a way to compute the distance from an object to the camera
	// ... and a way to sort the queue

	// todo: frustum culling
	// lets not draw stuff thats behind us at the very least
	// ... things need bounding boxes though.

	p.fbo.Bind()
	defer p.fbo.Unbind()
	p.fbo.Resize(args.Viewport.FrameWidth, args.Viewport.FrameHeight)
	p.fbo.DrawBuffers()

	// disable depth testing
	// todo: should be disabled for transparent things, not everything
	// render.DepthOutput(false)

	query := object.NewQuery(func(c object.Component) bool {
		_, ok := c.(ForwardDrawable)
		return ok
	})
	scene.Collect(&query)

	for _, component := range query.Results {
		drawable := component.(ForwardDrawable)
		drawable.DrawForward(args.Apply(component.Object().Transform().World()))
	}

	render.DepthOutput(true)

	render.CullFace(render.CullNone)
}
