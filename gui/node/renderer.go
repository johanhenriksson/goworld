package node

import (
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"

	"github.com/kjk/flex"
)

type RenderFunc func() T

type Renderer interface {
	Render(viewport vec2.T) widget.T
}

type renderer struct {
	key     string
	root    RenderFunc
	tree    T
	display widget.T
}

func NewRenderer(key string, app RenderFunc) Renderer {
	return &renderer{
		root: app,
	}
}

func (r *renderer) Render(viewport vec2.T) widget.T {
	r.tree = Reconcile(r.tree, r.root())
	r.display = r.tree.Hydrate(r.key)

	root := r.display.Flex()
	flex.CalculateLayout(root, viewport.X, viewport.Y, flex.DirectionLTR)

	return r.display
}
