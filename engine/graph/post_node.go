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
	var signal []sync.Semaphore
	if context.RenderComplete != nil {
		signal = []sync.Semaphore{context.RenderComplete}
	}

	// submit a dummy pass that waits for all previous passes to complete, then signals the render complete semaphore
	worker.Submit(command.SubmitInfo{
		Marker:   n.Name(),
		Commands: command.Empty,
		Wait:     n.waits(context.Index),
		Signal:   signal,
	})

	// present
	worker.Invoke(func() {
		n.target.Present(context)
	})

	// flush ensures all commands are submitted before we start rendering the next frame. otherwise, frame submissions may overlap.
	// todo: perhaps its possible to do this at a later stage? e.g. we could run update loop etc while waiting
	// note: this is only required if we use multiple/per-frame workers
	// worker.Flush()
}
