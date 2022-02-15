package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/instance"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/swapchain"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

type VkVertex struct {
	X, Y, Z float32
	R, G, B float32
}

type T interface {
	Instance() instance.T
	Device() device.T
	Surface() vk.Surface
	Swapchain() swapchain.T
	Destroy()

	GlfwHints(window.Args) []window.GlfwHint
	GlfwSetup(*glfw.Window, window.Args) error

	Resize(int, int)
	Aquire() (swapchain.Context, error)
	Present()
}

type backend struct {
	appName   string
	deviceIdx int
	swapcount int
	instance  instance.T
	device    device.T
	surface   vk.Surface
	swapchain swapchain.T
}

func New(appName string, deviceIndex int) T {
	return &backend{
		appName:   appName,
		deviceIdx: deviceIndex,
		swapcount: 2,
	}
}

func (b *backend) Instance() instance.T   { return b.instance }
func (b *backend) Device() device.T       { return b.device }
func (b *backend) Surface() vk.Surface    { return b.surface }
func (b *backend) Swapchain() swapchain.T { return b.swapchain }

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

	fmt.Println("window required extensions:", w.GetRequiredInstanceExtensions())

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

	// allocate swapchain
	b.swapchain = swapchain.New(b.device, b.swapcount, b.surface, surfaceFormat)

	return nil
}

func (b *backend) Destroy() {
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
	b.swapchain.Resize(width, height)
}

func (b *backend) Aquire() (swapchain.Context, error) {
	return b.swapchain.Aquire()
}

func (b *backend) Present() {
	b.swapchain.Present()
}
