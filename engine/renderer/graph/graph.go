package graph

import (
	"log"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type NodeFunc func(T, pass.RenderTarget, pass.GeometryBuffer)

// The render graph is responsible for synchronization between
// different render nodes.
type T interface {
	Node(pass NodePass) Node
	Recreate()
	Draw(scene object.Object, time, delta float32)
	Destroy()
}

type graph struct {
	app   vulkan.App
	pre   PreNode
	post  Node
	nodes []Node
	todo  map[Node]bool
	init  NodeFunc

	target  pass.RenderTarget
	gbuffer pass.GeometryBuffer
}

func New(app vulkan.App, init NodeFunc) T {
	g := &graph{
		app:   app,
		nodes: make([]Node, 0, 16),
		todo:  make(map[Node]bool, 16),
		init:  init,
	}
	g.Recreate()
	return g
}

func (g *graph) Recreate() {
	var err error

	g.Destroy()
	g.app.Pool().Recreate()

	g.target, err = pass.NewRenderTarget(g.app.Device(),
		g.app.Width(), g.app.Height(), g.app.Frames(),
		image.FormatRGBA8Unorm, g.app.Device().GetDepthFormat())
	if err != nil {
		panic(err)
	}

	g.gbuffer, err = pass.NewGbuffer(g.app.Device(), g.app.Width(), g.app.Height(), g.app.Frames())
	if err != nil {
		panic(err)
	}

	g.init(g, g.target, g.gbuffer)
	g.pre = newPreNode(g.app)
	g.post = newPostNode(g.app)
	g.connect()
}

func (g *graph) Node(pass NodePass) Node {
	nd := newNode(g.app, pass.Name(), pass)
	g.nodes = append(g.nodes, nd)
	return nd
}

func (g *graph) connect() {
	// use bottom of pipe so that subsequent passes start as soon as possible
	for _, node := range g.nodes {
		if len(node.Requires()) == 0 {
			node.After(g.pre, core1_0.PipelineStageTopOfPipe)
		}
	}
	for _, node := range g.nodes {
		if len(node.Dependants()) == 0 {
			g.post.After(node, core1_0.PipelineStageTopOfPipe)
		}
	}
}

func (g *graph) Draw(scene object.Object, time, delta float32) {
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

	// prepare
	args, err := g.pre.Prepare(scene, time, delta)
	if err != nil {
		log.Println("Render preparation error:", err)
		g.Recreate()
		return
	}

	// select a suitable worker for this frame
	worker := g.app.Worker(args.Context.Index)

	for len(g.todo) > 0 {
		progress := false
		for node := range g.todo {
			// check if ready
			if ready(node) {
				node.Draw(worker, *args, scene)
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
	g.post.Draw(worker, *args, scene)
}

func (g *graph) Destroy() {
	g.app.Flush()
	if g.gbuffer != nil {
		g.gbuffer.Destroy()
		g.gbuffer = nil
	}
	if g.target != nil {
		g.target.Destroy()
		g.target = nil
	}
	if g.pre != nil {
		g.pre.Destroy()
		g.pre = nil
	}
	if g.post != nil {
		g.post.Destroy()
		g.post = nil
	}
	for _, node := range g.nodes {
		node.Destroy()
	}
	g.nodes = g.nodes[:0]
}
