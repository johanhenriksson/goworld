package pipeline

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/vertex"
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

func New(device device.T, cache vk.PipelineCache, layout Layout, pass Pass, shaders []Shader, pointers vertex.Pointers) T {
	modules := util.Map(shaders, func(i int, shader Shader) vk.PipelineShaderStageCreateInfo {
		return vk.PipelineShaderStageCreateInfo{
			SType:  vk.StructureTypePipelineShaderStageCreateInfo,
			Module: shader.Ptr(),
			PName:  util.CString(shader.Entrypoint()),
			Stage:  shader.Stage(),
		}
	})

	attrs := pointersToVertexAttributes(pointers, 0)
	fmt.Println("attributes", attrs)

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
					Stride:    uint32(pointers.Stride()),
					InputRate: vk.VertexInputRateVertex,
				},
			},
			VertexAttributeDescriptionCount: uint32(len(attrs)),
			PVertexAttributeDescriptions:    attrs,
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

func pointersToVertexAttributes(ptrs vertex.Pointers, binding int) []vk.VertexInputAttributeDescription {
	attrs := make([]vk.VertexInputAttributeDescription, 0, len(ptrs))
	for _, ptr := range ptrs {
		if ptr.Binding < 0 {
			continue
		}
		attrs = append(attrs, vk.VertexInputAttributeDescription{
			Binding:  uint32(binding),
			Location: uint32(ptr.Binding),
			Format:   convertFormat(ptr),
			Offset:   uint32(ptr.Offset),
		})
	}
	return attrs
}

type ptrType struct {
	Source    types.Type
	Target    types.Type
	Elements  int
	Normalize bool
}

var formatMap = map[ptrType]vk.Format{
	{types.Float, types.Float, 1, false}:   vk.FormatR32Sfloat,
	{types.Float, types.Float, 2, false}:   vk.FormatR32g32Sfloat,
	{types.Float, types.Float, 3, false}:   vk.FormatR32g32b32Sfloat,
	{types.Float, types.Float, 4, false}:   vk.FormatR32g32b32a32Sfloat,
	{types.Int8, types.Int8, 1, false}:     vk.FormatR8Sint,
	{types.Int8, types.Int8, 2, false}:     vk.FormatR8g8Sint,
	{types.Int8, types.Int8, 3, false}:     vk.FormatR8g8b8Sint,
	{types.Int8, types.Int8, 4, false}:     vk.FormatR8g8b8a8Sint,
	{types.Int8, types.Float, 4, true}:     vk.FormatR8Snorm,
	{types.Int8, types.Float, 2, true}:     vk.FormatR8g8Snorm,
	{types.Int8, types.Float, 3, true}:     vk.FormatR8g8b8Snorm,
	{types.Int8, types.Float, 4, true}:     vk.FormatR8g8b8a8Snorm,
	{types.UInt8, types.UInt8, 1, false}:   vk.FormatR8Uint,
	{types.UInt8, types.UInt8, 2, false}:   vk.FormatR8g8Uint,
	{types.UInt8, types.UInt8, 3, false}:   vk.FormatR8g8b8Uint,
	{types.UInt8, types.UInt8, 4, false}:   vk.FormatR8g8b8a8Uint,
	{types.UInt8, types.Float, 1, false}:   vk.FormatR8Uscaled,
	{types.UInt8, types.Float, 2, false}:   vk.FormatR8g8Uscaled,
	{types.UInt8, types.Float, 3, false}:   vk.FormatR8g8b8Uscaled,
	{types.UInt8, types.Float, 4, false}:   vk.FormatR8g8b8a8Uscaled,
	{types.UInt8, types.Float, 1, true}:    vk.FormatR8Unorm,
	{types.UInt8, types.Float, 2, true}:    vk.FormatR8g8Unorm,
	{types.UInt8, types.Float, 3, true}:    vk.FormatR8g8b8Unorm,
	{types.UInt8, types.Float, 4, true}:    vk.FormatR8g8b8a8Unorm,
	{types.Int16, types.Int16, 1, false}:   vk.FormatR16Sint,
	{types.Int16, types.Int16, 2, false}:   vk.FormatR16g16Sint,
	{types.Int16, types.Int16, 3, false}:   vk.FormatR16g16b16Sint,
	{types.Int16, types.Int16, 4, false}:   vk.FormatR16g16b16a16Sint,
	{types.Int16, types.Float, 4, true}:    vk.FormatR16Snorm,
	{types.Int16, types.Float, 2, true}:    vk.FormatR16g16Snorm,
	{types.Int16, types.Float, 3, true}:    vk.FormatR16g16b16Snorm,
	{types.Int16, types.Float, 4, true}:    vk.FormatR16g16b16a16Snorm,
	{types.UInt16, types.UInt16, 1, false}: vk.FormatR16Uint,
	{types.UInt16, types.UInt16, 2, false}: vk.FormatR16g16Uint,
	{types.UInt16, types.UInt16, 3, false}: vk.FormatR16g16b16Uint,
	{types.UInt16, types.UInt16, 4, false}: vk.FormatR16g16b16a16Uint,
	{types.UInt16, types.Float, 1, true}:   vk.FormatR16Unorm,
	{types.UInt16, types.Float, 2, true}:   vk.FormatR16g16Unorm,
	{types.UInt16, types.Float, 3, true}:   vk.FormatR16g16b16Unorm,
	{types.UInt16, types.Float, 4, true}:   vk.FormatR16g16b16a16Unorm,
	{types.UInt16, types.Float, 1, false}:  vk.FormatR16Uscaled,
	{types.UInt16, types.Float, 2, false}:  vk.FormatR16g16Uscaled,
	{types.UInt16, types.Float, 3, false}:  vk.FormatR16g16b16Uscaled,
	{types.UInt16, types.Float, 4, false}:  vk.FormatR16g16b16a16Uscaled,
	{types.Int32, types.Int32, 1, false}:   vk.FormatR32Sint,
	{types.Int32, types.Int32, 2, false}:   vk.FormatR32g32Sint,
	{types.Int32, types.Int32, 3, false}:   vk.FormatR32g32b32Sint,
	{types.Int32, types.Int32, 4, false}:   vk.FormatR32g32b32a32Sint,
	{types.UInt32, types.UInt32, 1, false}: vk.FormatR32Uint,
	{types.UInt32, types.UInt32, 2, false}: vk.FormatR32g32Uint,
	{types.UInt32, types.UInt32, 3, false}: vk.FormatR32g32b32Uint,
	{types.UInt32, types.UInt32, 4, false}: vk.FormatR32g32b32a32Uint,
}

func convertFormat(ptr vertex.Pointer) vk.Format {
	kind := ptrType{ptr.Source, ptr.Destination, ptr.Elements, ptr.Normalize}
	if fmt, exists := formatMap[kind]; exists {
		return fmt
	}
	panic(fmt.Sprintf("illegal format in pointer %s from %s -> %s x%d (normalize: %t)", ptr.Name, ptr.Source, ptr.Destination, ptr.Elements, ptr.Normalize))
}
