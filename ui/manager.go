package ui

import (
	"github.com/johanhenriksson/goworld/render"

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

func NewManager(width, height float32) *Manager {
	m := &Manager{
		Width:    width,
		Height:   height,
		Viewport: mgl.Ortho(0, width, 0, height, 1000, -1000),
		Children: []render.Drawable{},
	}
	return m
}

func (m *Manager) Append(child render.Drawable) {
	m.Children = append(m.Children, child)
}

func (m *Manager) Draw() {
	/* create draw event args */
	p := m.Viewport
	v := mgl.Ident4()
	vp := p

	args := render.DrawArgs{
		Projection: p,
		View:       v,
		VP:         vp,
		MVP:        vp,
		Transform:  mgl.Ident4(),
	}

	render.ScreenBuffer.Bind()
	for _, el := range m.Children {
		el.Draw(args)
	}
}
