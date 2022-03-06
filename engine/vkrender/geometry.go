package vkrender

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/color"

	vk "github.com/vulkan-go/vulkan"
)

type GeometryPass struct {
	meshes  cache.Meshes
	backend vulkan.T
	ubo     []buffer.T
	layout  pipeline.Layout
	pipe    pipeline.T
	shaders []pipeline.Shader
	dlayout descriptor.T
	dpool   descriptor.Pool
	dsets   []descriptor.Set
}

func NewGeometryPass(backend vulkan.T, meshes cache.Meshes) *GeometryPass {
	p := &GeometryPass{
		backend: backend,
		meshes:  meshes,
	}

	ubosize := 3 * 16 * 4
	p.ubo = []buffer.T{
		buffer.NewUniform(p.backend.Device(), ubosize),
		buffer.NewUniform(p.backend.Device(), ubosize),
	}
	p.ubo[0].Write([]mat4.T{
		mat4.Ident(),
		mat4.Ident(),
		mat4.Rotate(vec3.New(0, 0, 45)),
	}, 0)
	p.ubo[1].Write([]mat4.T{
		mat4.Ident(),
		mat4.Ident(),
		mat4.Rotate(vec3.New(0, 0, 45)),
	}, 0)

	var cache vk.PipelineCache
	vk.CreatePipelineCache(p.backend.Device().Ptr(), &vk.PipelineCacheCreateInfo{
		SType: vk.StructureTypePipelineCacheCreateInfo,
	}, nil, &cache)
	defer vk.DestroyPipelineCache(p.backend.Device().Ptr(), cache, nil)

	p.shaders = []pipeline.Shader{
		pipeline.NewShader(p.backend.Device(), "assets/shaders/vk/color_f.vert.spv", vk.ShaderStageVertexBit),
		pipeline.NewShader(p.backend.Device(), "assets/shaders/vk/color_f.frag.spv", vk.ShaderStageFragmentBit),
	}

	p.dlayout = descriptor.New(p.backend.Device(), []descriptor.Binding{
		{
			Binding: 0,
			Type:    vk.DescriptorTypeUniformBuffer,
			Count:   1,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageVertexBit),
		},
	})
	dlayouts := []descriptor.T{p.dlayout, p.dlayout}

	p.dpool = descriptor.NewPool(p.backend.Device(), []vk.DescriptorPoolSize{
		{
			Type:            vk.DescriptorTypeUniformBuffer,
			DescriptorCount: 2,
		},
	})

	p.dsets = p.dpool.AllocateSets(dlayouts)

	vk.UpdateDescriptorSets(p.backend.Device().Ptr(), 2, []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          p.dsets[0].Ptr(),
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeUniformBuffer,
			PBufferInfo: []vk.DescriptorBufferInfo{
				{
					Buffer: p.ubo[0].Ptr(),
					Offset: 0,
					Range:  vk.DeviceSize(ubosize),
				},
			},
		},
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          p.dsets[1].Ptr(),
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeUniformBuffer,
			PBufferInfo: []vk.DescriptorBufferInfo{
				{
					Buffer: p.ubo[1].Ptr(),
					Offset: 0,
					Range:  vk.DeviceSize(ubosize),
				},
			},
		},
	}, 0, nil)

	p.layout = pipeline.NewLayout(p.backend.Device(), dlayouts)

	p.pipe = pipeline.New(p.backend.Device(), cache, p.layout, p.backend.Swapchain().Output(), p.shaders)

	return p
}

func (p *GeometryPass) Draw(args render.Args, scene object.T) {
	objects := query.New[mesh.T]().Where(isDrawDeferred).Collect(scene)
	for _, mesh := range objects {
		if err := p.DrawDeferred(args, mesh); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", mesh.Name(), err)
		}
	}
}

func (p *GeometryPass) DrawDeferred(args render.Args, mesh mesh.T) error {
	args = args.Apply(mesh.Transform().World())
	ctx := args.Context

	p.ubo[ctx.Index].Write([]mat4.T{
		args.Projection,
		args.View,
		mat4.Ident(),
	}, 0)

	vkmesh, ok := p.meshes.Fetch(mesh.Mesh(), nil).(*cache.VkMesh)
	if !ok {
		fmt.Println("mesh is nil")
		return nil
	}

	worker := ctx.Workers[0]

	worker.Queue(func(cmd command.Buffer) {
		clear := color.RGB(0.2, 0.2, 0.2)

		cmd.CmdBeginRenderPass(p.backend.Swapchain().Output(), ctx.Framebuffer, clear)
		cmd.CmdSetViewport(0, 0, ctx.Width, ctx.Height)
		cmd.CmdSetScissor(0, 0, ctx.Width, ctx.Height)

		// user draw calls
		cmd.CmdBindGraphicsPipeline(p.pipe)
		cmd.CmdBindGraphicsDescriptors(p.layout, p.dsets[ctx.Index:ctx.Index+1])

		cmd.CmdBindVertexBuffer(vkmesh.Buffer, 0)
		cmd.CmdDraw(vkmesh.Mesh.Elements(), vkmesh.Mesh.Elements()/3, 0, 0)

		cmd.CmdEndRenderPass()
	})

	worker.Submit(command.SubmitInfo{
		Wait:   []sync.Semaphore{ctx.ImageAvailable},
		Signal: []sync.Semaphore{ctx.RenderComplete},
		WaitMask: []vk.PipelineStageFlags{
			vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		},
	})

	return nil
}

func (p *GeometryPass) Destroy() {
	p.pipe.Destroy()
	p.layout.Destroy()
	p.dpool.Destroy()
	p.dlayout.Destroy()

	for _, shader := range p.shaders {
		shader.Destroy()
	}
	for _, ubo := range p.ubo {
		ubo.Destroy()
	}
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}
