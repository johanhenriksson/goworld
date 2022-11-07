package gui

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/util"

	"github.com/kjk/flex"
)

type Manager interface {
	object.Component

	DrawUI(widget.DrawArgs, object.T)
}

type manager struct {
	object.Component

	scale     float32
	fragments []widget.T
	root      node.RenderFunc

	tree node.T
	gui  widget.T
}

func New(root node.RenderFunc) Manager {
	return &manager{
		Component: object.NewComponent(),
		scale:     1,
		root:      root,
	}
}

func (m *manager) Name() string { return "UIManager" }

func (m *manager) DrawUI(args widget.DrawArgs, scene object.T) {
	viewport := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	m.scale = args.Viewport.Scale

	fragments := query.New[Fragment]().Collect(scene)

	fragmap := make(map[string][]Fragment)
	for _, fragment := range fragments {
		fragmap[fragment.Slot()] = append(fragmap[fragment.Slot()], fragment)
	}

	root := m.root()
	roots := []node.T{root}

	for len(fragmap) > 0 {
		changed := false
		for slot, fragments := range fragmap {
			parent := findNodeWithKey(roots, slot)
			if parent != nil {
				rendered := util.Map(fragments, func(f Fragment) node.T { return f.Render() })
				parent.SetChildren(append(parent.Children(), rendered...))
				delete(fragmap, slot)
				changed = true
			}
			break
		}
		if !changed {
			break
		}
	}
	for slot, left := range fragmap {
		fmt.Printf("dangling slot %s: %d elements\n", slot, len(left))
	}

	m.tree = node.Reconcile(m.tree, root)
	m.gui = m.tree.Hydrate()

	flexRoot := m.gui.Flex()
	flex.CalculateLayout(flexRoot, viewport.X, viewport.Y, flex.DirectionLTR)

	m.gui.Draw(widget.DrawArgs{
		Commands:  args.Commands,
		Meshes:    args.Meshes,
		Textures:  args.Textures,
		Transform: mat4.Ident(),
		ViewProj:  args.ViewProj,
		Viewport:  args.Viewport,
	})
}

func findNodeWithKey(nodes []node.T, key string) node.T {
	for _, n := range nodes {
		if n.Key() == key {
			return n
		}
		children := n.Children()
		if child := findNodeWithKey(children, key); child != nil {
			return child
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
