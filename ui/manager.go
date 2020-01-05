package ui

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

/** Main UI manager. Handles routing of events and drawing the UI. */
type Manager struct {
	/** Projection matrix - orthographic */
	Viewport mgl.Mat4
	Width    float32
	Height   float32

	Children []render.Drawable
}

func NewManager(app *engine.Application) *Manager {

	// ui manager
	//uimgr := ui.NewManager(float32(width)*scale, float32(height)*scale)

	width := float32(app.Window.Width) * app.Window.Scale()
	height := float32(app.Window.Height) * app.Window.Scale()

	m := &Manager{
		Width:    width,
		Height:   height,
		Viewport: mgl.Ortho(0, width, height, 0, 1000, -1000),
		Children: []render.Drawable{},
	}
	return m
}

func (m *Manager) Attach(child render.Drawable) {
	m.Children = append(m.Children, child)
}

func (m *Manager) DrawPass(scene *engine.Scene) {
	/* create draw event args */

	scaling := mgl.Scale3D(1, 1, 1)
	translation := mgl.Translate3D(0, 0, 0)
	root := translation.Mul4(scaling)

	p := m.Viewport // translation.Mul4(m.Viewport)

	v := mgl.Ident4()
	vp := p

	args := render.DrawArgs{
		Projection: p,
		View:       v,
		VP:         vp,
		MVP:        vp,
		Transform:  root, // mgl.Ident4(),
	}

	gl.Disable(gl.CULL_FACE)

	render.ScreenBuffer.Bind()
	for _, el := range m.Children {
		el.Draw(args)
	}
}

func (m *Manager) HandleMouse(x, y float32, button engine.MouseButton) {

}
