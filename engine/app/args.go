package app

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/graph"
)

type Args struct {
	Title    string
	Width    int
	Height   int
	Renderer engine.RendererFunc
}

func (a *Args) Defaults() *Args {
	if a.Title == "" {
		a.Title = "goworld"
	}
	if a.Width == 0 {
		a.Width = 800
	}
	if a.Height == 0 {
		a.Height = 600
	}
	if a.Renderer == nil {
		a.Renderer = graph.Default
	}
	return a
}
