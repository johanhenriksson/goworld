package main

import (
	"fmt"
	"runtime"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
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
	instance := instance.New("goworld")
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
	defer vk.DestroySurface(instance.Ptr(), surface, nil)

	chain := swapchain.New(window, device, surface)
	defer chain.Destroy()

	input := []int{1, 2, 3}
	stage := buffer.NewShared(device, 8*3)
	defer stage.Destroy()
	stage.Write(input, 0)

	remote := buffer.NewRemote(device, 8*3)
	defer remote.Destroy()

	output := buffer.NewShared(device, 8*3)
	defer output.Destroy()

	pool := command.NewPool(
		device,
		vk.CommandPoolCreateFlags(vk.CommandPoolCreateResetCommandBufferBit),
		vk.QueueFlags(vk.QueueGraphicsBit))
	defer pool.Destroy()

	cmdbuf := pool.Allocate(vk.CommandBufferLevelPrimary)
	defer cmdbuf.Destroy()
	cmdbuf.Begin()
	cmdbuf.CopyBuffer(stage, remote, vk.BufferCopy{SrcOffset: 0, DstOffset: 0, Size: 24})
	cmdbuf.CopyBuffer(remote, output, vk.BufferCopy{SrcOffset: 0, DstOffset: 0, Size: 24})
	cmdbuf.End()

	cmdbuf.SubmitSync(device.GetQueue(0, vk.QueueFlags(vk.QueueGraphicsBit)))

	result := make([]int, 3)
	output.Read(result, 0)
	fmt.Println("read back", result)

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
