package window

import (
	"fmt"
	"log"
	"runtime"

	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"

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

	Poll()
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
	wnd     *glfw.Window
	backend GlfwBackend
	mouse   mouse.MouseWrapper

	title           string
	width, height   int
	fwidth, fheight int
	scale           float32
}

func New(backend GlfwBackend, args Args) (T, error) {

	// window creation hints.
	for _, hint := range backend.GlfwHints(args) {
		glfw.WindowHint(hint.Hint, hint.Value)
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
		backend: backend,
		title:   args.Title,
		width:   width,
		height:  height,
		fwidth:  fwidth,
		fheight: fheight,
		scale:   scale,
	}

	// setup glfw window for use with the current backend
	if err := backend.GlfwSetup(wnd, args); err != nil {
		return nil, err
	}

	// attach default input handler, if provided
	if args.InputHandler != nil {
		window.SetInputHandler(args.InputHandler)
	}

	// set resize callback
	wnd.SetSizeCallback(window.onResize)

	return window, nil
}

func (w *window) Poll() {
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
	// very unclear if this works with vulkan
	w.width = width
	w.height = height
	w.fwidth, w.fheight = w.wnd.GetFramebufferSize()
	w.scale = float32(w.fwidth) / float32(w.width)

	w.backend.Resize(w.fwidth, w.fheight)
}

func (w *window) SetTitle(title string) {
	w.wnd.SetTitle(title)
	w.title = title
}
