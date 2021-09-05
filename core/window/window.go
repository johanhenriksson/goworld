package window

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
)

func init() {
	// glfw event handling must run on the main OS thread
	runtime.LockOSThread()

	// init glfw
	if err := glfw.Init(); err != nil {
		panic(err)
	}
}

type T interface {
	Size() (int, int)
	BufferSize() (int, int)

	SwapBuffers()
	ShouldClose() bool
}

type Args struct {
	Title        string
	Width        int
	Height       int
	Vsync        bool
	InputHandler input.Handler
}

type window struct {
	*glfw.Window
	mouse mouse.MouseWrapper
}

func New(args Args) (T, error) {
	// window creation hints.
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 1)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	// create a new GLFW window
	wnd, err := glfw.CreateWindow(args.Width, args.Height, args.Title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create glfw window: %w", err)
	}

	width, height := wnd.GetSize()
	fwidth, fheight := wnd.GetFramebufferSize()
	scale := float32(fwidth) / float32(width)
	log.Printf("Created window. Size %dx%d, Buffer Size %dx%d. Scale = %.0f%%\n",
		width, height, fwidth, fheight, scale*100)

	window := &window{
		Window: wnd,
	}

	if args.Vsync {
		glfw.SwapInterval(1)
	}

	// activate OpenGL context
	wnd.MakeContextCurrent()

	// initialize OpenGL. context must be active first
	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	// ensure the frame buffer is properly cleared
	gl.ClearColor(0, 0, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	wnd.SwapBuffers()
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// attach default input handler, if provided
	if args.InputHandler != nil {
		window.SetInputHandler(args.InputHandler)
	}

	wnd.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	return window, nil
}

func (w *window) SwapBuffers() {
	w.Window.SwapBuffers()

	glfw.PollEvents()
}

func (w *window) Size() (int, int)       { return w.GetSize() }
func (w *window) BufferSize() (int, int) { return w.GetFramebufferSize() }

func (w *window) SetInputHandler(handler input.Handler) {
	// keyboard events
	w.SetKeyCallback(keys.KeyCallbackWrapper(handler))
	w.SetCharCallback(keys.CharCallbackWrapper(handler))

	// mouse events
	w.mouse = mouse.NewWrapper(handler)
	w.SetMouseButtonCallback(w.mouse.Button)
	w.SetCursorPosCallback(w.mouse.Move)
	w.SetScrollCallback(w.mouse.Scroll)
}

func (w *window) onResize(_ *glfw.Window, width, height int) {

}
