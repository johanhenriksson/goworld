package graph

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type postNode struct {
	*node
	target vulkan.Target
}

func newPostNode(app vulkan.App, target vulkan.Target) Node {
	return &postNode{
		node:   newNode(app, "Post", nil),
		target: target,
	}
}

func (n *postNode) Draw(worker command.Worker, args render.Args, scene object.Component) {
	var signal []sync.Semaphore
	if args.Context.RenderComplete != nil {
		signal = []sync.Semaphore{args.Context.RenderComplete}
	}

	worker.Submit(command.SubmitInfo{
		Marker: n.Name(),
		Wait:   n.waits(args.Context.Index),
		Signal: signal,
		Callback: func() {
			args.Context.Release()
		},
	})

	// present
	n.target.Present(worker, args.Context)

	// flush ensures all commands are submitted before we start rendering the next frame. otherwise, frame submissions may overlap.
	// todo: perhaps its possible to do this at a later stage? e.g. we could run update loop etc while waiting
	// note: this is only required if we use multiple/per-frame workers
	// worker.Flush()
}
