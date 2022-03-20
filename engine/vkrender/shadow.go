package vkrender

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type ShadowPass interface {
	Pass

	Shadowmap() image.View
}

type ShadowDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[CameraData]
	Objects *descriptor.Storage[ObjectStorage]
}

type shadowpass struct {
	meshes    MeshCache
	backend   vulkan.T
	pass      renderpass.T
	mat       material.Instance[*ShadowDescriptors]
	completed sync.Semaphore
}

func NewShadowPass(backend vulkan.T, meshes MeshCache) ShadowPass {
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
			Pass:     pass,
			Pointers: vertex.ParsePointers(game.VoxelVertex{}),
		},
		&ShadowDescriptors{
			Camera: &descriptor.Uniform[CameraData]{
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

	camera := CameraData{
		ViewProj: lightDesc.ViewProj,
		Eye:      light.Transform().Position(),
	}
	desc.Camera.Set(camera)

	desc.Objects.Set(0, ObjectStorage{
		Model: mat4.Ident(),
	})

	desc.Objects.Set(1, ObjectStorage{
		Model: mat4.Translate(vec3.New(-16, 0, 0)),
	})

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
		p.mat.Bind(cmd)
	})

	objects := query.New[mesh.T]().Where(isDrawDeferred).Collect(scene)
	for _, mesh := range objects {
		if err := p.DrawDeferred(cmds, args, mesh); err != nil {
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

func (p *shadowpass) DrawDeferred(cmds command.Recorder, args render.Args, mesh mesh.T) error {
	args = args.Apply(mesh.Transform().World())

	vkmesh := p.meshes.Fetch(mesh.Mesh())
	if vkmesh == nil {
		fmt.Println("mesh is nil")
		return nil
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBindVertexBuffer(vkmesh.Vertices, 0)
		cmd.CmdBindIndexBuffers(vkmesh.Indices, 0, vk.IndexTypeUint16)

		// index of the object properties in the ssbo
		idx := 0
		count := 2
		cmd.CmdDrawIndexed(vkmesh.Mesh.Elements(), count, 0, 0, idx)
	})

	return nil
}

func (p *shadowpass) Destroy() {
	p.pass.Destroy()
	p.mat.Material().Destroy()
	p.completed.Destroy()
}
