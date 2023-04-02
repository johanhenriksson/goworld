package gui

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"

	"github.com/kjk/flex"
)

type Manager interface {
	object.T

	DrawUI(args widget.DrawArgs, quads *widget.QuadBuffer)
}

type manager struct {
	object.T

	viewport render.Screen
	render   node.RenderFunc
	tree     node.T
	gui      widget.T
	updated  chan struct{}

	supressNextRelease bool
}

func New(renderNodes node.RenderFunc) Manager {
	return object.New(&manager{
		viewport: render.Screen{
			Scale: 1,
		},
		updated: make(chan struct{}),
		render:  renderNodes,
	})
}

func (m *manager) Name() string { return "UIManager" }

func (m *manager) Update(scene object.T, dt float32) {
	// render root tree
	root := m.render()

	// find fragments
	fragments := object.NewQuery[Fragment]().Collect(scene)

	// populate with fragments
	for {
		changed := false
		for idx, fragment := range fragments {
			if fragment == nil {
				// nil fragments have already been processed
				continue
			}

			target := findNodeWithKey(root, fragment.Slot())
			if target == nil {
				// target slot is not available (yet)
				continue
			}

			frag := fragment.Render()
			switch fragment.Position() {
			case FragmentLast:
				target.Append(frag)
			case FragmentFirst:
				target.Prepend(frag)
			default:
				panic("invalid fragment position")
			}

			fragments[idx] = nil // set item to nil to mark it as completed
			changed = true
		}
		if !changed {
			// iterate until nothing changes, or the fragment map is empty
			break
		}
	}

	m.tree = node.Reconcile(m.tree, root)

	// 	go m.update(scene, dt, root, fragments)
	// }

	// func (m *manager) update(scene object.T, dt float32, root node.T, fragments []Fragment) {

	key := object.Key("gui", m)
	m.gui = m.tree.Hydrate(key)

	// update flexbox layout
	viewport := vec2.NewI(m.viewport.Width, m.viewport.Height)

	flexRoot := m.gui.Flex()
	flex.CalculateLayout(flexRoot, viewport.X, viewport.Y, flex.DirectionLTR)

	// signal background update completed
	// m.updated <- struct{}{}
}

func (m *manager) DrawUI(args widget.DrawArgs, quads *widget.QuadBuffer) {
	m.viewport = args.Viewport
	// draw
	m.gui.Draw(args, quads)
}

func findNodeWithKey(root node.T, key string) node.T {
	if root == nil {
		return nil
	}
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
	offset := e.Position().Sub(e.Position().Scaled(1 / m.viewport.Scale))
	ev := e.Project(offset)

	// pass on to gui fragments
	if handler, ok := m.gui.(mouse.Handler); ok {
		handler.MouseEvent(ev)
		if ev.Handled() {
			e.Consume()
			if ev.Action() == mouse.Press {
				m.supressNextRelease = true
			}
		}
	}

	// event has not been handled
	// unset keyboard focus
	if !ev.Handled() && (e.Action() == mouse.Press || e.Action() == mouse.Release) {
		keys.Focus(nil)
	}

	// if the UI captured a press event, we should make sure not to pass
	// the matching release event
	if ev.Action() == mouse.Release && m.supressNextRelease {
		m.supressNextRelease = false
		e.Consume()
	}
}
