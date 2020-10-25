package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

type ForwardDrawable interface {
	DrawForward(DrawArgs)
}

// ForwardPass holds information required to perform a forward rendering pass.
type ForwardPass struct {
	Output *render.ColorBuffer
	queue  *DrawQueue
}

// NewForwardPass sets up a forward pass.
func NewForwardPass(output *render.ColorBuffer) *ForwardPass {
	return &ForwardPass{
		Output: output,
		queue:  NewDrawQueue(),
	}
}

func (p *ForwardPass) Type() render.Pass {
	return render.Forward
}

func (p *ForwardPass) Resize(width, height int) {}

// DrawPass executes the forward pass
func (p *ForwardPass) Draw(scene *Scene) {
	scene.Camera.Use()
	p.Output.Bind()

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

	for _, cmd := range p.queue.items {
		drawable := cmd.Component.(ForwardDrawable)
		drawable.DrawForward(cmd.Args)
	}
}

func (p *ForwardPass) Visible(c Component, args DrawArgs) bool {
	_, ok := c.(ForwardDrawable)
	return ok
}

func (p *ForwardPass) Queue(c Component, args DrawArgs) {
	p.queue.Add(c, args)
}
