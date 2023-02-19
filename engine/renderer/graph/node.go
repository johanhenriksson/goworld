package graph

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type NodePass interface {
	Name() string
	Record(command.Recorder, render.Args, object.T)
	Destroy()
}

type Node interface {
	After(nd Node, mask core1_0.PipelineStageFlags)
	Before(nd Node, mask core1_0.PipelineStageFlags, signal []sync.Semaphore)
	Requires() []Node
	Dependants() []Node

	Name() string
	Draw(command.Worker, render.Args, object.T)
	Detach(Node)
	Destroy()
}

type node struct {
	name       string
	app        vulkan.App
	pass       NodePass
	after      map[string]edge
	before     map[string]edge
	requires   []Node
	dependants []Node
}

type edge struct {
	node   Node
	mask   core1_0.PipelineStageFlags
	signal []sync.Semaphore
}

func newNode(app vulkan.App, name string, pass NodePass) *node {
	return &node{
		app:        app,
		name:       name,
		pass:       pass,
		after:      make(map[string]edge, 4),
		before:     make(map[string]edge, 4),
		requires:   make([]Node, 0, 4),
		dependants: make([]Node, 0, 4),
	}
}

func (n *node) Requires() []Node   { return n.requires }
func (n *node) Dependants() []Node { return n.dependants }

func (n *node) After(nd Node, mask core1_0.PipelineStageFlags) {
	if _, exists := n.after[nd.Name()]; exists {
		return
	}
	signal := sync.NewSemaphoreArray(n.app.Device(), fmt.Sprintf("%s->%s", nd.Name(), n.name), 3)
	n.after[nd.Name()] = edge{
		node:   nd,
		mask:   mask,
		signal: signal,
	}
	nd.Before(n, mask, signal)

	n.refresh()
}

func (n *node) Before(nd Node, mask core1_0.PipelineStageFlags, signal []sync.Semaphore) {
	if _, exists := n.before[nd.Name()]; exists {
		return
	}
	n.before[nd.Name()] = edge{
		node:   nd,
		mask:   mask,
		signal: signal,
	}
	nd.After(n, mask)

	n.refresh()
}

func (n *node) refresh() {
	// recompute signals
	n.dependants = make([]Node, 0, len(n.after))
	for _, edge := range n.before {
		n.dependants = append(n.dependants, edge.node)
	}

	// recompute waits
	n.requires = make([]Node, 0, len(n.after))
	for _, edge := range n.after {
		if edge.signal == nil {
			// skip nil signals
			continue
		}
		n.requires = append(n.requires, edge.node)
	}
}

func (n *node) Detach(nd Node) {
	if _, exists := n.before[nd.Name()]; exists {
		delete(n.before, nd.Name())
		nd.Detach(n)
	}
	if edge, exists := n.after[nd.Name()]; exists {
		delete(n.after, nd.Name())
		nd.Detach(n)
		// free semaphores
		for _, signal := range edge.signal {
			signal.Destroy()
		}
	}
	n.refresh()
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Destroy() {
	for _, edge := range n.before {
		before := edge.node
		before.Detach(n)
		for _, s := range edge.signal {
			s.Destroy()
		}
	}
	for _, edge := range n.after {
		after := edge.node
		after.Detach(n)
	}
	if n.pass != nil {
		n.pass.Destroy()
		n.pass = nil
	}
	n.before = nil
	n.after = nil
}

func (n *node) waits(index int) []command.Wait {
	waits := make([]command.Wait, 0, len(n.after))
	for _, after := range n.after {
		if after.signal == nil {
			// skip nil signals
			continue
		}
		waits = append(waits, command.Wait{
			Semaphore: after.signal[index],
			Mask:      after.mask,
		})
	}
	return waits
}

func (n *node) signals(index int) []sync.Semaphore {
	signals := make([]sync.Semaphore, 0, len(n.before))
	for _, edge := range n.before {
		signals = append(signals, edge.signal[index])
	}
	return signals
}

func (n *node) Draw(worker command.Worker, args render.Args, scene object.T) {
	if n.pass == nil {
		return
	}
	cmds := command.NewRecorder()
	n.pass.Record(cmds, args, scene)
	worker.Queue(cmds.Apply)

	worker.Submit(command.SubmitInfo{
		Marker: n.pass.Name(),
		Wait:   n.waits(args.Context.Index),
		Signal: n.signals(args.Context.Index),
	})
}
