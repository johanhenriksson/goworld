package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/instance"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/swapchain"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	Instance() instance.T
	Device() device.T
	Surface() vk.Surface
	CmdPool() command.Pool
	Destroy()

	GlfwHints(window.Args) []window.GlfwHint
	GlfwSetup(*glfw.Window, window.Args) error

	Resize(int, int)
	Aquire()
	Present()
	Submit([]vk.CommandBuffer)
}

type backend struct {
	appName   string
	deviceIdx int
	instance  instance.T
	device    device.T
	surface   vk.Surface
	swapchain swapchain.T
	cmdpool   command.Pool
}

func New(appName string, deviceIndex int) T {
	return &backend{
		appName:   appName,
		deviceIdx: deviceIndex,
	}
}

func (b *backend) Instance() instance.T  { return b.instance }
func (b *backend) Device() device.T      { return b.device }
func (b *backend) Surface() vk.Surface   { return b.surface }
func (b *backend) CmdPool() command.Pool { return b.cmdpool }

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

	// create instance
	b.instance = instance.New(b.appName)

	// create device
	var err error
	physDevices := b.instance.EnumeratePhysicalDevices()
	b.device, err = device.New(physDevices[b.deviceIdx])
	if err != nil {
		panic(err)
	}

	// surface
	surfPtr, err := w.CreateWindowSurface(b.instance.Ptr(), nil)
	if err != nil {
		panic(err)
	}

	b.surface = vk.SurfaceFromPointer(surfPtr)

	width, height := w.GetFramebufferSize()
	b.swapchain = swapchain.New(b.device, width, height, b.surface)

	b.cmdpool = command.NewPool(
		b.device,
		vk.CommandPoolCreateFlags(vk.CommandPoolCreateResetCommandBufferBit),
		vk.QueueFlags(vk.QueueGraphicsBit))

	return nil
}

func (b *backend) Destroy() {
	if b.cmdpool != nil {
		b.cmdpool.Destroy()
		b.cmdpool = nil
	}
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

func (b *backend) Aquire() {
	b.swapchain.Aquire()
}

func (b *backend) Present() {
	b.swapchain.Present()
}

func (b *backend) Submit(cmdBuffers []vk.CommandBuffer) {
	b.swapchain.Submit(cmdBuffers)
}
