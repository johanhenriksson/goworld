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

func New(device device.T, cache vk.PipelineCache, layout Layout, pass Pass, shaders []Shader) T {
	modules := util.Map(shaders, func(i int, shader Shader) vk.PipelineShaderStageCreateInfo {
		return vk.PipelineShaderStageCreateInfo{
			SType:  vk.StructureTypePipelineShaderStageCreateInfo,
			Module: shader.Ptr(),
			PName:  util.CString(shader.Entrypoint()),
			Stage:  shader.Stage(),
		}
	})

	info := vk.GraphicsPipelineCreateInfo{
		SType: vk.StructureTypeGraphicsPipelineCreateInfo,

		// layout
		Layout: layout.Ptr(),

		// render pass
		RenderPass: pass.Ptr(),

		// Stages
		StageCount: uint32(len(modules)),
		PStages:    modules,

		// Vertex input state
		PVertexInputState: &vk.PipelineVertexInputStateCreateInfo{
			SType:                         vk.StructureTypePipelineVertexInputStateCreateInfo,
			VertexBindingDescriptionCount: 1,
			PVertexBindingDescriptions: []vk.VertexInputBindingDescription{
				{
					Binding:   0,
					Stride:    10 * 4,
					InputRate: vk.VertexInputRateVertex,
				},
			},
			VertexAttributeDescriptionCount: 2,
			PVertexAttributeDescriptions: []vk.VertexInputAttributeDescription{
				{
					Binding:  0,
					Location: 0, // vec3 position
					Format:   vk.FormatR32g32b32Sfloat,
					Offset:   0,
				},
				{
					Binding:  0,
					Location: 1, // vec3 color
					Format:   vk.FormatR32g32b32Sfloat,
					Offset:   6 * 4,
				},
			},
		},

		// Input assembly
		PInputAssemblyState: &vk.PipelineInputAssemblyStateCreateInfo{
			SType:    vk.StructureTypePipelineInputAssemblyStateCreateInfo,
			Topology: vk.PrimitiveTopologyTriangleList,
		},

		// viewport state
		PViewportState: &vk.PipelineViewportStateCreateInfo{
			SType:         vk.StructureTypePipelineViewportStateCreateInfo,
			ViewportCount: 1,
			PViewports: []vk.Viewport{
				{
					Width:    1000,
					Height:   1000,
					MinDepth: 0,
					MaxDepth: 1,
				},
			},
			ScissorCount: 1,
			PScissors: []vk.Rect2D{
				// scissor
				{
					Offset: vk.Offset2D{},
					Extent: vk.Extent2D{
						Width:  1000,
						Height: 1000,
					},
				},
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
			LineWidth:               1,
		},

		// multisample
		PMultisampleState: &vk.PipelineMultisampleStateCreateInfo{
			SType:                vk.StructureTypePipelineMultisampleStateCreateInfo,
			RasterizationSamples: vk.SampleCountFlagBits(vk.SampleCount1Bit),
		},

		// depth & stencil
		PDepthStencilState: &vk.PipelineDepthStencilStateCreateInfo{
			// enable depth testing with less or
			SType:                 vk.StructureTypePipelineDepthStencilStateCreateInfo,
			DepthTestEnable:       vk.True,
			DepthWriteEnable:      vk.True,
			DepthCompareOp:        vk.CompareOpLessOrEqual,
			DepthBoundsTestEnable: vk.False,
			Back: vk.StencilOpState{
				FailOp:    vk.StencilOpKeep,
				PassOp:    vk.StencilOpKeep,
				CompareOp: vk.CompareOpAlways,
			},
			StencilTestEnable: vk.False,
			Front: vk.StencilOpState{
				FailOp:    vk.StencilOpKeep,
				PassOp:    vk.StencilOpKeep,
				CompareOp: vk.CompareOpAlways,
			},
		},

		// color blending
		PColorBlendState: &vk.PipelineColorBlendStateCreateInfo{
			SType:           vk.StructureTypePipelineColorBlendStateCreateInfo,
			LogicOpEnable:   vk.False,
			LogicOp:         vk.LogicOpClear,
			AttachmentCount: 1,
			PAttachments: []vk.PipelineColorBlendAttachmentState{
				{
					// additive blending
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
				vk.DynamicStateViewport,
				vk.DynamicStateScissor,
			},
		},
	}

	ptrs := make([]vk.Pipeline, 1)
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
