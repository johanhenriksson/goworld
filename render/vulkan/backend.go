package vulkan

import (
	"log"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/johanhenriksson/goworld/render/vulkan/instance"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	window.GlfwBackend

	Instance() instance.T
	Device() device.T
	Swapchain() swapchain.T
	Frames() int
	Width() int
	Height() int
	Destroy()

	Worker(int) command.Worker
	Transferer() command.Worker
}

type backend struct {
	appName   string
	deviceIdx int
	frames    int
	width     int
	height    int
	instance  instance.T
	device    device.T
	surface   vk.Surface
	swapchain swapchain.T
	transfer  command.Worker
	workers   []command.Worker
}

func New(appName string, deviceIndex int) T {
	return &backend{
		appName:   appName,
		deviceIdx: deviceIndex,
		frames:    3,
	}
}

func (b *backend) Instance() instance.T   { return b.instance }
func (b *backend) Device() device.T       { return b.device }
func (b *backend) Swapchain() swapchain.T { return b.swapchain }
func (b *backend) Frames() int            { return b.frames }
func (b *backend) Width() int             { return b.width }
func (b *backend) Height() int            { return b.height }

func (b *backend) Transferer() command.Worker {
	return b.transfer
}

func (b *backend) Worker(frame int) command.Worker {
	return b.workers[frame]
}

func (b *backend) GlfwHints(args window.Args) []window.GlfwHint {
	return []window.GlfwHint{
		{Hint: glfw.ClientAPI, Value: glfw.NoAPI},
	}
}

func (b *backend) GlfwSetup(w *glfw.Window, args window.Args) error {
	// initialize vulkan
	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err := vk.Init(); err != nil {
		panic(err)
	}

	log.Println("window required extensions:", w.GetRequiredInstanceExtensions())

	// create instance * device
	b.instance = instance.New(b.appName)
	b.device = b.instance.GetDevice(b.deviceIdx)

	// surface
	surfPtr, err := w.CreateWindowSurface(b.instance.Ptr(), nil)
	if err != nil {
		panic(err)
	}

	b.surface = vk.SurfaceFromPointer(surfPtr)
	surfaceFormat := b.device.GetSurfaceFormats(b.surface)[0]

	b.width, b.height = w.GetFramebufferSize()

	// allocate swapchain
	b.swapchain = swapchain.New(b.device, b.frames, b.width, b.height, b.surface, surfaceFormat)

	// transfer worker
	b.transfer = command.NewWorker(b.device, vk.QueueFlags(vk.QueueTransferBit))

	// per frame graphics workers
	b.workers = make([]command.Worker, b.frames)
	for i := 0; i < b.frames; i++ {
		b.workers[i] = command.NewWorker(b.device, vk.QueueFlags(vk.QueueGraphicsBit))
	}

	// init global descriptor pool
	descriptor.InitGlobalPool(b.device)

	return nil
}

func (b *backend) Destroy() {
	// clean up global descriptor pool
	descriptor.DestroyGlobalPool()

	for i := 0; i < b.frames; i++ {
		b.workers[i].Destroy()
	}
	b.transfer.Destroy()

	if b.swapchain != nil {
		b.swapchain.Destroy()
		b.swapchain = nil
	}
	if b.surface != nil {
		vk.DestroySurface(b.instance.Ptr(), b.surface, nil)
		b.surface = nil
	}
	if b.device != nil {
		b.device.Destroy()
		b.device = nil
	}
	if b.instance != nil {
		b.instance.Destroy()
		b.instance = nil
	}
}

func (b *backend) Resize(width, height int) {
	b.width = width
	b.height = height
	b.swapchain.Resize(width, height)
}

func (b *backend) Aquire() (swapchain.Context, error) {
	return b.swapchain.Aquire()
}

func (b *backend) Present() {
	b.swapchain.Present()
}