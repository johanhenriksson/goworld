package gui

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"

	"github.com/kjk/flex"
)

type Manager interface {
	object.Component

	DrawUI(widget.DrawArgs, object.T)
}

type manager struct {
	object.Component

	scale  float32
	render node.RenderFunc
	tree   node.T
	gui    widget.T
}

func New(render node.RenderFunc) Manager {
	return &manager{
		Component: object.NewComponent(),
		scale:     1,
		render:    render,
	}
}

func (m *manager) Name() string { return "UIManager" }

func (m *manager) DrawUI(args widget.DrawArgs, scene object.T) {
	// render root tree
	root := m.render()

	// populate with fragments
	fragments := query.New[Fragment]().Collect(scene)
	for {
		changed := false
		for idx, fragment := range fragments {
			if fragment == nil {
				// nil fragments have already been processed
				continue
			}

			parent := findNodeWithKey(root, fragment.Slot())
			if parent != nil {
				parent.SetChildren(append(parent.Children(), fragment.Render()))
				fragments[idx] = nil // set item to nil to mark it as completed
				changed = true
			}
		}
		if !changed {
			// iterate until nothing changes, or the fragment map is empty
			break
		}
	}

	// reconcile & hydrate tree
	m.tree = node.Reconcile(m.tree, root)
	m.gui = m.tree.Hydrate()

	// update flexbox layout
	viewport := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	m.scale = args.Viewport.Scale

	flexRoot := m.gui.Flex()
	flex.CalculateLayout(flexRoot, viewport.X, viewport.Y, flex.DirectionLTR)

	// draw
	m.gui.Draw(widget.DrawArgs{
		Commands:  args.Commands,
		Meshes:    args.Meshes,
		Textures:  args.Textures,
		Transform: mat4.Ident(),
		ViewProj:  args.ViewProj,
		Viewport:  args.Viewport,
	})
}

func findNodeWithKey(root node.T, key string) node.T {
	if root.Key() == key {
		return root
	}
	for _, child := range root.Children() {
		if hit := findNodeWithKey(child, key); hit != nil {
			return hit
		}
	}
	return nil
}

func (m *manager) MouseEvent(e mouse.Event) {
	// if the cursor is locked, we consider the game to have focus

	if e.Locked() {
		return
	}

	if m.gui == nil {
		// no rendered gui
		return
	}

	// apply UI scaling to cursor position
	offset := e.Position().Sub(e.Position().Scaled(1 / m.scale))
	ev := e.Project(offset)

	// pass on to gui fragments
	if handler, ok := m.gui.(mouse.Handler); ok {
		handler.MouseEvent(ev)
		if ev.Handled() {
			e.Consume()
			return
		}
	}

	// event has not been handled
	// unset keyboard focus
	if e.Action() == mouse.Press || e.Action() == mouse.Release {
		keys.Focus(nil)
	}
}
