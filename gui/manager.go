package gui

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/gui/hooks"
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Manager interface {
	object.T
	Draw(render.Args, scene.T)
}

type manager struct {
	object.T
	scale float32

	tree node.T
	gui  widget.T
	root func() node.T
}

func New() Manager {
	root := func() node.T {
		f := TestUI()
		// f.Move(vec2.New(500, 300))
		scene := rect.New("GUI", &rect.Props{
			Layout:   layout.Absolute{},
			Children: []node.T{f},
		})
		// scene.Resize(vec2.New(1600, 1200))
		return scene
	}

	mgr := &manager{
		T:     object.New("GUI Manager"),
		root:  root,
		scale: 1,
	}

	return mgr
}

func (m *manager) Draw(args render.Args, scene scene.T) {
	width, height := float32(args.Viewport.FrameWidth), float32(args.Viewport.FrameHeight)
	m.scale = width / float32(args.Viewport.Width)

	// todo: resize if changed
	// perhaps the root component always accepts screen size etc

	m.tree = node.Reconcile(m.tree, m.root())
	m.gui = m.tree.Hydrate()
	m.gui.Resize(vec2.NewI(args.Viewport.Width, args.Viewport.Height))

	proj := mat4.Orthographic(0, width, height, 0, 1000, -1000)
	view := mat4.Scale(vec3.New(m.scale, m.scale, 1)) // todo: ui scaling
	vp := proj.Mul(&view)

	gl.DepthFunc(gl.LEQUAL)

	uiArgs := render.Args{
		Projection: proj,
		View:       view,
		VP:         vp,
		MVP:        vp,
		Transform:  mat4.Ident(),
		Viewport:   args.Viewport,
	}

	m.gui.Draw(uiArgs)

	hooks.SetScene(scene)
}

func (m *manager) MouseEvent(e mouse.Event) {
	// if the cursor is locked, we consider the game to have focus
	if e.Locked() {
		return
	}

	// scale down to low dpi.
	ev := e.Project(e.Position().Scaled(1 / m.scale))

	hit := false
	for _, frame := range m.gui.Children() {
		if handler, ok := frame.(mouse.Handler); ok {
			fev := ev.Project(frame.Position())
			target := fev.Position()
			size := frame.Size()
			if target.X < 0 || target.X > size.X || target.Y < 0 || target.Y > size.Y {
				// outside
				continue
			}

			hit = true

			handler.MouseEvent(fev)
			if fev.Handled() {
				e.Consume()
				break
			}
		}
	}

	// consume the event if it hits any UI element
	if hit {
		e.Consume()
	}
}
