package graph

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	Node(pass NodePass) Node
	Connect()
	Draw(worker command.Worker, args render.Args, scene object.T)
	Destroy()
}

type graph struct {
	device device.T
	pre    Node
	post   Node
	nodes  []Node
	todo   map[Node]bool
}

func New(dev device.T) T {
	return &graph{
		device: dev,
		nodes:  make([]Node, 0, 16),
		todo:   make(map[Node]bool, 16),
		pre:    newPreNode(dev),
		post:   newPostNode(dev),
	}
}

func (g *graph) Node(pass NodePass) Node {
	nd := newNode(g.device, pass.Name(), pass)
	g.nodes = append(g.nodes, nd)
	return nd
}

func (g *graph) Connect() {
	for _, node := range g.nodes {
		if len(node.Requires()) == 0 {
			node.After(g.pre, vk.PipelineStageTopOfPipeBit)
		}
	}
	for _, node := range g.nodes {
		if len(node.Dependants()) == 0 {
			g.post.After(node, vk.PipelineStageTopOfPipeBit)
		}
	}
}

func (g *graph) Draw(worker command.Worker, args render.Args, scene object.T) {
	// put all nodes in a todo list
	// for each node in todo list
	//   if all Before nodes are not in todo list
	//     record node
	//     remove node from todo list
	for _, n := range g.nodes {
		g.todo[n] = true
	}

	ready := func(n Node) bool {
		for _, req := range n.Requires() {
			if g.todo[req] {
				return false
			}
		}
		return true
	}

	g.pre.Draw(worker, args, scene)
	for len(g.todo) > 0 {
		progress := false
		for node := range g.todo {
			// check if ready
			if ready(node) {
				node.Draw(worker, args, scene)
				delete(g.todo, node)
				progress = true
				break
			}
		}
		if !progress {
			// dependency error
			panic("unable to make progress in render graph")
		}
	}
	g.post.Draw(worker, args, scene)
}

func (g *graph) Destroy() {
	for _, node := range g.nodes {
		node.Destroy()
	}
	g.device = nil
	g.nodes = nil
}
