package gui

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type Manager interface {
	object.Component

	DrawUI(render.Args, object.T)
}

type manager struct {
	object.Component

	gui      widget.T
	renderer node.Renderer
	scale    float32
}

func New(app node.RenderFunc) Manager {
	return &manager{
		Component: object.NewComponent(),
		renderer:  node.NewRenderer(app),
		scale:     1,
	}
}

func (m *manager) Name() string { return "GUIManager" }

func (m *manager) DrawUI(args render.Args, scene object.T) {
	viewport := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	m.scale = args.Viewport.Scale
	m.gui = m.renderer.Render(viewport)
	m.gui.Draw(args)
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

	if handler, ok := m.gui.(mouse.Handler); ok {
		handler.MouseEvent(ev)
		if ev.Handled() {
			e.Consume()
		}
	}
}
