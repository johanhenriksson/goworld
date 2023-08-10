package graph

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type postNode struct {
	*node
	target vulkan.Target
}

func newPostNode(app vulkan.App, target vulkan.Target) *postNode {
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

	worker.Submit(command.SubmitInfo{
		Marker: n.Name(),
		Wait:   n.waits(context.Index),
		Signal: signal,
		Callback: func() {
			context.Release()
		},
	})

	// present
	n.target.Present(worker, context)

	// flush ensures all commands are submitted before we start rendering the next frame. otherwise, frame submissions may overlap.
	// todo: perhaps its possible to do this at a later stage? e.g. we could run update loop etc while waiting
	// note: this is only required if we use multiple/per-frame workers
	// worker.Flush()
}
