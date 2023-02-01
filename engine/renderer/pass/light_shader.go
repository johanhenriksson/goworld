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

	"github.com/vkngwrapper/core/v2/core1_0"
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

type LightShader material.Instance[*LightDescriptors]

func NewLightShader(device device.T, pool descriptor.Pool, pass renderpass.T) LightShader {
	mat := material.New(
		device,
		material.Args{
			Shader:   shader.New(device, "vk/light"),
			Pass:     pass,
			Subpass:  LightingSubpass,
			Pointers: vertex.ParsePointers(vertex.T{}),
			Constants: []pipeline.PushConstant{
				{
					Stages: core1_0.StageFragment,
					Type:   LightConst{},
				},
			},
			DepthTest: true,
		},
		&LightDescriptors{
			Diffuse: &descriptor.InputAttachment{
				Stages: core1_0.StageFragment,
			},
			Normal: &descriptor.InputAttachment{
				Stages: core1_0.StageFragment,
			},
			Position: &descriptor.InputAttachment{
				Stages: core1_0.StageFragment,
			},
			Depth: &descriptor.InputAttachment{
				Stages: core1_0.StageFragment,
			},
			Camera: &descriptor.Uniform[uniform.Camera]{
				Stages: core1_0.StageFragment,
			},
			Shadow: &descriptor.SamplerArray{
				Stages: core1_0.StageFragment,
				Count:  16,
			},
		})
	return mat.Instantiate(pool)
}
