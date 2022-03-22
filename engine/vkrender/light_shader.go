package vkrender

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type LightDescriptors struct {
	descriptor.Set
	Diffuse  *descriptor.InputAttachment
	Normal   *descriptor.InputAttachment
	Position *descriptor.InputAttachment
	Depth    *descriptor.InputAttachment
	Camera   *descriptor.Uniform[CameraData]
	Shadow   *descriptor.SamplerArray
}

type LightConst struct {
	ViewProj    mat4.T
	Color       color.T
	Position    vec4.T
	Type        light.Type
	Shadowmap   uint32
	Range       float32
	Intensity   float32
	Attenuation light.Attenuation
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
					"Shadow":   5,
				},
			),
			Pass:     pass,
			Subpass:  "lighting",
			Pointers: vertex.ParsePointers(vertex.T{}),
			Constants: []pipeline.PushConstant{
				{
					Stages: vk.ShaderStageFragmentBit,
					Type:   LightConst{},
				},
			},
			DepthTest: true,
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
			Shadow: &descriptor.SamplerArray{
				Stages: vk.ShaderStageFragmentBit,
				Count:  16,
			},
		})
	return mat.Instantiate()
}
