package pipeline

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/types"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/samber/lo"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Pipeline struct {
	ptr    core1_0.Pipeline
	device *device.Device
	args   Args
}

func New(device *device.Device, args Args) *Pipeline {
	if device == nil {
		panic("device is nil")
	}
	if args.Shader == nil {
		panic("shader is nil")
	}

	args.defaults()

	key := fmt.Sprintf("%s/%s", args.Pass.Name(), args.Shader.Name())
	log.Println("creating pipeline", key)

	// todo: pipeline cache

	for i, ptr := range args.Pointers {
		if index, kind, exists := args.Shader.Input(ptr.Name); exists {
			ptr.Bind(index, kind)
			args.Pointers[i] = ptr
		} else {
			log.Printf("no attribute in shader %s\n", ptr.Name)
		}
	}

	modules := lo.Map(args.Shader.Modules(), func(shader shader.Module, _ int) core1_0.PipelineShaderStageCreateInfo {
		return core1_0.PipelineShaderStageCreateInfo{
			Module: shader.Ptr(),
			Name:   shader.Entrypoint(),
			Stage:  core1_0.ShaderStageFlags(shader.Stage()),
		}
	})

	log.Println("  depth test:", args.DepthTest)
	log.Println("  depth func:", args.DepthFunc)
	log.Println("  depth write:", args.DepthWrite)

	log.Println("  attributes", args.Pointers)
	vertexInput := pointersToVertexInput(args.Pointers, 0)

	subpass := args.Pass.Subpass(args.Subpass)
	log.Println("  subpass:", subpass.Name, subpass.Index())

	blendStates := lo.Map(subpass.ColorAttachments, func(name attachment.Name, _ int) core1_0.PipelineColorBlendAttachmentState {
		attach := args.Pass.Attachment(name)
		// todo: move into attachment object
		// or into the material/pipeline object?
		blend := attach.Blend()
		return core1_0.PipelineColorBlendAttachmentState{
			// additive blending
			BlendEnabled:        blend.Enabled,
			ColorBlendOp:        blend.Color.Operation,
			SrcColorBlendFactor: blend.Color.SrcFactor,
			DstColorBlendFactor: blend.Color.DstFactor,
			AlphaBlendOp:        blend.Alpha.Operation,
			SrcAlphaBlendFactor: blend.Alpha.SrcFactor,
			DstAlphaBlendFactor: blend.Alpha.DstFactor,
			ColorWriteMask: core1_0.ColorComponentRed | core1_0.ColorComponentGreen |
				core1_0.ColorComponentBlue | core1_0.ColorComponentAlpha,
		}
	})

	info := core1_0.GraphicsPipelineCreateInfo{
		// layout
		Layout:  args.Layout.Ptr(),
		Subpass: subpass.Index(),

		// render pass
		RenderPass: args.Pass.Ptr(),

		// Stages
		Stages: modules,

		// Vertex input state
		VertexInputState: vertexInput,

		// Input assembly
		InputAssemblyState: &core1_0.PipelineInputAssemblyStateCreateInfo{
			Topology: core1_0.PrimitiveTopology(args.Primitive),
		},

		// viewport state
		// does not seem to matter so much since we set it dynamically every frame
		ViewportState: &core1_0.PipelineViewportStateCreateInfo{
			Viewports: []core1_0.Viewport{
				{
					Width:    1000,
					Height:   1000,
					MinDepth: 0,
					MaxDepth: 1,
				},
			},
			Scissors: []core1_0.Rect2D{
				// scissor
				{
					Offset: core1_0.Offset2D{},
					Extent: core1_0.Extent2D{
						Width:  1000,
						Height: 1000,
					},
				},
			},
		},

		// rasterization state
		RasterizationState: &core1_0.PipelineRasterizationStateCreateInfo{
			DepthClampEnable:        args.DepthClamp,
			DepthBiasEnable:         false,
			RasterizerDiscardEnable: false,
			PolygonMode:             args.PolygonFillMode,
			CullMode:                core1_0.CullModeFlags(args.CullMode),
			LineWidth:               1,

			// clockwise in vulkans right-handed coordinates is equivalent to the
			// traditional opengl counter-clockwise winding, which is in line with
			// the left-handed world space coordinate system.
			FrontFace: core1_0.FrontFaceClockwise,
		},

		// multisample
		MultisampleState: &core1_0.PipelineMultisampleStateCreateInfo{
			RasterizationSamples: core1_0.Samples1,
		},

		// depth & stencil
		DepthStencilState: &core1_0.PipelineDepthStencilStateCreateInfo{
			// enable depth testing with less or
			DepthTestEnable:       args.DepthTest,
			DepthWriteEnable:      args.DepthWrite,
			DepthCompareOp:        args.DepthFunc,
			DepthBoundsTestEnable: false,
			Back: core1_0.StencilOpState{
				FailOp:    core1_0.StencilKeep,
				PassOp:    core1_0.StencilKeep,
				CompareOp: core1_0.CompareOpAlways,
			},
			StencilTestEnable: args.StencilTest,
			Front: core1_0.StencilOpState{
				FailOp:    core1_0.StencilKeep,
				PassOp:    core1_0.StencilKeep,
				CompareOp: core1_0.CompareOpAlways,
			},
		},

		// color blending
		ColorBlendState: &core1_0.PipelineColorBlendStateCreateInfo{
			LogicOpEnabled: false,
			LogicOp:        core1_0.LogicOpClear,
			Attachments:    blendStates,
		},

		// dynamic state: viewport & scissor
		DynamicState: &core1_0.PipelineDynamicStateCreateInfo{
			DynamicStates: []core1_0.DynamicState{
				core1_0.DynamicStateViewport,
				core1_0.DynamicStateScissor,
			},
		},
	}

	ptrs, result, err := device.Ptr().CreateGraphicsPipelines(nil, nil, []core1_0.GraphicsPipelineCreateInfo{info})
	if err != nil {
		panic(err)
	}
	if result != core1_0.VKSuccess {
		panic("failed to create pipeline")
	}

	if key != "" {
		device.SetDebugObjectName(driver.VulkanHandle(ptrs[0].Handle()), core1_0.ObjectTypePipeline, key)
	}

	return &Pipeline{
		ptr:    ptrs[0],
		device: device,
		args:   args,
	}
}

func (p *Pipeline) Ptr() core1_0.Pipeline {
	return p.ptr
}

func (p *Pipeline) Shader() *shader.Shader {
	return p.args.Shader
}

func (p *Pipeline) Layout() *Layout {
	return p.args.Layout
}

func (p *Pipeline) Destroy() {
	if p.ptr != nil {
		p.ptr.Destroy(nil)
		p.ptr = nil
	}
}

func pointersToVertexInput(ptrs vertex.Pointers, binding int) *core1_0.PipelineVertexInputStateCreateInfo {
	bindings := make([]core1_0.VertexInputBindingDescription, 0, 1)
	attrs := make([]core1_0.VertexInputAttributeDescription, 0, len(ptrs))
	for _, ptr := range ptrs {
		if ptr.Binding < 0 {
			continue
		}
		attrs = append(attrs, core1_0.VertexInputAttributeDescription{
			Binding:  binding,
			Location: uint32(ptr.Binding),
			Format:   convertFormat(ptr),
			Offset:   ptr.Offset,
		})
	}
	if len(attrs) > 0 {
		bindings = append(bindings, core1_0.VertexInputBindingDescription{
			Binding:   0,
			Stride:    ptrs.Stride(),
			InputRate: core1_0.VertexInputRateVertex,
		})
	}
	return &core1_0.PipelineVertexInputStateCreateInfo{
		VertexBindingDescriptions:   bindings,
		VertexAttributeDescriptions: attrs,
	}
}

type ptrType struct {
	Source    types.Type
	Target    types.Type
	Elements  int
	Normalize bool
}

var formatMap = map[ptrType]core1_0.Format{
	{types.Float, types.Float, 1, false}:   core1_0.FormatR32SignedFloat,
	{types.Float, types.Float, 2, false}:   core1_0.FormatR32G32SignedFloat,
	{types.Float, types.Float, 3, false}:   core1_0.FormatR32G32B32SignedFloat,
	{types.Float, types.Float, 4, false}:   core1_0.FormatR32G32B32A32SignedFloat,
	{types.Int8, types.Int8, 1, false}:     core1_0.FormatR8SignedInt,
	{types.Int8, types.Int8, 2, false}:     core1_0.FormatR8G8SignedInt,
	{types.Int8, types.Int8, 3, false}:     core1_0.FormatR8G8B8SignedInt,
	{types.Int8, types.Int8, 4, false}:     core1_0.FormatR8G8B8A8SignedInt,
	{types.Int8, types.Float, 4, true}:     core1_0.FormatR8SignedNormalized,
	{types.Int8, types.Float, 2, true}:     core1_0.FormatR8G8SignedNormalized,
	{types.Int8, types.Float, 3, true}:     core1_0.FormatR8G8B8SignedNormalized,
	{types.Int8, types.Float, 4, true}:     core1_0.FormatR8G8B8A8SignedNormalized,
	{types.UInt8, types.UInt8, 1, false}:   core1_0.FormatR8UnsignedInt,
	{types.UInt8, types.UInt8, 2, false}:   core1_0.FormatR8G8UnsignedInt,
	{types.UInt8, types.UInt8, 3, false}:   core1_0.FormatR8G8B8UnsignedInt,
	{types.UInt8, types.UInt8, 4, false}:   core1_0.FormatR8G8B8A8UnsignedInt,
	{types.UInt8, types.Float, 1, false}:   core1_0.FormatR8UnsignedScaled,
	{types.UInt8, types.Float, 2, false}:   core1_0.FormatR8G8UnsignedScaled,
	{types.UInt8, types.Float, 3, false}:   core1_0.FormatR8G8B8UnsignedScaled,
	{types.UInt8, types.Float, 4, false}:   core1_0.FormatR8G8B8A8UnsignedScaled,
	{types.UInt8, types.Float, 1, true}:    core1_0.FormatR8UnsignedNormalized,
	{types.UInt8, types.Float, 2, true}:    core1_0.FormatR8G8UnsignedNormalized,
	{types.UInt8, types.Float, 3, true}:    core1_0.FormatR8G8B8UnsignedNormalized,
	{types.UInt8, types.Float, 4, true}:    core1_0.FormatR8G8B8A8UnsignedNormalized,
	{types.Int16, types.Int16, 1, false}:   core1_0.FormatR16SignedInt,
	{types.Int16, types.Int16, 2, false}:   core1_0.FormatR16G16SignedInt,
	{types.Int16, types.Int16, 3, false}:   core1_0.FormatR16G16B16SignedInt,
	{types.Int16, types.Int16, 4, false}:   core1_0.FormatR16G16B16A16SignedInt,
	{types.Int16, types.Float, 4, true}:    core1_0.FormatR16SignedNormalized,
	{types.Int16, types.Float, 2, true}:    core1_0.FormatR16G16SignedNormalized,
	{types.Int16, types.Float, 3, true}:    core1_0.FormatR16G16B16SignedNormalized,
	{types.Int16, types.Float, 4, true}:    core1_0.FormatR16G16B16A16SignedNormalized,
	{types.UInt16, types.UInt16, 1, false}: core1_0.FormatR16UnsignedInt,
	{types.UInt16, types.UInt16, 2, false}: core1_0.FormatR16G16UnsignedInt,
	{types.UInt16, types.UInt16, 3, false}: core1_0.FormatR16G16B16UnsignedInt,
	{types.UInt16, types.UInt16, 4, false}: core1_0.FormatR16G16B16A16UnsignedInt,
	{types.UInt16, types.Float, 1, true}:   core1_0.FormatR16UnsignedNormalized,
	{types.UInt16, types.Float, 2, true}:   core1_0.FormatR16G16UnsignedNormalized,
	{types.UInt16, types.Float, 3, true}:   core1_0.FormatR16G16B16UnsignedNormalized,
	{types.UInt16, types.Float, 4, true}:   core1_0.FormatR16G16B16A16UnsignedNormalized,
	{types.UInt16, types.Float, 1, false}:  core1_0.FormatR16UnsignedScaled,
	{types.UInt16, types.Float, 2, false}:  core1_0.FormatR16G16UnsignedScaled,
	{types.UInt16, types.Float, 3, false}:  core1_0.FormatR16G16B16UnsignedScaled,
	{types.UInt16, types.Float, 4, false}:  core1_0.FormatR16G16B16A16UnsignedScaled,
	{types.Int32, types.Int32, 1, false}:   core1_0.FormatR32SignedInt,
	{types.Int32, types.Int32, 2, false}:   core1_0.FormatR32G32SignedInt,
	{types.Int32, types.Int32, 3, false}:   core1_0.FormatR32G32B32SignedInt,
	{types.Int32, types.Int32, 4, false}:   core1_0.FormatR32G32B32A32SignedInt,
	{types.UInt32, types.UInt32, 1, false}: core1_0.FormatR32UnsignedInt,
	{types.UInt32, types.UInt32, 2, false}: core1_0.FormatR32G32UnsignedInt,
	{types.UInt32, types.UInt32, 3, false}: core1_0.FormatR32G32B32UnsignedInt,
	{types.UInt32, types.UInt32, 4, false}: core1_0.FormatR32G32B32A32UnsignedInt,
}

func convertFormat(ptr vertex.Pointer) core1_0.Format {
	kind := ptrType{ptr.Source, ptr.Destination, ptr.Elements, ptr.Normalize}
	if fmt, exists := formatMap[kind]; exists {
		return fmt
	}
	panic(fmt.Sprintf("illegal format in pointer %s from %s -> %s x%d (normalize: %t)", ptr.Name, ptr.Source, ptr.Destination, ptr.Elements, ptr.Normalize))
}
