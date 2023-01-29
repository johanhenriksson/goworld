package graph

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/sync"

	vk "github.com/vulkan-go/vulkan"
)

type NodePass interface {
	Name() string
	Record(command.Recorder, render.Args, object.T)
	Destroy()
}

type Node interface {
	After(nd Node, mask vk.PipelineStageFlagBits)
	Before(nd Node, mask vk.PipelineStageFlagBits, signal sync.Semaphore)
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
	before     map[Node]sync.Semaphore
	requires   []Node
	dependants []Node
	waits      []command.Wait
	signals    []sync.Semaphore
}

type edge struct {
	mask   vk.PipelineStageFlagBits
	signal sync.Semaphore
}

func newNode(dev device.T, name string, pass NodePass) *node {
	return &node{
		device:     dev,
		name:       name,
		pass:       pass,
		after:      make(map[Node]edge, 4),
		before:     make(map[Node]sync.Semaphore, 4),
		waits:      make([]command.Wait, 0, 4),
		signals:    make([]sync.Semaphore, 0, 4),
		requires:   make([]Node, 0, 4),
		dependants: make([]Node, 0, 4),
	}
}

func (n *node) Requires() []Node   { return n.requires }
func (n *node) Dependants() []Node { return n.dependants }

func (n *node) After(nd Node, mask vk.PipelineStageFlagBits) {
	if _, exists := n.after[nd]; exists {
		return
	}
	signal := sync.NewSemaphore(n.device)
	n.after[nd] = edge{
		mask:   mask,
		signal: signal,
	}
	nd.Before(n, mask, signal)

	n.refresh()
}

func (n *node) Before(nd Node, mask vk.PipelineStageFlagBits, signal sync.Semaphore) {
	if _, exists := n.before[nd]; exists {
		return
	}
	n.before[nd] = signal
	nd.After(n, vk.PipelineStageTopOfPipeBit)

	n.refresh()
}

func (n *node) refresh() {
	// recompute signals
	n.signals = make([]sync.Semaphore, 0, len(n.before))
	n.dependants = make([]Node, 0, len(n.after))
	for node, before := range n.before {
		n.signals = append(n.signals, before)
		n.dependants = append(n.dependants, node)
	}

	// recompute waits
	n.waits = make([]command.Wait, 0, len(n.after))
	n.requires = make([]Node, 0, len(n.after))
	for node, after := range n.after {
		if after.signal == nil {
			// skip nil signals
			continue
		}
		n.waits = append(n.waits, command.Wait{
			Semaphore: after.signal,
			Mask:      after.mask,
		})
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
		signal.Destroy()
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
	n.signals = nil
	n.waits = nil
}

func (n *node) Draw(worker command.Worker, args render.Args, scene object.T) {
	cmds := command.NewRecorder()
	n.pass.Record(cmds, args, scene)

	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Marker: n.pass.Name(),
		Wait:   n.waits,
		Signal: n.signals,
	})
}
