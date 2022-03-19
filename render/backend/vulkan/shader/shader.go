package shader

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Input struct {
	Index int
	Type  types.Type
}

type Inputs map[string]Input
type Descriptors map[string]int

type T interface {
	Modules() []Module
	Destroy()
	Input(name string) (int, types.Type, bool)
	Descriptor(name string) (int, bool)
}

type shader struct {
	modules     []Module
	inputs      Inputs
	descriptors Descriptors
}

func New(device device.T, path string, inputs Inputs, descriptors Descriptors) T {
	// todo: inputs & descriptors should be obtained from SPIR-V reflection
	modules := []Module{
		NewModule(device, fmt.Sprintf("assets/shaders/%s.vert", path), vk.ShaderStageVertexBit),
		NewModule(device, fmt.Sprintf("assets/shaders/%s.frag", path), vk.ShaderStageFragmentBit),
	}
	return &shader{
		modules:     modules,
		inputs:      inputs,
		descriptors: descriptors,
	}
}

func (s *shader) Modules() []Module {
	return s.modules
}

func (s *shader) Destroy() {
	for _, module := range s.modules {
		module.Destroy()
	}
}

func (s *shader) Input(name string) (int, types.Type, bool) {
	input, exists := s.inputs[name]
	return input.Index, input.Type, exists
}

func (s *shader) Descriptor(name string) (int, bool) {
	index, exists := s.descriptors[name]
	return index, exists
}
