package gui

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Manager interface {
	object.T
	Draw(render.Args, object.T)
}

type manager struct {
	object.T
	scale float32

	renderer node.Renderer
	gui      widget.T
}

func New(app node.RenderFunc) Manager {
	root := func() node.T {
		return rect.New("GUI", rect.Props{
			Children: []node.T{app()},
		})
	}

	mgr := &manager{
		T:        object.New("GUI Manager"),
		renderer: node.NewRenderer(root),
		scale:    1,
	}

	return mgr
}

func (m *manager) Draw(args render.Args, scene object.T) {
	width, height := float32(args.Viewport.FrameWidth), float32(args.Viewport.FrameHeight)
	m.scale = width / float32(args.Viewport.Width)

	// render GUI elements
	viewport := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	m.gui = m.renderer.Render(viewport)

	proj := mat4.Orthographic(0, width, height, 0, 1000, -1000)
	view := mat4.Scale(vec3.New(m.scale, m.scale, 1)) // todo: ui scaling
	model := mat4.Translate(vec3.New(0, 0, -100))
	vp := proj.Mul(&view)
	mvp := vp.Mul(&model)

	gl.DepthFunc(gl.LEQUAL)

	uiArgs := render.Args{
		Projection: proj,
		View:       view,
		VP:         vp,
		MVP:        mvp,
		Transform:  model,
		Viewport:   args.Viewport,
	}

	m.gui.Draw(uiArgs)
}

func (m *manager) MouseEvent(e mouse.Event) {
	// if the cursor is locked, we consider the game to have focus
	if e.Locked() {
		return
	}

	// scale down to low dpi.
	ev := e.Project(e.Position().Scaled(1 / m.scale))
	if handler, ok := m.gui.(mouse.Handler); ok {
		handler.MouseEvent(ev)
	}

	// check if the event actually hit something
	for _, child := range m.gui.Children() {
		if _, ok := child.(mouse.Handler); ok {
			ev := ev.Project(child.Position())
			target := ev.Position()
			size := child.Size()
			if target.X < 0 || target.X > size.X || target.Y < 0 || target.Y > size.Y {
				// outside
				continue
			}

			// we hit something
			e.Consume()
			break
		}
	}
}
