package window

import (
	"fmt"
	"log"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
	Title() string
	SetTitle(string)

	Size() (int, int)
	BufferSize() (int, int)
	Scale() float32

	SwapBuffers()
	ShouldClose() bool

	SetInputHandler(input.Handler)
}

type Args struct {
	Title        string
	Width        int
	Height       int
	Vsync        bool
	Debug        bool
	InputHandler input.Handler
}

type window struct {
	wnd   *glfw.Window
	mouse mouse.MouseWrapper

	title           string
	width, height   int
	fwidth, fheight int
	scale           float32
}

func New(args Args) (T, error) {
	// window creation hints.
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 1)

	if args.Debug {
		glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	}

	// create a new GLFW window
	wnd, err := glfw.CreateWindow(args.Width, args.Height, args.Title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create glfw window: %w", err)
	}

	// retrieve window & framebuffer size
	width, height := wnd.GetSize()
	fwidth, fheight := wnd.GetFramebufferSize()
	scale := float32(fwidth) / float32(width)
	log.Printf("Created window. Size %dx%d, Buffer Size %dx%d. Scale = %.0f%%\n",
		width, height, fwidth, fheight, scale*100)

	window := &window{
		wnd:     wnd,
		title:   args.Title,
		width:   width,
		height:  height,
		fwidth:  fwidth,
		fheight: fheight,
		scale:   scale,
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

	// attach default input handler, if provided
	if args.InputHandler != nil {
		window.SetInputHandler(args.InputHandler)
	}

	// set resize callback
	wnd.SetSizeCallback(window.onResize)

	// set up debugging
	if args.Debug {
		var flags int32
		gl.GetIntegerv(gl.CONTEXT_FLAGS, &flags)
		if flags&gl.CONTEXT_FLAG_DEBUG_BIT == gl.CONTEXT_FLAG_DEBUG_BIT {
			gl.Enable(gl.DEBUG_OUTPUT)
			gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
			gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)
			gl.DebugMessageCallback(window.onDebugMessage, nil)
		} else {
			fmt.Println("warning: failed to enable opengl debugging")
		}
	}

	return window, nil
}

func (w *window) SwapBuffers() {
	w.wnd.SwapBuffers()
	glfw.PollEvents()
}

func (w *window) Size() (int, int)       { return w.width, w.height }
func (w *window) BufferSize() (int, int) { return w.fwidth, w.fheight }
func (w *window) Scale() float32         { return w.scale }
func (w *window) ShouldClose() bool      { return w.wnd.ShouldClose() }
func (w *window) Title() string          { return w.title }

func (w *window) SetInputHandler(handler input.Handler) {
	// keyboard events
	w.wnd.SetKeyCallback(keys.KeyCallbackWrapper(handler))
	w.wnd.SetCharCallback(keys.CharCallbackWrapper(handler))

	// mouse events
	w.mouse = mouse.NewWrapper(handler)
	w.wnd.SetMouseButtonCallback(w.mouse.Button)
	w.wnd.SetCursorPosCallback(w.mouse.Move)
	w.wnd.SetScrollCallback(w.mouse.Scroll)
}

func (w *window) onResize(_ *glfw.Window, width, height int) {
	w.width = width
	w.height = height
	w.fwidth, w.fheight = w.wnd.GetFramebufferSize()
	w.scale = float32(w.fwidth) / float32(w.width)
}

func (w *window) onDebugMessage(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {
	// todo: proper messages
	// see https://learnopengl.com/In-Practice/Debugging
	fmt.Println("GL Debug:", message)
}

func (w *window) SetTitle(title string) {
	w.wnd.SetTitle(title)
	w.title = title
}
