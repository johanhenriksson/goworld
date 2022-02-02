package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

// Manager is the main UI manager. Handles routing of events and drawing of the UI.
type Manager struct {
	Viewport mat4.T
	Width    float32
	Height   float32
	Focused  Component

	Children []Component
}

// NewManager creates a new UI manager.
func NewManager(width, height float32) *Manager {
	// grab UI dimensions from application window
	// width := float32(app.Window.Width)   // * app.Window.Scale
	// height := float32(app.Window.Height) // * app.Window.Scale

	m := &Manager{
		Width:    width,
		Height:   height,
		Viewport: mat4.Orthographic(0, width, height, 0, 1000, -1000),
		Children: []Component{},
	}

	// hook GLFW input event callbacks - this allows the UI to capture events
	// not very elegant, would be cool to do this in a cleaner way
	// app.Window.Wnd.SetKeyCallback(m.glfwKeyCallback)
	// app.Window.Wnd.SetMouseButtonCallback(m.glfwMouseButtonCallback)
	// app.Window.Wnd.SetCharCallback(m.glfwInputCallback)

	// watermark / fps text
	// m.Attach(NewWatermark(app.Window))

	return m
}

// Attach a child component
func (m *Manager) Attach(child Component) {
	m.Children = append(m.Children, child)
}

// DrawPass draws the UI
func (m *Manager) Draw(args render.Args, scene scene.T) {
	p := m.Viewport
	v := mat4.Ident() // unused by UI
	vp := p           // unused by UI

	args.Projection = p
	args.View = v
	args.VP = vp
	args.MVP = vp
	args.Transform = mat4.Ident()

	// ensure back face culling is disabled
	// since UI is scaled by Y-1, we only want back faces
	gl.Disable(gl.CULL_FACE)
	render.Blend(true)

	// clear depth buffer
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	render.BindScreenBuffer()
	render.SetViewport(0, 0, args.Viewport.FrameWidth, args.Viewport.FrameHeight)
	for _, el := range m.Children {
		el.Draw(args)
	}
}

// Focus the given component
func (m *Manager) Focus(target Component) {
	if m.Focused != nil {
		m.Focused.Blur()
	}
	m.Focused = target
	if target != nil {
		target.Focus()
	}
}

func (m *Manager) glfwMouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	// we're only interested in mouse down events at this time
	if action != glfw.Release {
		// supress if the event was handled by an element
		// if m.handleMouse(mouse.Position, mouse.Button(button)) {
		// 	return
		// }
	}
}

func (m *Manager) handleMouse(position vec2.T, button mouse.Button) bool {
	// reset focus
	m.Focus(nil)

	event := MouseEvent{
		UI:     m,
		Point:  position,
		Button: button,
	}
	for _, el := range m.Children {
		handled := el.HandleMouse(event)
		if handled {
			return true
		}
	}

	return false
}

func (m *Manager) glfwKeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// supress key events from engine while an element is focused
	if m.Focused != nil {
		ev := KeyEvent{
			UI:    m,
			Key:   keys.Code(key),
			Press: action == glfw.Press,
		}
		m.Focused.HandleKey(ev)
	}
}

func (m *Manager) glfwInputCallback(w *glfw.Window, char rune) {
	// pass to focused element
	if m.Focused != nil {
		m.Focused.HandleInput(char)
	}
}

func (m *Manager) Resize(width, height int) {}
