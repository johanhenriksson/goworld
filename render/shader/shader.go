package shader

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/types"
)

type Input struct {
	Index int
	Type  types.Type
}

type Inputs map[string]Input

// Input returns the index and type of a shader input by name, and a bool indicating wheter its valid.
func (i Inputs) Input(name string) (int, types.Type, bool) {
	input, exists := i[name]
	return input.Index, input.Type, exists
}

type Bindings map[string]int

// Descriptor returns the index of a descriptor by name, and a bool indicating wheter its valid.
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
	Textures() []texture.Slot
}

type shader struct {
	name     string
	modules  []Module
	inputs   Inputs
	bindings Bindings
	textures []texture.Slot
}

func New(device device.T, path string) T {
	// todo: inputs & descriptors should be obtained from SPIR-V reflection
	detailsPath := fmt.Sprintf("shaders/%s.json", path)
	details, err := ReadDetails(detailsPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load shader details %s: %s", detailsPath, err))
	}

	inputs, err := details.ParseInputs()
	if err != nil {
		panic(fmt.Sprintf("failed to parse shader inputs: %s", err))
	}

	modules := []Module{
		NewModule(device, fmt.Sprintf("shaders/%s.vs.glsl", path), StageVertex),
		NewModule(device, fmt.Sprintf("shaders/%s.fs.glsl", path), StageFragment),
	}

	return &shader{
		name:     path,
		modules:  modules,
		inputs:   inputs,
		bindings: details.Bindings,
		textures: details.Textures,
	}
}

// Name returns the file name of the shader
func (s *shader) Name() string {
	return s.name
}

func (s *shader) Modules() []Module {
	return s.modules
}

// Destroy the shader and its modules.
func (s *shader) Destroy() {
	for _, module := range s.modules {
		module.Destroy()
	}
}

func (s *shader) Input(name string) (int, types.Type, bool) {
	return s.inputs.Input(name)
}

func (s *shader) Textures() []texture.Slot {
	return s.textures
}

func (s *shader) Descriptor(name string) (int, bool) {
	return s.bindings.Descriptor(name)
}
