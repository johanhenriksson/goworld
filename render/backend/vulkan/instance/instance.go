package instance

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

var extensions = []string{
	vk.KhrSurfaceExtensionName,
	vk.ExtDebugReportExtensionName,
	"VK_EXT_metal_surface",
}

var layers = []string{
	"VK_LAYER_KHRONOS_validation",
	// "VK_LAYER_LUNARG_api_dump",
}

type T interface {
	device.Resource[vk.Instance]
	EnumeratePhysicalDevices() []vk.PhysicalDevice
	GetDevice(int) device.T
}

type instance struct {
	ptr vk.Instance
}

func New(appName string) T {
	appInfo := vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PApplicationName:   util.CString(appName),
		ApplicationVersion: vk.MakeVersion(0, 1, 0),
		PEngineName:        util.CString("goworld"),
		EngineVersion:      vk.MakeVersion(0, 2, 0),
		ApiVersion:         vk.MakeVersion(1, 1, 0),
	}

	createInfo := vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        &appInfo,
		PpEnabledExtensionNames: util.CStrings(extensions),
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledLayerNames:     util.CStrings(layers),
		EnabledLayerCount:       uint32(len(layers)),
	}

	var ptr vk.Instance
	r := vk.CreateInstance(&createInfo, nil, &ptr)
	if r != vk.Success {
		panic(fmt.Sprintf("create instance returned %d", r))
	}

	if err := vk.InitInstance(ptr); err != nil {
		panic(err)
	}

	return &instance{
		ptr: ptr,
	}
}

func (i *instance) Ptr() vk.Instance {
	return i.ptr
}

func (i *instance) Destroy() {
	vk.DestroyInstance(i.ptr, nil)
	i.ptr = nil
}

func (i *instance) GetDevice(index int) device.T {
	physDevices := i.EnumeratePhysicalDevices()
	device, err := device.New(physDevices[index])
	if err != nil {
		panic(err)
	}
	return device
}

func (i *instance) EnumeratePhysicalDevices() []vk.PhysicalDevice {
	count := uint32(0)
	vk.EnumeratePhysicalDevices(i.ptr, &count, nil)
	devices := make([]vk.PhysicalDevice, count)
	vk.EnumeratePhysicalDevices(i.ptr, &count, devices)
	return devices
}
