package vulkan

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
	"github.com/vkngwrapper/extensions/v2/khr_surface"
	khr_surface_driver "github.com/vkngwrapper/extensions/v2/khr_surface/driver"
)

type ResizeHandler func(width, height int)

type Window interface {
	Target

	Title() string
	SetTitle(string)

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
	Frames        int
	Vsync         bool
	Debug         bool
	InputHandler  input.Handler
	ResizeHandler ResizeHandler
}

type window struct {
	T
	wnd   *glfw.Window
	mouse mouse.MouseWrapper

	title         string
	width, height int
	frames        int
	scale         float32
	swap          swapchain.T
	surface       khr_surface.Surface
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
	surfPtr, err := wnd.CreateWindowSurface((*driver.VkInstance)(unsafe.Pointer(backend.Instance().Ptr().Handle())), nil)
	if err != nil {
		panic(err)
	}

	surfaceHandle := (*khr_surface_driver.VkSurfaceKHR)(unsafe.Pointer(surfPtr))
	surfaceExt := khr_surface.CreateExtensionFromInstance(backend.Instance().Ptr())
	surface, err := surfaceExt.CreateSurfaceFromHandle(*surfaceHandle)
	if err != nil {
		panic(err)
	}

	surfaceFormat, _, err := surface.PhysicalDeviceSurfaceFormats(backend.Device().Physical())
	if err != nil {
		panic(err)
	}

	// allocate swapchain
	swap := swapchain.New(backend.Device(), args.Frames, width, height, surface, surfaceFormat[0])

	window := &window{
		T:       backend,
		wnd:     wnd,
		title:   args.Title,
		width:   width,
		height:  height,
		frames:  args.Frames,
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
		window.width = width
		window.height = height
		window.swap.Resize(width, height)
	})

	return window, nil
}

func (w *window) Poll() {
	glfw.PollEvents()
}

func (w *window) Width() int        { return w.width }
func (w *window) Height() int       { return w.height }
func (w *window) Frames() int       { return w.frames }
func (w *window) Scale() float32    { return w.scale }
func (w *window) ShouldClose() bool { return w.wnd.ShouldClose() }
func (w *window) Title() string     { return w.title }

func (w *window) Surfaces() []image.T           { return w.swap.Images() }
func (w *window) SurfaceFormat() core1_0.Format { return w.swap.SurfaceFormat() }
func (w *window) Swapchain() swapchain.T        { return w.swap }

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

func (w *window) Aquire() (swapchain.Context, error) {
	return w.swap.Aquire()
}

func (w *window) Destroy() {
	w.swap.Destroy()
	// vk.DestroySurface(w.T.Instance().Ptr(), w.surface, nil)
}
