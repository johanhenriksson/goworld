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

	title         string
	width, height int
	scale         float32
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

	// find the scaling of the current monitor
	// what do we do if the user moves it to a different monitor with different scaling?
	monitor := GetCurrentMonitor(wnd)
	scale, _ := monitor.GetContentScale()

	// retrieve window & framebuffer size
	width, height := wnd.GetFramebufferSize()
	Scale = scale
	log.Printf("Created window with size %dx%d and content scale %.0f%%\n",
		width, height, scale*100)

	window := &window{
		wnd:     wnd,
		backend: backend,
		title:   args.Title,
		width:   width,
		height:  height,
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
	wnd.SetFramebufferSizeCallback(window.onResize)

	return window, nil
}

func (w *window) Poll() {
	glfw.PollEvents()
}

func (w *window) Size() (int, int)  { return w.width, w.height }
func (w *window) Scale() float32    { return w.scale }
func (w *window) ShouldClose() bool { return w.wnd.ShouldClose() }
func (w *window) Title() string     { return w.title }

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
	w.width, w.height = w.wnd.GetFramebufferSize()
	w.backend.Resize(w.width, w.height)
}

func (w *window) SetTitle(title string) {
	w.wnd.SetTitle(title)
	w.title = title
}
