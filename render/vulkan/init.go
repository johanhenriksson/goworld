package vulkan

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

func init() {
	// glfw event handling must run on the main OS thread
	runtime.LockOSThread()

	// init glfw
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	// initialize vulkan
	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err := vk.Init(); err != nil {
		panic(err)
	}
}
