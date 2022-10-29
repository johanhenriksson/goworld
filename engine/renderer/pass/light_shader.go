package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type LightDescriptors struct {
	descriptor.Set
	Diffuse  *descriptor.InputAttachment
	Normal   *descriptor.InputAttachment
	Position *descriptor.InputAttachment
	Depth    *descriptor.InputAttachment
	Camera   *descriptor.Uniform[uniform.Camera]
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
			Shader:   shader.New(device, "vk/light"),
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
			Camera: &descriptor.Uniform[uniform.Camera]{
				Stages: vk.ShaderStageFragmentBit,
			},
			Shadow: &descriptor.SamplerArray{
				Stages: vk.ShaderStageFragmentBit,
				Count:  16,
			},
		})
	return mat.Instantiate()
}
