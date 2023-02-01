package graph

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type preNode struct {
	*node
}

func newPreNode(dev device.T) Node {
	return &preNode{
		node: newNode(dev, "Pre", nil),
	}
}

func (n *preNode) Draw(worker command.Worker, args render.Args, scene object.T) {
	var waits []command.Wait
	if args.Context.ImageAvailable != nil {
		waits = []command.Wait{
			{
				Semaphore: args.Context.ImageAvailable,
				Mask:      core1_0.PipelineStageColorAttachmentOutput,
			},
		}
	}
	args.Context.InFlight.Lock()
	worker.Submit(command.SubmitInfo{
		Marker: n.Name(),
		Wait:   waits,
		Signal: n.signals,
	})
}
