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
	Prepare(scene object.T, time, delta float32) (*render.Args, error)
}

type preNode struct {
	*node
}

func newPreNode(app vulkan.App) PreNode {
	return &preNode{
		node: newNode(app, "Pre", nil),
	}
}

func (n *preNode) Prepare(scene object.T, time, delta float32) (*render.Args, error) {
	screen := render.Screen{
		Width:  n.app.Width(),
		Height: n.app.Height(),
		Scale:  n.app.Scale(),
	}

	// find the first active camera
	camera, cameraExists := object.Query[camera.T]().First(scene)
	if !cameraExists {
		return nil, ErrNoCamera
	}

	// aquire next frame
	context, err := n.app.Aquire()
	if err != nil {
		return nil, ErrRecreate
	}

	// cache ticks
	n.app.Meshes().Tick(context.Index)
	n.app.Textures().Tick(context.Index)

	// create render arguments
	args := render.Args{
		Time:       time,
		Delta:      delta,
		Context:    context,
		Viewport:   screen,
		Near:       camera.Near(),
		Far:        camera.Far(),
		Fov:        camera.Fov(),
		Projection: camera.Projection(),
		View:       camera.View(),
		ViewInv:    camera.ViewInv(),
		VP:         camera.ViewProj(),
		VPInv:      camera.ViewProjInv(),
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

	worker := n.app.Worker(context.Index)
	worker.Submit(command.SubmitInfo{
		Marker: n.Name(),
		Wait:   waits,
		Signal: n.signals(context.Index),
	})

	return &args, nil
}
