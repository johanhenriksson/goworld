package engine

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	/* GLFW event handling must run on the main OS thread */
	runtime.LockOSThread()
}

// UpdateCallback defines the window update callback function
type UpdateCallback func(float32)

// RenderCallback defines the window render callback function
type RenderCallback func(*Window, float32)

// Window represents the main engine window
type Window struct {
	Wnd           *glfw.Window
	Width         int
	Height        int
	updateCb      UpdateCallback
	renderCb      RenderCallback
	maxFrameTime  float64
	lastFrameTime float64
}

// CreateWindow creates the main engine window, and the OpenGL
func CreateWindow(title string, width int, height int) *Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}

	/* GLFW Window settings */
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(10)

	/* Initialize OpenGL */
	if err := gl.Init(); err != nil {
		panic(err)
	}

	w := &Window{
		Width:         width,
		Height:        height,
		Wnd:           window,
		maxFrameTime:  0.0,
		lastFrameTime: glfw.GetTime(),
	}
	w.SetMaxFps(60)

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetKeyCallback(KeyCallback)
	window.SetMouseButtonCallback(MouseButtonCallback)
	window.SetCursorPosCallback(func(wnd *glfw.Window, x, y float64) {
		MouseMoveCallback(wnd, x, y, w.Scale())
	})

	return w
}

// Closed returns true if the window will close.
func (wnd *Window) Closed() bool {
	return wnd.Wnd.ShouldClose()
}

// SetMaxFps adjusts the maximum frame rate.
func (wnd *Window) SetMaxFps(fps int) {
	wnd.maxFrameTime = 1.0 / float64(fps)
}

// LockCursor locks cursor movement.
func (wnd *Window) LockCursor() {
	wnd.Wnd.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

// ReleaseCursor unlocks cursor movement.
func (wnd *Window) ReleaseCursor() {
	wnd.Wnd.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}

// SetRenderCallback sets the engine render callback.
func (wnd *Window) SetRenderCallback(cb RenderCallback) {
	wnd.renderCb = cb
}

// SetUpdateCallback sets the engine update callback.
func (wnd *Window) SetUpdateCallback(cb UpdateCallback) {
	wnd.updateCb = cb
}

// Loop runs the main engine loop.
func (wnd *Window) Loop() {
	for !wnd.Closed() {
		// calculate frame delta time
		t := glfw.GetTime()
		dt := float32(t - wnd.lastFrameTime)
		wnd.lastFrameTime = t

		// todo: this is part of the FPS camera.
		// move it somewhere more appropriate
		if MouseDown(MouseButton1) {
			wnd.LockCursor()
		} else {
			wnd.ReleaseCursor()
		}

		// update scene
		if wnd.updateCb != nil {
			wnd.updateCb(dt)
		}

		// render scene
		if wnd.renderCb != nil {
			wnd.renderCb(wnd, dt)
		}

		// end scene
		wnd.Wnd.SwapBuffers()
		glfw.PollEvents()

		UpdateMouse(dt)

		// wait a bit if we're running faster than the maximum fps
		if wnd.maxFrameTime > 0 {
			elapsed := glfw.GetTime() - t
			dur := wnd.maxFrameTime - elapsed
			time.Sleep(time.Duration(dur) * time.Second)
		}
	}
}

// Terminate kills the GLFW window.
func (wnd *Window) Terminate() {
	glfw.Terminate()
}

// GetBufferSize returns the window framebuffer size.
func (wnd *Window) GetBufferSize() (int, int) {
	return wnd.Wnd.GetFramebufferSize()
}

// Scale returns the window DPI scale relative to the framebuffer.
func (wnd *Window) Scale() float32 {
	fw, _ := wnd.Wnd.GetFramebufferSize()
	return float32(fw) / float32(wnd.Width)
}
