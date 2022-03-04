package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/game"
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

var cpuprof = flag.String("cpuprof", "", "write cpu profile to file")

func main() {
	defer func() {
		log.Println("Clean exit")
	}()

	flag.Parse()
	if *cpuprof != "" {
		os.MkdirAll("profiling", 0755)
		ppath := fmt.Sprintf("profiling/%s", *cpuprof)
		f, err := os.Create(ppath)
		if err != nil {
			panic(err)
		}
		fmt.Println("writing cpu profiling output to", ppath)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	running := true
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		for range sigint {
			if !running {
				log.Println("Kill")
				os.Exit(1)
			} else {
				log.Println("Interrupt")
				running = false
			}
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
	// view := mat4.LookAtVK(vec3.New(-5, 8, -5), vec3.New(0, 0, 0))
	view := mat4.LookAtLH(vec3.New(-11, 17, -11), vec3.New(8, 3, 8))
	// view.Invert()

	ubosize := 3 * 16 * 4
	ubo := []buffer.T{
		buffer.NewUniform(backend.Device(), ubosize),
		buffer.NewUniform(backend.Device(), ubosize),
	}
	defer ubo[0].Destroy()
	defer ubo[1].Destroy()
	ubo[0].Write([]mat4.T{
		proj,
		view,
		mat4.Rotate(vec3.New(0, 0, 45)),
	}, 0)
	ubo[1].Write([]mat4.T{
		proj,
		view,
		mat4.Rotate(vec3.New(0, 0, 45)),
	}, 0)

	world := game.NewWorld(31481234, 16)
	chunk := world.AddChunk(0, 0)
	vertexdata := game.ComputeVertexData(chunk)

	vertices := len(vertexdata)
	triangles := vertices / 3
	log.Println("vertices:", vertices, "triangles:", triangles)
	bufsize := vertices * 10 * 4

	vtx := buffer.NewRemote(backend.Device(), bufsize, vk.BufferUsageFlags(vk.BufferUsageVertexBufferBit))
	defer vtx.Destroy()
	vtxstage := buffer.NewShared(backend.Device(), bufsize)
	vtxstage.Write(vertexdata, 0)
	// vtxstage.Write(makeColorCube(1), 0)

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
	dlayouts := []descriptor.T{ubolayout, ubolayout}

	dpool := descriptor.NewPool(backend.Device(), []vk.DescriptorPoolSize{
		{
			Type:            vk.DescriptorTypeUniformBuffer,
			DescriptorCount: 2,
		},
	})
	defer dpool.Destroy()

	desc := dpool.AllocateSets(dlayouts)
	// uboset := desc[0]

	vk.UpdateDescriptorSets(backend.Device().Ptr(), 2, []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          desc[0].Ptr(),
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeUniformBuffer,
			PBufferInfo: []vk.DescriptorBufferInfo{
				{
					Buffer: ubo[0].Ptr(),
					Offset: 0,
					Range:  vk.DeviceSize(ubosize),
				},
			},
		},
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          desc[1].Ptr(),
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeUniformBuffer,
			PBufferInfo: []vk.DescriptorBufferInfo{
				{
					Buffer: ubo[1].Ptr(),
					Offset: 0,
					Range:  vk.DeviceSize(ubosize),
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
		ctx, err := backend.Aquire()
		if err != nil {
			log.Println("Aquire() failed:", err)
			wnd.Poll()
			continue
		}

		// update ubo
		t += 0.016
		ubo[ctx.Index].Write([]mat4.T{
			mat4.PerspectiveVK(45, float32(ctx.Width)/float32(ctx.Height), 0.1, 100),
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
			cmd.CmdBindGraphicsDescriptors(layout, []descriptor.Set{desc[ctx.Index]})

			cmd.CmdBindVertexBuffer(vtx, 0)
			cmd.CmdDraw(vertices, triangles, 0, 0)

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

		// this call may cause resize events
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
