package instance

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vkngwrapper/core/v2"
	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type T interface {
	EnumeratePhysicalDevices() []core1_0.PhysicalDevice
	Destroy()
	Ptr() core1_0.Instance
}

type instance struct {
	ptr core1_0.Instance
}

func New(appName string) T {
	loader, err := core.CreateLoaderFromProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err != nil {
		panic(err)
	}
	handle, _, err := loader.CreateInstance(nil, core1_0.InstanceCreateInfo{
		APIVersion:            common.APIVersion(common.CreateVersion(1, 1, 0)),
		ApplicationName:       appName,
		ApplicationVersion:    common.CreateVersion(0, 1, 0),
		EngineName:            "goworld",
		EngineVersion:         common.CreateVersion(0, 2, 1),
		EnabledLayerNames:     layers,
		EnabledExtensionNames: extensions,
	})
	if err != nil {
		panic(err)
	}
	return &instance{
		ptr: handle,
	}
}

func (i *instance) Ptr() core1_0.Instance {
	return i.ptr
}

func (i *instance) Destroy() {
	i.ptr.Destroy(nil)
	i.ptr = nil
}

func (i *instance) EnumeratePhysicalDevices() []core1_0.PhysicalDevice {
	r, _, err := i.ptr.EnumeratePhysicalDevices()
	if err != nil {
		panic(err)
	}
	return r
}
