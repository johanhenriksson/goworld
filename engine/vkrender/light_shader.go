package vkrender

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_shader"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	vk "github.com/vulkan-go/vulkan"
)

type LightDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[CameraData]
	Light    *descriptor.UniformArray[light.Descriptor]
	Diffuse  *descriptor.InputAttachment
	Normal   *descriptor.InputAttachment
	Position *descriptor.InputAttachment
	Depth    *descriptor.InputAttachment
}

type LShader struct {
	vk_shader.T[*LightDescriptors]
}

func NewLightShader(backend vulkan.T, pass renderpass.T) vk_shader.T[*LightDescriptors] {
	return vk_shader.New(
		backend,
		vk_shader.Args{
			Path:     "vk/light",
			Frames:   1,
			Pass:     pass,
			Subpass:  "lighting",
			Pointers: vertex.ParsePointers(vertex.T{}),
			Attributes: shader.AttributeMap{
				"position": {
					Loc:  0,
					Type: types.Float,
				},
			},
			Constants: []pipeline.PushConstant{
				{
					Stages: vk.ShaderStageFragmentBit,
					Offset: 0,
					Size:   4,
				},
			},
		},
		&LightDescriptors{
			Diffuse: &descriptor.InputAttachment{
				Binding: 0,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Normal: &descriptor.InputAttachment{
				Binding: 1,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Position: &descriptor.InputAttachment{
				Binding: 2,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Depth: &descriptor.InputAttachment{
				Binding: 3,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Camera: &descriptor.Uniform[CameraData]{
				Binding: 4,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Light: &descriptor.UniformArray[light.Descriptor]{
				Binding: 5,
				Size:    10,
				Stages:  vk.ShaderStageFragmentBit,
			},
		})
}
