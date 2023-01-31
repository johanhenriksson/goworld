package graph

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"
)

type postNode struct {
	*node
}

func newPostNode(dev device.T) Node {
	return &postNode{
		node: newNode(dev, "Post", nil),
	}
}

func (n *postNode) Draw(worker command.Worker, args render.Args, scene object.T) {
	var signal []sync.Semaphore
	if args.Context.RenderComplete != nil {
		signal = []sync.Semaphore{args.Context.RenderComplete}
	}

	worker.Submit(command.SubmitInfo{
		Marker: n.Name(),
		Wait:   n.waits,
		Signal: signal,
		Then: func() {
			args.Context.InFlight.Unlock()
		},
	})
}
