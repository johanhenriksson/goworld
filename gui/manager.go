package gui

import (
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Manager interface {
	object.Component

	DrawUI(widget.DrawArgs, object.T)
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

func (m *manager) DrawUI(args widget.DrawArgs, scene object.T) {
	viewport := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	m.scale = args.Viewport.Scale
	m.gui = m.renderer.Render(viewport)
	m.gui.Draw(widget.DrawArgs{
		Commands:  args.Commands,
		Meshes:    args.Meshes,
		Textures:  args.Textures,
		Transform: mat4.Ident(),
		ViewProj:  args.ViewProj,
		Viewport:  args.Viewport,
	})
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
		} else {
			// unset keyboard focus
			if e.Action() == mouse.Press || e.Action() == mouse.Release {
				keys.Focus(nil)
			}
		}
	}
}
