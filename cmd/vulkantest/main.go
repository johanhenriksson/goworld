package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/vertex"

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
	defer func() {
		log.Println("Clean exit")
	}()

	running := true
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		for range sigint {
			log.Println("Interrupt")
			running = false
		}
	}()

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

	proj := mat4.PerspectiveVK(45, 1, 0.1, 100)
	view := mat4.Translate(vec3.New(0, 0, -6))

	ubo := buffer.NewUniform(backend.Device(), 3*16*4)
	defer ubo.Destroy()
	ubo.Write([]mat4.T{
		proj,
		view,
		mat4.Rotate(vec3.New(0, 0, 45)),
	}, 0)

	vtx := buffer.NewRemote(backend.Device(), 4096, vk.BufferUsageFlags(vk.BufferUsageVertexBufferBit))
	defer vtx.Destroy()
	vtxstage := buffer.NewShared(backend.Device(), 4096)
	vtxstage.Write(makeColorCube(1), 0)

	transferer := command.NewWorker(backend.Device(), vk.QueueFlags(vk.QueueTransferBit))
	defer transferer.Destroy()
	transferer.Queue(func(cmd command.Buffer) {
		cmd.CmdCopyBuffer(vtxstage, vtx)
	})
	transferer.Submit(command.SubmitInfo{})
	log.Println("waiting for transfers...")
	transferer.Wait()
	log.Println("transfers completed")

	vtxstage.Destroy()

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

	t := 0.0
	for running && !wnd.ShouldClose() {
		ctx := backend.Aquire()

		// update ubo
		t += 0.016
		ubo.Write([]mat4.T{
			mat4.Perspective(45, float32(ctx.Width)/float32(ctx.Height), 0.1, 100),
			view,
			mat4.Rotate(vec3.New(0, float32(44*t), float32(90*t))),
		}, 0)

		worker := ctx.Workers[0]

		worker.Queue(func(cmd command.Buffer) {
			clear := color.RGB(0.2, 0.2, 0.2)

			cmd.CmdBeginRenderPass(backend.Swapchain().Output(), ctx.Framebuffer, clear)
			cmd.CmdSetViewport(0, 0, ctx.Width, ctx.Height)
			cmd.CmdSetScissor(0, 0, ctx.Width, ctx.Height)

			// user draw calls
			cmd.CmdBindGraphicsPipeline(pipe)
			cmd.CmdBindGraphicsDescriptors(layout, desc)

			cmd.CmdBindVertexBuffer(vtx, 0)
			cmd.CmdDraw(36, 12, 0, 0)

			cmd.CmdEndRenderPass()
		})

		worker.Submit(command.SubmitInfo{
			Wait:   []sync.Semaphore{ctx.ImageAvailable},
			Signal: []sync.Semaphore{ctx.RenderComplete},
			WaitMask: []vk.PipelineStageFlags{
				vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
			},
		})

		backend.Present()

		wnd.Poll()
	}

	vk.DeviceWaitIdle(backend.Device().Ptr())
}

func makeColorCube(s float32) []vertex.C {
	return []vertex.C{
		{P: vec3.New(s, -s, s), N: vec3.UnitX, C: vec4.New(1, 0, 1, 1)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitX, C: vec4.New(1, 0, 0, 1)},
		{P: vec3.New(s, s, -s), N: vec3.UnitX, C: vec4.New(1, 1, 0, 1)},
		{P: vec3.New(s, -s, s), N: vec3.UnitX, C: vec4.New(1, 0, 1, 1)},
		{P: vec3.New(s, s, -s), N: vec3.UnitX, C: vec4.New(1, 1, 0, 1)},
		{P: vec3.New(s, s, s), N: vec3.UnitX, C: vec4.New(1, 1, 1, 1)},

		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, C: vec4.New(0, 0, 1, 1)},
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, C: vec4.New(0, 1, 0, 1)},
		{P: vec3.New(-s, -s, -s), N: vec3.UnitXN, C: vec4.New(0, 0, 0, 1)},
		{P: vec3.New(-s, -s, s), N: vec3.UnitXN, C: vec4.New(0, 0, 1, 1)},
		{P: vec3.New(-s, s, s), N: vec3.UnitXN, C: vec4.New(0, 1, 1, 1)},
		{P: vec3.New(-s, s, -s), N: vec3.UnitXN, C: vec4.New(0, 1, 0, 1)},

		{P: vec3.New(-s, s, -s), N: vec3.UnitY, C: vec4.New(0, 1, 0, 1)},
		{P: vec3.New(-s, s, s), N: vec3.UnitY, C: vec4.New(0, 1, 1, 1)},
		{P: vec3.New(s, s, -s), N: vec3.UnitY, C: vec4.New(1, 1, 0, 1)},
		{P: vec3.New(s, s, -s), N: vec3.UnitY, C: vec4.New(1, 1, 0, 1)},
		{P: vec3.New(-s, s, s), N: vec3.UnitY, C: vec4.New(0, 1, 1, 1)},
		{P: vec3.New(s, s, s), N: vec3.UnitY, C: vec4.New(1, 1, 1, 1)},

		{P: vec3.New(-s, -s, -s), N: vec3.UnitYN, C: vec4.New(0, 0, 0, 1)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, C: vec4.New(1, 0, 0, 1)},
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, C: vec4.New(0, 0, 1, 1)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitYN, C: vec4.New(1, 0, 0, 1)},
		{P: vec3.New(s, -s, s), N: vec3.UnitYN, C: vec4.New(1, 0, 1, 1)},
		{P: vec3.New(-s, -s, s), N: vec3.UnitYN, C: vec4.New(0, 0, 1, 1)},

		{P: vec3.New(-s, -s, s), N: vec3.UnitZ, C: vec4.New(0, 0, 1, 1)},
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, C: vec4.New(1, 0, 1, 1)},
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, C: vec4.New(0, 1, 1, 1)},
		{P: vec3.New(s, -s, s), N: vec3.UnitZ, C: vec4.New(1, 0, 1, 1)},
		{P: vec3.New(s, s, s), N: vec3.UnitZ, C: vec4.New(1, 1, 1, 1)},
		{P: vec3.New(-s, s, s), N: vec3.UnitZ, C: vec4.New(0, 1, 1, 1)},

		{P: vec3.New(-s, -s, -s), N: vec3.UnitZN, C: vec4.New(0, 0, 0, 1)},
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, C: vec4.New(0, 1, 0, 1)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, C: vec4.New(1, 0, 0, 1)},
		{P: vec3.New(s, -s, -s), N: vec3.UnitZN, C: vec4.New(1, 0, 0, 1)},
		{P: vec3.New(-s, s, -s), N: vec3.UnitZN, C: vec4.New(0, 1, 0, 1)},
		{P: vec3.New(s, s, -s), N: vec3.UnitZN, C: vec4.New(1, 1, 0, 1)},
	}
}
