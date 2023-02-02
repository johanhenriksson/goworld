package graph

import (
	"errors"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

var ErrNoCamera = errors.New("no active camera")
var ErrRecreate = errors.New("recreate renderer")

type PreDrawable interface {
	object.T
	PreDraw(render.Args, object.T) error
}

type PreNode interface {
	Node
	Prepare(object.T) (*render.Args, error)
}

type preNode struct {
	*node
}

func newPreNode(target vulkan.Target) PreNode {
	return &preNode{
		node: newNode(target, "Pre", nil),
	}
}

func (n *preNode) Prepare(scene object.T) (*render.Args, error) {
	screen := render.Screen{
		Width:  n.target.Width(),
		Height: n.target.Height(),
		Scale:  n.target.Scale(),
	}

	// find the first active camera
	camera := object.Query[camera.T]().First(scene)
	if camera == nil {
		return nil, ErrNoCamera
	}

	// aquire next frame
	context, err := n.target.Aquire()
	if err != nil {
		return nil, ErrRecreate
	}

	// create render arguments
	args := render.Args{
		Context:    context,
		Viewport:   screen,
		Projection: camera.Projection(),
		View:       camera.View(),
		VP:         camera.ViewProj(),
		MVP:        camera.ViewProj(),
		Position:   camera.Transform().WorldPosition(),
		Clear:      camera.ClearColor(),
		Forward:    camera.Transform().Forward(),
		Transform:  mat4.Ident(),
	}

	// execute pre-draw pass
	objects := object.Query[PreDrawable]().Collect(scene)
	for _, object := range objects {
		object.PreDraw(args.Apply(object.Transform().World()), scene)
	}

	// fire off render start signals
	var waits []command.Wait
	if args.Context.ImageAvailable != nil {
		waits = []command.Wait{
			{
				Semaphore: args.Context.ImageAvailable,
				Mask:      core1_0.PipelineStageColorAttachmentOutput,
			},
		}
	}

	worker := n.target.Worker(context.Index)
	worker.Submit(command.SubmitInfo{
		Marker: n.Name(),
		Wait:   waits,
		Signal: n.signals(context.Index),
	})

	return &args, nil
}
