package pipeline

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Pipeline]

	Layout() Layout
}

type pipeline struct {
	ptr    vk.Pipeline
	device device.T
	layout Layout
}

func New(device device.T, args Args) T {
	args.defaults()
	log.Println("creating pipeline")

	// todo: pipeline cache
	// could probably be controlled a global setting?

	modules := util.Map(args.Shader.Modules(), func(shader shader.Module) vk.PipelineShaderStageCreateInfo {
		return vk.PipelineShaderStageCreateInfo{
			SType:  vk.StructureTypePipelineShaderStageCreateInfo,
			Module: shader.Ptr(),
			PName:  util.CString(shader.Entrypoint()),
			Stage:  shader.Stage(),
		}
	})

	log.Println("  attributes", args.Pointers)
	attrs := pointersToVertexAttributes(args.Pointers, 0)

	subpass := args.Pass.Subpass(args.Subpass)
	log.Println("  subpass:", subpass.Name, subpass.Index())

	blendStates := util.Map(subpass.ColorAttachments, func(name string) vk.PipelineColorBlendAttachmentState {
		attach := args.Pass.Attachment(name)
		// todo: move into attachment object
		// or into the material/pipeline object?
		blend := attach.Blend()
		return vk.PipelineColorBlendAttachmentState{
			// additive blending
			BlendEnable:         vkBool(blend.Enabled),
			ColorBlendOp:        blend.Color.Operation,
			SrcColorBlendFactor: blend.Color.SrcFactor,
			DstColorBlendFactor: blend.Color.DstFactor,
			AlphaBlendOp:        blend.Alpha.Operation,
			SrcAlphaBlendFactor: blend.Alpha.SrcFactor,
			DstAlphaBlendFactor: blend.Alpha.DstFactor,
			ColorWriteMask: vk.ColorComponentFlags(
				vk.ColorComponentRBit | vk.ColorComponentGBit |
					vk.ColorComponentBBit | vk.ColorComponentABit),
		}
	})

	info := vk.GraphicsPipelineCreateInfo{
		SType: vk.StructureTypeGraphicsPipelineCreateInfo,

		// layout
		Layout:  args.Layout.Ptr(),
		Subpass: uint32(subpass.Index()),

		// render pass
		RenderPass: args.Pass.Ptr(),

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
					Stride:    uint32(args.Pointers.Stride()),
					InputRate: vk.VertexInputRateVertex,
				},
			},
			VertexAttributeDescriptionCount: uint32(len(attrs)),
			PVertexAttributeDescriptions:    attrs,
		},

		// Input assembly
		PInputAssemblyState: &vk.PipelineInputAssemblyStateCreateInfo{
			SType:    vk.StructureTypePipelineInputAssemblyStateCreateInfo,
			Topology: vkPrimitiveTopology(args.Primitive),
		},

		// viewport state
		// does not seem to matter so much since we set it dynamically every frame
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
			PolygonMode:             args.PolygonFillMode,
			CullMode:                vk.CullModeFlags(args.CullMode),
			FrontFace:               vk.FrontFaceCounterClockwise,
			LineWidth:               1,
		},

		// multisample
		PMultisampleState: &vk.PipelineMultisampleStateCreateInfo{
			SType:                vk.StructureTypePipelineMultisampleStateCreateInfo,
			RasterizationSamples: vk.SampleCount1Bit,
		},

		// depth & stencil
		PDepthStencilState: &vk.PipelineDepthStencilStateCreateInfo{
			// enable depth testing with less or
			SType:                 vk.StructureTypePipelineDepthStencilStateCreateInfo,
			DepthTestEnable:       vkBool(args.DepthTest),
			DepthWriteEnable:      vkBool(args.DepthWrite),
			DepthCompareOp:        args.DepthFunc,
			DepthBoundsTestEnable: vk.False,
			Back: vk.StencilOpState{
				FailOp:    vk.StencilOpKeep,
				PassOp:    vk.StencilOpKeep,
				CompareOp: vk.CompareOpAlways,
			},
			StencilTestEnable: vkBool(args.StencilTest),
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
			AttachmentCount: uint32(len(blendStates)),
			PAttachments:    blendStates,
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
	result := vk.CreateGraphicsPipelines(device.Ptr(), nil, 1, []vk.GraphicsPipelineCreateInfo{info}, nil, ptrs)
	if result != vk.Success {
		panic("failed to create pipeline")
	}

	return &pipeline{
		ptr:    ptrs[0],
		device: device,
		layout: args.Layout,
	}
}

func (p *pipeline) Ptr() vk.Pipeline {
	return p.ptr
}

func (p *pipeline) Layout() Layout {
	return p.layout
}

func (p *pipeline) Destroy() {
	if p.ptr != nil {
		vk.DestroyPipeline(p.device.Ptr(), p.ptr, nil)
		p.ptr = nil
	}
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

func vkBool(v bool) vk.Bool32 {
	if v {
		return vk.True
	}
	return vk.False
}

func vkPrimitiveTopology(primitive vertex.Primitive) vk.PrimitiveTopology {
	switch primitive {
	case vertex.Triangles:
		return vk.PrimitiveTopologyTriangleList
	case vertex.Lines:
		return vk.PrimitiveTopologyLineList
	case vertex.Points:
		return vk.PrimitiveTopologyPointList
	}
	panic("unknown primitive")
}
