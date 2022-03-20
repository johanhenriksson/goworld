package vkrender

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
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
	Shadow   *descriptor.Sampler
}

func NewLightShader(device device.T, pass renderpass.T) material.Instance[*LightDescriptors] {
	mat := material.New(
		device,
		material.Args{
			Shader: shader.New(
				device,
				"vk/light",
				shader.Inputs{
					"position": {
						Index: 0,
						Type:  types.Float,
					},
				},
				shader.Descriptors{
					"Diffuse":  0,
					"Normal":   1,
					"Position": 2,
					"Depth":    3,
					"Camera":   4,
					"Light":    5,
					"Shadow":   6,
				},
			),
			Pass:     pass,
			Subpass:  "lighting",
			Pointers: vertex.ParsePointers(vertex.T{}),
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
				Stages: vk.ShaderStageFragmentBit,
			},
			Normal: &descriptor.InputAttachment{
				Stages: vk.ShaderStageFragmentBit,
			},
			Position: &descriptor.InputAttachment{
				Stages: vk.ShaderStageFragmentBit,
			},
			Depth: &descriptor.InputAttachment{
				Stages: vk.ShaderStageFragmentBit,
			},
			Camera: &descriptor.Uniform[CameraData]{
				Stages: vk.ShaderStageFragmentBit,
			},
			Light: &descriptor.UniformArray[light.Descriptor]{
				Size:   10,
				Stages: vk.ShaderStageFragmentBit,
			},
			Shadow: &descriptor.Sampler{
				Stages: vk.ShaderStageFragmentBit,
			},
		})
	return mat.Instantiate()
}
