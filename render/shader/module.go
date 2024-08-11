package shader

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Module interface {
	device.Resource[core1_0.ShaderModule]

	Entrypoint() string
	Stage() ShaderStage
}

type shader_module struct {
	device *device.Device
	ptr    core1_0.ShaderModule
	stage  ShaderStage
}

func NewModule(device *device.Device, path string, stage ShaderStage) Module {
	if device == nil {
		panic("device is nil")
	}

	bytecode, err := LoadOrCompile(path, stage)
	if err != nil {
		panic(err)
	}

	ptr, result, err := device.Ptr().CreateShaderModule(nil, core1_0.ShaderModuleCreateInfo{
		Code: sliceUint32(bytecode),
	})
	if err != nil {
		panic(err)
	}
	if result != core1_0.VKSuccess {
		panic("failed to create shader")
	}
	device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()), core1_0.ObjectTypeShaderModule, path)

	return &shader_module{
		device: device,
		ptr:    ptr,
		stage:  stage,
	}
}

func (b *shader_module) VkType() core1_0.ObjectType { return core1_0.ObjectTypeShaderModule }

func (s *shader_module) Ptr() core1_0.ShaderModule {
	return s.ptr
}

func (s *shader_module) Stage() ShaderStage {
	return s.stage
}

func (s *shader_module) Entrypoint() string {
	return "main"
}

func (s *shader_module) Destroy() {
	s.ptr.Destroy(nil)
	s.ptr = nil
}
