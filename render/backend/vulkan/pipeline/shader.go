package pipeline

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Shader interface {
	Modules() []ShaderModule
	Destroy()
}

type shader struct {
	modules []ShaderModule
}

func NewShader(device device.T, path string) Shader {
	modules := []ShaderModule{
		NewShaderModule(device, fmt.Sprintf("assets/shaders/%s.vert", path), vk.ShaderStageVertexBit),
		NewShaderModule(device, fmt.Sprintf("assets/shaders/%s.frag", path), vk.ShaderStageFragmentBit),
	}
	return &shader{
		modules: modules,
	}
}

func (s *shader) Modules() []ShaderModule {
	return s.modules
}

func (s *shader) Destroy() {
	for _, module := range s.modules {
		module.Destroy()
	}
}
