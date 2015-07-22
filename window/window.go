package window

import (
    "log"
    "time"
    "runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

    "github.com/johanhenriksson/goworld/engine"
)

/* GLFW event handling must run on the main OS thread */
func init() {
	runtime.LockOSThread()
}

type UpdateCallback func(float32)
type RenderCallback func(*Window, float32)

type Window struct {
    Wnd             *glfw.Window
    updateCb        UpdateCallback
    renderCb        RenderCallback
    maxFrameTime    float64
    lastFrameTime   float64
}

func Create(title string, width int, height int) *Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}

    /* GLFW Window settings */
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
    glfw.SwapInterval(1);

	/* Initialize OpenGL */
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

    window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
    window.SetKeyCallback(engine.KeyCallback)
    window.SetCursorPosCallback(engine.MouseMoveCallback)
    window.SetMouseButtonCallback(engine.MouseButtonCallback)

    w := &Window {
        Wnd:            window,
        maxFrameTime:   0.0,
        lastFrameTime:  glfw.GetTime(),
    }
    w.SetMaxFps(60)
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
        t  := glfw.GetTime()
        dt := float32(t - wnd.lastFrameTime)
        wnd.lastFrameTime = t

        engine.UpdateMouse(dt)
        if engine.MouseDown(engine.MouseButton1) {
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
