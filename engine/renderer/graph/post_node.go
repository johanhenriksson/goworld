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
}

func newPostNode(app vulkan.App) Node {
	return &postNode{
		node: newNode(app, "Post", nil),
	}
}

func (n *postNode) Draw(worker command.Worker, args render.Args, scene object.T) {
	var signal []sync.Semaphore
	if args.Context.RenderComplete != nil {
		signal = []sync.Semaphore{args.Context.RenderComplete}
	}

	barrier := make(chan struct{})
	worker.OnComplete(func() {
		<-barrier
		args.Context.Release()
	})
	worker.Submit(command.SubmitInfo{
		Marker: n.Name(),
		Wait:   n.waits(args.Context.Index),
		Signal: signal,
	})

	// submit cache work
	n.app.Meshes().Submit()
	n.app.Textures().Submit()

	// present
	n.app.Present(worker, args.Context)
	barrier <- struct{}{}
}
