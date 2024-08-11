package graph

import (
	"errors"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/vkngwrapper/core/v2/core1_0"
)

var ErrRecreate = errors.New("recreate renderer")

type PreDrawable interface {
	object.Component
	PreDraw(draw.Args, object.Object) error
}

type preNode struct {
	*node
	target       engine.Target
	cameraQuery  *object.Query[*camera.Camera]
	predrawQuery *object.Query[PreDrawable]
}

func newPreNode(app engine.App, target engine.Target) *preNode {
	return &preNode{
		node:         newNode(app, "Pre", nil),
		target:       target,
		cameraQuery:  object.NewQuery[*camera.Camera](),
		predrawQuery: object.NewQuery[PreDrawable](),
	}
}

func (n *preNode) Prepare(scene object.Object, time, delta float32) (*draw.Args, *swapchain.Context, error) {
	viewport := draw.Viewport{
		Width:  n.target.Width(),
		Height: n.target.Height(),
		Scale:  n.target.Scale(),
	}

	// todo: cache handling does not really belong in the render graph
	// ensure the default white texture is always available
	n.app.Textures().Fetch(color.White)
	// cache ticks
	n.app.Meshes().Tick()
	n.app.Textures().Tick()

	// aquire next frame
	context, err := n.target.Aquire(n.app.Worker())
	if err != nil {
		return nil, nil, ErrRecreate
	}

	// create render arguments
	args := draw.Args{}

	// find the first active cam
	if cam, exists := n.cameraQuery.Reset().First(scene); exists {
		args.Camera = cam.Refresh(viewport)
	} else {
		args.Camera.Viewport = viewport
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

	// pre-node submits a dummy pass that does nothing
	// except signal that any pass without dependencies can start
	worker := n.app.Worker()
	worker.Submit(command.SubmitInfo{
		Commands: command.Empty,
		Marker:   n.Name(),
		Signal:   n.signals(context.Index),
		Wait: []command.Wait{
			{
				Semaphore: context.ImageAvailable,
				Mask:      core1_0.PipelineStageColorAttachmentOutput,
			},
		},
	})

	return &args, context, nil
}
