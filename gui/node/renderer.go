package node

import "github.com/johanhenriksson/goworld/gui/widget"

type RenderFunc func() T

type Renderer interface {
	Render() widget.T
}

type renderer struct {
	root    RenderFunc
	tree    T
	display widget.T
}

func NewRenderer(app RenderFunc) Renderer {
	return &renderer{
		root: app,
	}
}

func (r *renderer) Render() widget.T {
	r.tree = Reconcile(r.tree, r.root())
	r.display = r.tree.Hydrate()
	return r.display
}
