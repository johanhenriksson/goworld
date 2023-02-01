package shader

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/types"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Input struct {
	Index int
	Type  types.Type
}

type Inputs map[string]Input

func (i Inputs) Input(name string) (int, types.Type, bool) {
	input, exists := i[name]
	return input.Index, input.Type, exists
}

type Bindings map[string]int

func (d Bindings) Descriptor(name string) (int, bool) {
	index, exists := d[name]
	return index, exists
}

type T interface {
	Name() string
	Modules() []Module
	Destroy()
	Input(name string) (int, types.Type, bool)
	Descriptor(name string) (int, bool)
}

type shader struct {
	name     string
	modules  []Module
	inputs   Inputs
	bindings Bindings
}

func New(device device.T, path string) T {
	// todo: inputs & descriptors should be obtained from SPIR-V reflection
	details, err := ReadDetails(fmt.Sprintf("assets/shaders/%s.json", path))
	if err != nil {
		panic(fmt.Sprintf("failed to load shader details: %s", err))
	}

	inputs, err := details.ParseInputs()
	if err != nil {
		panic(fmt.Sprintf("failed to parse shader inputs: %s", err))
	}

	modules := []Module{
		NewModule(device, fmt.Sprintf("assets/shaders/%s.vert", path), core1_0.StageVertex),
		NewModule(device, fmt.Sprintf("assets/shaders/%s.frag", path), core1_0.StageFragment),
	}

	return &shader{
		name:     path,
		modules:  modules,
		inputs:   inputs,
		bindings: details.Bindings,
	}
}

func (s *shader) Name() string {
	return s.name
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
	return s.inputs.Input(name)
}

func (s *shader) Descriptor(name string) (int, bool) {
	return s.bindings.Descriptor(name)
}
