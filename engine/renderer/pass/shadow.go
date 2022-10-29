package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type ShadowPass interface {
	Pass

	Shadowmap() image.View
}

type ShadowDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[uniform.Camera]
	Objects *descriptor.Storage[uniform.Object]
}

type shadowpass struct {
	meshes    cache.MeshCache
	backend   vulkan.T
	pass      renderpass.T
	passes    []DeferredSubpass
	completed sync.Semaphore
}

func NewShadowPass(backend vulkan.T, meshes cache.MeshCache, passes []DeferredSubpass) ShadowPass {
	log.Println("create shadow pass")
	size := 1024

	subpasses := make([]renderpass.Subpass, 0, len(passes))
	dependencies := make([]renderpass.SubpassDependency, 0, 2*len(passes))
	for _, gpass := range passes {
		subpasses = append(subpasses, renderpass.Subpass{
			Name:  gpass.Name(),
			Depth: true,
		})
		dependencies = append(dependencies, renderpass.SubpassDependency{
			Src: "external",
			Dst: gpass.Name(),

			SrcStageMask:  vk.PipelineStageBottomOfPipeBit,
			DstStageMask:  vk.PipelineStageColorAttachmentOutputBit,
			SrcAccessMask: vk.AccessMemoryReadBit,
			DstAccessMask: vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit,
			Flags:         vk.DependencyByRegionBit,
		})
		dependencies = append(dependencies, renderpass.SubpassDependency{
			Src: gpass.Name(),
			Dst: "external",

			SrcStageMask:  vk.PipelineStageColorAttachmentOutputBit,
			DstStageMask:  vk.PipelineStageFragmentShaderBit,
			SrcAccessMask: vk.AccessColorAttachmentWriteBit,
			DstAccessMask: vk.AccessShaderReadBit,
			Flags:         vk.DependencyByRegionBit,
		})
	}

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
		Subpasses:    subpasses,
		Dependencies: dependencies,
	})

	// instantiate geometry subpasses
	for _, gpass := range passes {
		gpass.Instantiate(pass)
	}

	return &shadowpass{
		backend:   backend,
		meshes:    meshes,
		pass:      pass,
		passes:    passes,
		completed: sync.NewSemaphore(backend.Device()),
	}
}

func (p *shadowpass) Completed() sync.Semaphore {
	return p.completed
}

func (p *shadowpass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	cmds := command.NewRecorder()

	light := query.New[light.T]().Where(func(lit light.T) bool { return lit.Type() == light.Directional }).First(scene)
	lightDesc := light.LightDescriptor()

	camera := uniform.Camera{
		ViewProj: lightDesc.ViewProj,
		Eye:      light.Transform().Position(),
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
	})

	for i, pass := range p.passes {
		pass.Record(cmds, camera, scene)

		if i < len(p.passes)-1 {
			cmds.Record(func(cmd command.Buffer) {
				cmd.CmdNextSubpass()
			})
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
}

func (p *shadowpass) Shadowmap() image.View {
	return p.pass.Attachment("depth").View(0)
}

func (p *shadowpass) Destroy() {
	for _, gpass := range p.passes {
		gpass.Destroy()
	}

	p.pass.Destroy()
	p.completed.Destroy()
}
