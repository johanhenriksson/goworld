package shader

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Module interface {
	device.Resource[core1_0.ShaderModule]

	Entrypoint() string
	Stage() core1_0.ShaderStageFlags
}

type shader_module struct {
	device device.T
	ptr    core1_0.ShaderModule
	stage  core1_0.ShaderStageFlags
}

func NewModule(device device.T, path string, stage core1_0.ShaderStageFlags) Module {
	bytecode, err := LoadOrCompile(path)
	if err != nil {
		panic(err)
	}

	ptr, _, err := device.Ptr().CreateShaderModule(nil, core1_0.ShaderModuleCreateInfo{
		Code: sliceUint32(bytecode),
	})
	if err != nil {
		panic(err)
	}

	return &shader_module{
		device: device,
		ptr:    ptr,
		stage:  stage,
	}
}

func (s *shader_module) Ptr() core1_0.ShaderModule {
	return s.ptr
}

func (s *shader_module) Stage() core1_0.ShaderStageFlags {
	return s.stage
}

func (s *shader_module) Entrypoint() string {
	return "main"
}

func (s *shader_module) Destroy() {
	s.ptr.Destroy(nil)
	s.ptr = nil
}
