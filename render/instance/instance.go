package instance

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vkngwrapper/core/v2"
	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/khr_portability_enumeration"
)

type Instance struct {
	ptr core1_0.Instance
}

func New(appName string) *Instance {
	loader, err := core.CreateLoaderFromProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err != nil {
		panic(err)
	}
	log.Println("creating instance with extensions", extensions)
	handle, _, err := loader.CreateInstance(nil, core1_0.InstanceCreateInfo{
		APIVersion:            common.APIVersion(common.CreateVersion(1, 2, 295)),
		ApplicationName:       appName,
		ApplicationVersion:    common.CreateVersion(0, 1, 0),
		EngineName:            "goworld",
		EngineVersion:         common.CreateVersion(0, 2, 1),
		EnabledLayerNames:     layers,
		EnabledExtensionNames: extensions,

		Flags: khr_portability_enumeration.InstanceCreateEnumeratePortability,
	})
	if err != nil {
		panic(err)
	}
	return &Instance{
		ptr: handle,
	}
}

func (i *Instance) Ptr() core1_0.Instance {
	return i.ptr
}

func (i *Instance) Destroy() {
	i.ptr.Destroy(nil)
	i.ptr = nil
}

func (i *Instance) EnumeratePhysicalDevices() []core1_0.PhysicalDevice {
	r, _, err := i.ptr.EnumeratePhysicalDevices()
	if err != nil {
		panic(err)
	}
	return r
}
