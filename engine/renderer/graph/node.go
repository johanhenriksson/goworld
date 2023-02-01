package graph

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"

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
	device     device.T
	pass       NodePass
	after      map[Node]edge
	before     map[Node][]sync.Semaphore
	requires   []Node
	dependants []Node
}

type edge struct {
	mask   core1_0.PipelineStageFlags
	signal []sync.Semaphore
}

func newNode(dev device.T, name string, pass NodePass) *node {
	return &node{
		device:     dev,
		name:       name,
		pass:       pass,
		after:      make(map[Node]edge, 4),
		before:     make(map[Node][]sync.Semaphore, 4),
		requires:   make([]Node, 0, 4),
		dependants: make([]Node, 0, 4),
	}
}

func (n *node) Requires() []Node   { return n.requires }
func (n *node) Dependants() []Node { return n.dependants }

func (n *node) After(nd Node, mask core1_0.PipelineStageFlags) {
	if _, exists := n.after[nd]; exists {
		return
	}
	signal := sync.NewSemaphoreArray(n.device, fmt.Sprintf("%s->%s", nd.Name(), n.name), 3)
	n.after[nd] = edge{
		mask:   mask,
		signal: signal,
	}
	nd.Before(n, mask, signal)

	n.refresh()
}

func (n *node) Before(nd Node, mask core1_0.PipelineStageFlags, signal []sync.Semaphore) {
	if _, exists := n.before[nd]; exists {
		return
	}
	n.before[nd] = signal
	nd.After(n, core1_0.PipelineStageTopOfPipe)

	n.refresh()
}

func (n *node) refresh() {
	// recompute signals
	n.dependants = make([]Node, 0, len(n.after))
	for node := range n.before {
		n.dependants = append(n.dependants, node)
	}

	// recompute waits
	n.requires = make([]Node, 0, len(n.after))
	for node, after := range n.after {
		if after.signal == nil {
			// skip nil signals
			continue
		}
		n.requires = append(n.requires, node)
	}
}

func (n *node) Detach(nd Node) {
	delete(n.before, nd)
	delete(n.after, nd)
	n.refresh()
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Destroy() {
	for before, signal := range n.before {
		before.Detach(n)
		for _, s := range signal {
			s.Destroy()
		}
	}
	for after := range n.after {
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
	for _, before := range n.before {
		signals = append(signals, before[index])
	}
	return signals
}

func (n *node) Draw(worker command.Worker, args render.Args, scene object.T) {
	cmds := command.NewRecorder()
	n.pass.Record(cmds, args, scene)

	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Marker: n.pass.Name(),
		Wait:   n.waits(args.Context.Index),
		Signal: n.signals(args.Context.Index),
	})
}
