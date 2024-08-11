package glfw

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/window"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/instance"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
	"github.com/vkngwrapper/extensions/v2/khr_surface"
	khr_surface_driver "github.com/vkngwrapper/extensions/v2/khr_surface/driver"
)

type glfwWindow struct {
	wnd   *glfw.Window
	mouse mouse.MouseWrapper

	title         string
	width, height int
	frames        int
	scale         float32
	swap          swapchain.T
	surface       khr_surface.Surface
}

func NewWindow(vulkan instance.T, device device.T, args window.WindowArgs) (window.Window, error) {
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
	surfPtr, err := wnd.CreateWindowSurface((*driver.VkInstance)(unsafe.Pointer(vulkan.Ptr().Handle())), nil)
	if err != nil {
		panic(err)
	}

	surfaceHandle := (*khr_surface_driver.VkSurfaceKHR)(unsafe.Pointer(surfPtr))
	surfaceExt := khr_surface.CreateExtensionFromInstance(vulkan.Ptr())
	surface, err := surfaceExt.CreateSurfaceFromHandle(*surfaceHandle)
	if err != nil {
		panic(err)
	}

	surfaceFormat, _, err := surface.PhysicalDeviceSurfaceFormats(device.Physical())
	if err != nil {
		panic(err)
	}

	// allocate swapchain
	swap := swapchain.New(device, args.Frames, width, height, surface, surfaceFormat[0])

	window := &glfwWindow{
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
		// update window scaling
		monitor := GetCurrentMonitor(wnd)
		window.scale, _ = monitor.GetContentScale()

		window.width = width
		window.height = height
		window.swap.Resize(width, height)
	})

	return window, nil
}

func (w *glfwWindow) Poll() {
	glfw.PollEvents()
}

func (w *glfwWindow) Size() engine.TargetSize {
	return engine.TargetSize{
		Width:  w.width,
		Height: w.height,
		Frames: w.frames,
		Scale:  w.scale,
	}
}

func (w *glfwWindow) Width() int        { return w.width }
func (w *glfwWindow) Height() int       { return w.height }
func (w *glfwWindow) Frames() int       { return w.frames }
func (w *glfwWindow) Scale() float32    { return w.scale }
func (w *glfwWindow) ShouldClose() bool { return w.wnd.ShouldClose() }
func (w *glfwWindow) Title() string     { return w.title }

func (w *glfwWindow) Surfaces() []image.T           { return w.swap.Images() }
func (w *glfwWindow) SurfaceFormat() core1_0.Format { return w.swap.SurfaceFormat() }

func (w *glfwWindow) SetInputHandler(handler input.Handler) {
	// keyboard events
	w.wnd.SetKeyCallback(keys.KeyCallbackWrapper(handler))
	w.wnd.SetCharCallback(keys.CharCallbackWrapper(handler))

	// mouse events
	w.mouse = mouse.NewWrapper(handler)
	w.wnd.SetMouseButtonCallback(w.mouse.Button)
	w.wnd.SetCursorPosCallback(w.mouse.Move)
	w.wnd.SetScrollCallback(w.mouse.Scroll)
}

func (w *glfwWindow) SetTitle(title string) {
	w.wnd.SetTitle(title)
	w.title = title
}

func (w *glfwWindow) Aquire(worker command.Worker) (*swapchain.Context, error) {
	return w.swap.Aquire(worker)
}

func (w *glfwWindow) Present(worker command.Worker, ctx *swapchain.Context) {
	w.swap.Present(worker, ctx)
}

func (w *glfwWindow) Destroy() {
	w.swap.Destroy()
	w.surface.Destroy(nil)
}
