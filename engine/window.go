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

type UpdateCallback func(float32)
type RenderCallback func(*Window, float32)

type Window struct {
	Wnd           *glfw.Window
	Width         int
	Height        int
	updateCb      UpdateCallback
	renderCb      RenderCallback
	maxFrameTime  float64
	lastFrameTime float64
}

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

func (wnd *Window) Closed() bool {
	return wnd.Wnd.ShouldClose()
}

func (wnd *Window) SetMaxFps(fps int) {
	wnd.maxFrameTime = 1.0 / float64(fps)
}

func (wnd *Window) LockCursor() {
	wnd.Wnd.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func (wnd *Window) ReleaseCursor() {
	wnd.Wnd.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}

func (wnd *Window) SetRenderCallback(cb RenderCallback) {
	wnd.renderCb = cb
}

func (wnd *Window) SetUpdateCallback(cb UpdateCallback) {
	wnd.updateCb = cb
}

func (wnd *Window) Loop() {
	for !wnd.Closed() {
		t := glfw.GetTime()
		dt := float32(t - wnd.lastFrameTime)
		wnd.lastFrameTime = t

		if MouseDown(MouseButton1) {
			wnd.LockCursor()
		} else {
			wnd.ReleaseCursor()
		}

		if wnd.updateCb != nil {
			wnd.updateCb(dt)
		}

		if wnd.renderCb != nil {
			wnd.renderCb(wnd, dt)
		}

		UpdateMouse(dt)

		wnd.EndFrame()
		if wnd.maxFrameTime > 0 {
			elapsed := glfw.GetTime() - t
			dur := wnd.maxFrameTime - elapsed
			time.Sleep(time.Duration(dur) * time.Second)
		}
	}
}

func (wnd *Window) EndFrame() {
	wnd.Wnd.SwapBuffers()
	glfw.PollEvents()
}

func (wnd *Window) Terminate() {
	glfw.Terminate()
}

func (wnd *Window) GetBufferSize() (int, int) {
	return wnd.Wnd.GetFramebufferSize()
}

func (wnd *Window) Scale() float32 {
	fw, _ := wnd.Wnd.GetFramebufferSize()
	return float32(fw) / float32(wnd.Width)
}
