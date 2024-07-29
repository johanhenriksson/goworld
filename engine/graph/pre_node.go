package graph

import (
	"errors"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

var ErrRecreate = errors.New("recreate renderer")

type PreDrawable interface {
	object.Component
	PreDraw(render.Args, object.Object) error
}

type preNode struct {
	*node
	target       vulkan.Target
	cameraQuery  *object.Query[*camera.Camera]
	predrawQuery *object.Query[PreDrawable]
}

func newPreNode(app vulkan.App, target vulkan.Target) *preNode {
	return &preNode{
		node:         newNode(app, "Pre", nil),
		target:       target,
		cameraQuery:  object.NewQuery[*camera.Camera](),
		predrawQuery: object.NewQuery[PreDrawable](),
	}
}

func (n *preNode) Prepare(scene object.Object, time, delta float32) (*render.Args, *swapchain.Context, error) {
	screen := render.Screen{
		Width:  n.target.Width(),
		Height: n.target.Height(),
		Scale:  n.target.Scale(),
	}

	// aquire next frame
	ctxAvailable := make(chan *swapchain.Context)
	n.app.Worker().Invoke(func() {
		context, err := n.target.Aquire()
		if err != nil {
			ctxAvailable <- nil
		} else {
			ctxAvailable <- context
		}
	})

	// ensure the default white texture is always available
	n.app.Textures().Fetch(color.White)

	// cache ticks
	n.app.Meshes().Tick()
	n.app.Textures().Tick()

	// create render arguments
	args := render.Args{}

	// find the first active camera
	if camera, exists := n.cameraQuery.Reset().First(scene); exists {
		args = camera.RenderArgs(screen)
	} else {
		args.Viewport = screen
	}

	// wait for context
	context := <-ctxAvailable
	if context == nil {
		return nil, nil, ErrRecreate
	}

	// fill in time & swapchain context
	args.Frame = context.Index
	args.Time = time
	args.Delta = delta
	args.Transform = mat4.Ident()

	// execute pre-draw pass
	objects := n.predrawQuery.Reset().Collect(scene)
	for _, object := range objects {
		object.PreDraw(args.Apply(object.Transform().Matrix()), scene)
	}

	// fire off render start signals
	var waits []command.Wait
	if context.ImageAvailable != nil {
		// why would this be nil?
		waits = []command.Wait{
			{
				Semaphore: context.ImageAvailable,
				Mask:      core1_0.PipelineStageColorAttachmentOutput,
			},
		}
	}

	// pre-node submits a dummy pass that does nothing
	// except signal that any pass without dependencies can start
	worker := n.app.Worker()
	worker.Submit(command.SubmitInfo{
		Commands: command.Empty,
		Marker:   n.Name(),
		Wait:     waits,
		Signal:   n.signals(context.Index),
	})

	return &args, context, nil
}
