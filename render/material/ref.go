package material

import (
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type Ref interface {
	cache.Item
}

type Standard Instance[*Descriptors]

type Descriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type Def struct {
	Shader       string
	Subpass      renderpass.Name
	VertexFormat any
}

func FromDef(backend vulkan.T, rpass renderpass.T, def *Def) Standard {
	desc := &Descriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: vk.ShaderStageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: vk.ShaderStageAll,
			Size:   100,
		},
		Textures: &descriptor.SamplerArray{
			Stages: vk.ShaderStageFragmentBit,
			Count:  100,
		},
	}

	pointers := vertex.ParsePointers(def.VertexFormat)

	return New(
		backend.Device(),
		Args{
			Shader:     shader.New(backend.Device(), def.Shader),
			Pass:       rpass,
			Subpass:    def.Subpass,
			Pointers:   pointers,
			DepthTest:  true,
			DepthWrite: true,
		},
		desc).Instantiate()
}