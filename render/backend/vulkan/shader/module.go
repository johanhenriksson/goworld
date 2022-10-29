package shader

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Module interface {
	device.Resource[vk.ShaderModule]

	Entrypoint() string
	Stage() vk.ShaderStageFlagBits
}

type shader_module struct {
	device device.T
	ptr    vk.ShaderModule
	stage  vk.ShaderStageFlagBits
}

func NewModule(device device.T, path string, stage vk.ShaderStageFlagBits) Module {
	bytecode, err := LoadOrCompile(path)
	if err != nil {
		panic(err)
	}

	info := vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(bytecode)),
		PCode:    sliceUint32(bytecode),
	}

	var ptr vk.ShaderModule
	vk.CreateShaderModule(device.Ptr(), &info, nil, &ptr)

	return &shader_module{
		device: device,
		ptr:    ptr,
		stage:  stage,
	}
}

func (s *shader_module) Ptr() vk.ShaderModule {
	return s.ptr
}

func (s *shader_module) Stage() vk.ShaderStageFlagBits {
	return s.stage
}

func (s *shader_module) Entrypoint() string {
	return "main"
}

func (s *shader_module) Destroy() {
	vk.DestroyShaderModule(s.device.Ptr(), s.ptr, nil)
	s.ptr = nil
}
