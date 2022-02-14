package main

import (
	"log"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"

	vk "github.com/vulkan-go/vulkan"
)

type VkVertex struct {
	X, Y, Z float32
	R, G, B float32
}

type Uniforms struct {
	Proj  mat4.T
	View  mat4.T
	Model mat4.T
}

func main() {
	backend := vulkan.New("goworld: vulkan", 0)
	defer backend.Destroy()

	wnd, err := window.New(backend, window.Args{
		Title:  "goworld: vulkan",
		Width:  500,
		Height: 500,
	})
	if err != nil {
		panic(err)
	}
	queue := backend.Device().GetQueue(0, vk.QueueFlags(vk.QueueGraphicsBit))

	proj := mat4.Perspective(50, 1, 0.1, 100)
	proj[5] *= -1

	ubo := buffer.NewRemote(backend.Device(), 3*16*4, vk.BufferUsageFlags(vk.BufferUsageUniformBufferBit))
	defer ubo.Destroy()
	ubostage := buffer.NewShared(backend.Device(), 3*16*4)
	ubostage.Write([]Uniforms{
		{
			Proj:  proj,
			View:  mat4.Translate(vec3.New(0, 0, -15)),
			Model: mat4.Ident(),
		},
	}, 0)

	vtx := buffer.NewRemote(backend.Device(), 3*24, vk.BufferUsageFlags(vk.BufferUsageVertexBufferBit))
	defer vtx.Destroy()
	vtxstage := buffer.NewShared(backend.Device(), 3*24)
	vtxstage.Write([]VkVertex{
		{1.0, 1.0, 0.0, 1.0, 0.0, 0.0},
		{-1.0, 1.0, 0.0, 0.0, 1.0, 0.0},
		{0.0, -1.0, 0.0, 0.0, 0.0, 1.0},
	}, 0)

	idx := buffer.NewRemote(backend.Device(), 3*4, vk.BufferUsageFlags(vk.BufferUsageIndexBufferBit))
	defer idx.Destroy()
	idxstage := buffer.NewShared(backend.Device(), 3*4)
	idxstage.Write([]uint32{0, 1, 2}, 0)

	transferer := command.NewWorker(backend.Device())
	transferer.Queue(func(cmd command.Buffer) {
		cmd.CmdCopyBuffer(vtxstage, vtx)
		cmd.CmdCopyBuffer(idxstage, idx)
		cmd.CmdCopyBuffer(ubostage, ubo)
	})
	transferer.Submit(command.SubmitInfo{
		Queue: queue,
	})
	log.Println("waiting for transfers...")
	transferer.Wait()
	log.Println("transfers completed")

	vtxstage.Destroy()
	idxstage.Destroy()
	ubostage.Destroy()

	var cache vk.PipelineCache
	vk.CreatePipelineCache(backend.Device().Ptr(), &vk.PipelineCacheCreateInfo{
		SType: vk.StructureTypePipelineCacheCreateInfo,
	}, nil, &cache)
	defer vk.DestroyPipelineCache(backend.Device().Ptr(), cache, nil)

	shaders := []pipeline.Shader{
		pipeline.NewShader(backend.Device(), "assets/shaders/vk/color_f.vert.spv", vk.ShaderStageVertexBit),
		pipeline.NewShader(backend.Device(), "assets/shaders/vk/color_f.frag.spv", vk.ShaderStageFragmentBit),
	}
	defer shaders[0].Destroy()
	defer shaders[1].Destroy()

	ubolayout := descriptor.New(backend.Device(), []descriptor.Binding{
		{
			Binding: 0,
			Type:    vk.DescriptorTypeUniformBuffer,
			Count:   1,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageVertexBit),
		},
	})
	defer ubolayout.Destroy()
	dlayouts := []descriptor.T{ubolayout}

	dpool := descriptor.NewPool(backend.Device(), []vk.DescriptorPoolSize{
		{
			Type:            vk.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
		},
	})
	defer dpool.Destroy()

	desc := dpool.AllocateSets(dlayouts)
	// uboset := desc[0]

	vk.UpdateDescriptorSets(backend.Device().Ptr(), 1, []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          desc[0].Ptr(),
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeUniformBuffer,
			PBufferInfo: []vk.DescriptorBufferInfo{
				{
					Buffer: ubo.Ptr(),
					Offset: 0,
					Range:  vk.DeviceSize(ubo.Size()),
				},
			},
		},
	}, 0, nil)

	layout := pipeline.NewLayout(backend.Device(), dlayouts)
	defer layout.Destroy()

	pipe := pipeline.New(backend.Device(), cache, layout, backend.Swapchain().Output(), shaders)
	defer pipe.Destroy()

	f := 0
	for !wnd.ShouldClose() {
		// update ubo

		backend.Aquire()

		backend.Present(func(buf command.Buffer) {
			buf.CmdBindGraphicsPipeline(pipe)
			buf.CmdBindGraphicsDescriptors(layout, desc)

			buf.CmdBindVertexBuffer(vtx, 0)
			buf.CmdBindIndexBuffers(idx, 0, vk.IndexTypeUint32)

			buf.CmdDrawIndexed(3, 1, 0, 0, 0)
		})

		wnd.Poll()
		f++

	}

	vk.DeviceWaitIdle(backend.Device().Ptr())
}
