package vkrender

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/cache"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/types"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type ShadowPass interface {
	Pass

	Shadowmap() image.View
}

type ShadowDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[Camera]
	Objects *descriptor.Storage[ObjectStorage]
}

type shadowpass struct {
	meshes    cache.MeshCache
	backend   vulkan.T
	pass      renderpass.T
	mat       material.Instance[*ShadowDescriptors]
	completed sync.Semaphore
}

func NewShadowPass(backend vulkan.T, meshes cache.MeshCache) ShadowPass {
	log.Println("create shadow pass")
	size := 1024
	pass := renderpass.New(backend.Device(), renderpass.Args{
		Frames: 1,
		Width:  size,
		Height: size,

		DepthAttachment: &renderpass.DepthAttachment{
			LoadOp:      vk.AttachmentLoadOpClear,
			StoreOp:     vk.AttachmentStoreOpStore,
			FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			Usage:       vk.ImageUsageSampledBit,
			ClearDepth:  1,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  "shadows",
				Depth: true,
			},
		},
		Dependencies: []renderpass.SubpassDependency{},
	})

	mat := material.New(
		backend.Device(),
		material.Args{
			Shader: shader.New(
				backend.Device(),
				"vk/shadow",
				shader.Inputs{
					"position": {
						Index: 0,
						Type:  types.Float,
					},
				},
				shader.Descriptors{
					"Camera":  0,
					"Objects": 1,
				},
			),
			Pass:       pass,
			Pointers:   vertex.ParsePointers(game.VoxelVertex{}),
			DepthTest:  true,
			DepthWrite: true,
		},
		&ShadowDescriptors{
			Camera: &descriptor.Uniform[Camera]{
				Stages: vk.ShaderStageAll,
			},
			Objects: &descriptor.Storage[ObjectStorage]{
				Stages: vk.ShaderStageAll,
				Size:   10,
			},
		}).Instantiate()

	return &shadowpass{
		backend:   backend,
		meshes:    meshes,
		mat:       mat,
		pass:      pass,
		completed: sync.NewSemaphore(backend.Device()),
	}
}

func (p *shadowpass) Completed() sync.Semaphore {
	return p.completed
}

func (p *shadowpass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	cmds := command.NewRecorder()

	desc := p.mat.Descriptors()

	light := query.New[light.T]().Where(func(lit light.T) bool { return lit.Type() == light.Directional }).First(scene)
	lightDesc := light.LightDescriptor()

	camera := Camera{
		ViewProj: lightDesc.ViewProj,
		Eye:      light.Transform().Position(),
	}
	desc.Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
		p.mat.Bind(cmd)
	})

	objects := query.New[mesh.T]().Where(isDrawDeferred).Collect(scene)
	for index, mesh := range objects {
		if err := p.DrawShadow(cmds, index, args, mesh); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", mesh.Name(), err)
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})

	worker := p.backend.Worker(ctx.Index)
	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Signal: []sync.Semaphore{p.completed},
		Wait: []command.Wait{
			{
				Semaphore: ctx.ImageAvailable,
				Mask:      vk.PipelineStageColorAttachmentOutputBit,
			},
		},
	})
	// worker.Wait()
}

func (p *shadowpass) Shadowmap() image.View {
	return p.pass.Attachment("depth").View(0)
}

func (p *shadowpass) DrawShadow(cmds command.Recorder, index int, args render.Args, mesh mesh.T) error {
	vkmesh := p.meshes.Fetch(mesh.Mesh())
	if vkmesh == nil {
		fmt.Println("mesh is nil")
		return nil
	}

	p.mat.Descriptors().Objects.Set(index, ObjectStorage{
		Model: mesh.Transform().World(),
	})

	cmds.Record(func(cmd command.Buffer) {
		vkmesh.Draw(cmd, index)
	})

	return nil
}

func (p *shadowpass) Destroy() {
	p.pass.Destroy()
	p.mat.Material().Destroy()
	p.completed.Destroy()
}
