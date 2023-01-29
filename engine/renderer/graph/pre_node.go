package graph

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"

	vk "github.com/vulkan-go/vulkan"
)

type preNode struct {
	*node
}

func newPreNode(dev device.T) Node {
	return &preNode{
		node: newNode(dev, nil),
	}
}

func (n *preNode) Name() string {
	return "Pre"
}

func (n *preNode) Draw(worker command.Worker, args render.Args, scene object.T) {
	var waits []command.Wait
	if args.Context.ImageAvailable != nil {
		waits = []command.Wait{
			{
				Semaphore: args.Context.ImageAvailable,
				Mask:      vk.PipelineStageColorAttachmentOutputBit,
			},
		}
	}
	worker.Submit(command.SubmitInfo{
		Marker: n.Name(),
		Wait:   waits,
		Signal: n.signals,
	})
}
