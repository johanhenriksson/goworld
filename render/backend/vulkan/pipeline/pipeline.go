package pipeline

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Pipeline]
}

type pipeline struct {
	ptr    vk.Pipeline
	device device.T
}

func New(device device.T, cache vk.PipelineCache, shaders []Shader) T {
	modules := util.Map(shaders, func(i int, shader Shader) vk.PipelineShaderStageCreateInfo {
		return vk.PipelineShaderStageCreateInfo{
			SType:  vk.StructureTypePipelineShaderStageCreateInfo,
			Flags:  vk.PipelineShaderStageCreateFlags(vk.PipelineStageVertexShaderBit),
			Module: shader.Ptr(),
			PName:  util.CString(shader.Entrypoint()),
		}
	})

	info := vk.GraphicsPipelineCreateInfo{
		SType: vk.StructureTypeGraphicsPipelineCreateInfo,

		// Stages
		StageCount: uint32(len(modules)),
		PStages:    modules,

		// Vertex input state

		// Input assembly
		PInputAssemblyState: &vk.PipelineInputAssemblyStateCreateInfo{
			SType:    vk.StructureTypePipelineInputAssemblyStateCreateInfo,
			Topology: vk.PrimitiveTopologyTriangleList,
		},

		// viewport state
		PViewportState: &vk.PipelineViewportStateCreateInfo{
			SType:         vk.StructureTypePipelineViewportStateCreateInfo,
			ViewportCount: 1,
			PViewports:    []vk.Viewport{
				// viewport
			},
			ScissorCount: 1,
			PScissors:    []vk.Rect2D{
				// scissor
			},
		},

		// rasterization state
		PRasterizationState: &vk.PipelineRasterizationStateCreateInfo{
			SType:                   vk.StructureTypePipelineRasterizationStateCreateInfo,
			DepthClampEnable:        vk.False,
			RasterizerDiscardEnable: vk.False,
			PolygonMode:             vk.PolygonModeFill,
			CullMode:                vk.CullModeFlags(vk.CullModeNone),
			FrontFace:               vk.FrontFaceCounterClockwise,
		},

		// multisample
		PMultisampleState: &vk.PipelineMultisampleStateCreateInfo{
			SType:                vk.StructureTypePipelineMultisampleStateCreateInfo,
			RasterizationSamples: vk.SampleCountFlagBits(vk.SampleCount1Bit),
		},

		// depth & stencil

		// color blending
		PColorBlendState: &vk.PipelineColorBlendStateCreateInfo{
			SType:           vk.StructureTypePipelineColorBlendStateCreateInfo,
			LogicOpEnable:   vk.False,
			LogicOp:         vk.LogicOpClear,
			AttachmentCount: 1,
			PAttachments: []vk.PipelineColorBlendAttachmentState{
				{
					BlendEnable:         vk.False,
					ColorBlendOp:        vk.BlendOpAdd,
					SrcColorBlendFactor: vk.BlendFactorZero,
					DstColorBlendFactor: vk.BlendFactorOne,
					AlphaBlendOp:        vk.BlendOpAdd,
					SrcAlphaBlendFactor: vk.BlendFactorZero,
					DstAlphaBlendFactor: vk.BlendFactorOne,
					ColorWriteMask: vk.ColorComponentFlags(
						vk.ColorComponentRBit | vk.ColorComponentGBit |
							vk.ColorComponentBBit | vk.ColorComponentABit),
				},
			},
		},

		// dynamic state: viewport & scissor
		PDynamicState: &vk.PipelineDynamicStateCreateInfo{
			SType:             vk.StructureTypePipelineDynamicStateCreateInfo,
			DynamicStateCount: 2,
			PDynamicStates: []vk.DynamicState{
				vk.DynamicStateScissor,
				vk.DynamicStateViewport,
			},
		},
	}

	var ptrs []vk.Pipeline
	vk.CreateGraphicsPipelines(device.Ptr(), cache, 1, []vk.GraphicsPipelineCreateInfo{info}, nil, ptrs)

	return &pipeline{
		ptr:    ptrs[0],
		device: device,
	}
}

func (p *pipeline) Ptr() vk.Pipeline {
	return p.ptr
}

func (p *pipeline) Destroy() {
	vk.DestroyPipeline(p.device.Ptr(), p.ptr, nil)
	p.ptr = nil
}
