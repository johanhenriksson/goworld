package graph

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/johanhenriksson/goworld/render/sync"
)

type postNode struct {
	*node
	target engine.Target
}

func newPostNode(app engine.App, target engine.Target) *postNode {
	return &postNode{
		node:   newNode(app, "Post", nil),
		target: target,
	}
}

func (n *postNode) Present(worker command.Worker, context *swapchain.Context) {
	// submit a dummy pass that waits for all previous passes to complete, then signals the render complete semaphore
	worker.Submit(command.SubmitInfo{
		Marker:   n.Name(),
		Commands: command.Empty,
		Wait:     n.waits(context.Index),
		Signal:   []*sync.Semaphore{context.RenderComplete},
		Callback: func() {
			// release the context, allowing it to be reused
			context.Release()

			// todo: collect actual gpu timings here
			// todo: this is a reasonable place to release per-frame resources
		},
	})

	// present
	n.target.Present(worker, context)
}
