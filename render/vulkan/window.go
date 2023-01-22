package vulkan

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

type ResizeHandler func(width, height int)

type Window interface {
	Title() string
	SetTitle(string)

	Size() (int, int)
	Scale() float32

	Poll()
	ShouldClose() bool
	Destroy()

	SetInputHandler(input.Handler)
	// SetResizeHandler(ResizeHandler)

	Swapchain() swapchain.T
}

type WindowArgs struct {
	Title         string
	Width         int
	Height        int
	Vsync         bool
	Debug         bool
	InputHandler  input.Handler
	ResizeHandler ResizeHandler
}

type window struct {
	wnd     *glfw.Window
	backend T
	mouse   mouse.MouseWrapper

	title         string
	width, height int
	scale         float32
	swap          swapchain.T
	surface       vk.Surface
}

func NewWindow(backend T, args WindowArgs) (Window, error) {
	// window creation hints.
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)

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
	log.Printf("Created window with size %dx%d and content scale %.0f%%\n",
		width, height, scale*100)

	// create window surface
	surfPtr, err := wnd.CreateWindowSurface(backend.Instance().Ptr(), nil)
	if err != nil {
		panic(err)
	}

	surface := vk.SurfaceFromPointer(surfPtr)
	surfaceFormat := backend.Device().GetSurfaceFormats(surface)[0]

	// allocate swapchain
	swap := swapchain.New(backend.Device(), backend.Frames(), width, height, surface, surfaceFormat)

	window := &window{
		wnd:     wnd,
		backend: backend,
		title:   args.Title,
		width:   width,
		height:  height,
		scale:   scale,
		swap:    swap,
		surface: surface,
	}

	// attach default input handler, if provided
	if args.InputHandler != nil {
		window.SetInputHandler(args.InputHandler)
	}

	// set resize callback
	wnd.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		args.ResizeHandler(width, height)
	})

	return window, nil
}

func (w *window) Poll() {
	glfw.PollEvents()
}

func (w *window) Size() (int, int)  { return w.width, w.height }
func (w *window) Scale() float32    { return w.scale }
func (w *window) ShouldClose() bool { return w.wnd.ShouldClose() }
func (w *window) Title() string     { return w.title }

func (w *window) Swapchain() swapchain.T { return w.swap }

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

func (w *window) SetTitle(title string) {
	w.wnd.SetTitle(title)
	w.title = title
}

func (w *window) Destroy() {
	w.swap.Destroy()
	vk.DestroySurface(w.backend.Instance().Ptr(), w.surface, nil)
}
