package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render"
)

type ForwardDrawable interface {
	DrawForward(DrawArgs)
}

// ForwardPass holds information required to perform a forward rendering pass.
type ForwardPass struct {
	output  *render.ColorBuffer
	gbuffer *render.GeometryBuffer
	queue   *DrawQueue
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
		queue:   NewDrawQueue(),
	}
}

func (p *ForwardPass) Type() render.Pass {
	return render.Forward
}

func (p *ForwardPass) Resize(width, height int) {}

// DrawPass executes the forward pass
func (p *ForwardPass) Draw(scene *Scene) {
	scene.Camera.Use()

	// setup rendering
	render.Blend(true)
	render.BlendMultiply()
	render.CullFace(render.CullBack)

	// draw scene
	p.queue.Clear()
	scene.Collect(p)

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

	for _, cmd := range p.queue.items {
		drawable := cmd.Component.(ForwardDrawable)
		drawable.DrawForward(cmd.Args)
	}

	render.DepthOutput(true)

	render.CullFace(render.CullNone)
}

func (p *ForwardPass) Visible(c Component, args DrawArgs) bool {
	_, ok := c.(ForwardDrawable)
	return ok
}

func (p *ForwardPass) Queue(c Component, args DrawArgs) {
	p.queue.Add(c, args)
}
