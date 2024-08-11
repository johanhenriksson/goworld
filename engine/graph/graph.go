package graph

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render/upload"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type GraphFunc func(*Graph, engine.Target) []Resource

type Resource interface {
	Destroy()
}

type Graph struct {
	app       engine.App
	target    engine.Target
	pre       *preNode
	post      *postNode
	nodes     []Node
	todo      map[Node]bool
	init      GraphFunc
	resources []Resource
}

func New(app engine.App, output engine.Target, init GraphFunc) *Graph {
	g := &Graph{
		app:    app,
		target: output,
		nodes:  make([]Node, 0, 16),
		todo:   make(map[Node]bool, 16),
		init:   init,
	}
	g.Recreate()
	return g
}

func (g *Graph) Recreate() {
	g.Destroy()
	g.app.Pool().Recreate()

	g.resources = g.init(g, g.target)

	g.pre = newPreNode(g.app, g.target)
	g.post = newPostNode(g.app, g.target)
	g.connect()
}

func (g *Graph) Node(pass draw.Pass) Node {
	nd := newNode(g.app, pass.Name(), pass)
	g.nodes = append(g.nodes, nd)
	return nd
}

func (g *Graph) connect() {
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

func (g *Graph) Draw(scene object.Object, time, delta float32) {
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
	args, context, err := g.pre.Prepare(scene, time, delta)
	if err != nil {
		log.Println("Render preparation error:", err)
		g.Recreate()
		return
	}

	// select a suitable worker for this frame
	worker := g.app.Worker()

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

	g.post.Present(worker, context)
}

func (g *Graph) Screengrab() *image.RGBA {
	idx := 0
	g.app.Device().WaitIdle()
	source := g.target.Surfaces()[idx]
	ss, err := upload.DownloadImage(g.app.Device(), g.app.Worker(), source)
	if err != nil {
		panic(err)
	}
	return ss
}

func (g *Graph) Screenshot() {
	img := g.Screengrab()
	filename := fmt.Sprintf("Screenshot-%s.png", time.Now().Format("2006-01-02_15-04-05"))
	if err := upload.SavePng(img, filename); err != nil {
		panic(err)
	}
	log.Println("saved screenshot", filename)
}

func (g *Graph) Destroy() {
	g.app.Flush()
	for _, resource := range g.resources {
		resource.Destroy()
	}
	g.resources = nil
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
