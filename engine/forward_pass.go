package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/render"
)

type ForwardDrawable interface {
	DrawForward(DrawArgs)
}

// ForwardPass holds information required to perform a forward rendering pass.
type ForwardPass struct {
	output  *render.ColorBuffer
	gbuffer *render.GeometryBuffer
	fbo     *render.FrameBuffer
}

// NewForwardPass sets up a forward pass.
func NewForwardPass(gbuffer *render.GeometryBuffer, output *render.ColorBuffer) *ForwardPass {
	fbo := render.CreateFrameBuffer(gbuffer.Width, gbuffer.Height)
	fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, output.Texture)
	fbo.AttachBuffer(gl.COLOR_ATTACHMENT1, gbuffer.Normal)
	fbo.AttachBuffer(gl.COLOR_ATTACHMENT2, gbuffer.Position)
	fbo.AttachBuffer(gl.DEPTH_ATTACHMENT, gbuffer.Depth)

	return &ForwardPass{
		fbo:     fbo,
		output:  output,
		gbuffer: gbuffer,
	}
}

func (p *ForwardPass) Resize(width, height int) {}

// DrawPass executes the forward pass
func (p *ForwardPass) Draw(scene *Scene) {
	scene.Camera.Use()

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
	p.fbo.DrawBuffers()

	// disable depth testing
	// todo: should be disabled for transparent things, not everything
	// render.DepthOutput(false)

	query := object.NewQuery(func(c object.Component) bool {
		_, ok := c.(ForwardDrawable)
		return ok
	})
	scene.Collect(&query)

	args := scene.Camera.DrawArgs()
	for _, component := range query.Results {
		drawable := component.(ForwardDrawable)
		drawable.DrawForward(args.Apply(component.Parent().Transform()))
	}

	render.DepthOutput(true)

	render.CullFace(render.CullNone)
}
