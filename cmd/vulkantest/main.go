package main

import (
	"fmt"
	"runtime"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/instance"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/swapchain"
	"github.com/johanhenriksson/goworld/util"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

func init() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())

	if err := vk.Init(); err != nil {
		panic(err)
	}
}

func main() {
	// create a window
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	window, err := glfw.CreateWindow(500, 500, "moltenvk", nil, nil)
	fmt.Println("window required extensions:", window.GetRequiredInstanceExtensions())
	defer window.Destroy()

	// create instance
	instance := instance.New("gulkan")
	defer instance.Destroy()

	// create device
	physDevices := instance.EnumeratePhysicalDevices()
	device, err := device.New(physDevices[0])
	if err != nil {
		panic(err)
	}
	defer device.Destroy()

	// surface
	surfPtr, err := window.CreateWindowSurface(instance.Ptr(), nil)
	if err != nil {
		panic(err)
	}
	surface := vk.SurfaceFromPointer(surfPtr)

	// window.Show()

	chain := swapchain.New(window, device, surface)
	defer chain.Destroy()

	for !window.ShouldClose() {
		// aquire backbuffer image
		chain.Aquire()

		// cmds := make([]vk.CommandBuffer, 1)
		// vk.AllocateCommandBuffers(device.Ptr(), &vk.CommandBufferAllocateInfo{
		// 	SType: vk.StructureTypeCommandBufferAllocateInfo,
		// }, cmds)

		// draw
		chain.Submit([]vk.CommandBuffer{})

		chain.Present()

		glfw.PollEvents()
	}
}

func GetDeviceNames(devices []vk.PhysicalDevice) []string {
	return util.Map(devices, func(i int, device vk.PhysicalDevice) string {
		var properties vk.PhysicalDeviceProperties
		vk.GetPhysicalDeviceProperties(device, &properties)
		properties.Deref()
		return vk.ToString(properties.DeviceName[:])
	})
}
