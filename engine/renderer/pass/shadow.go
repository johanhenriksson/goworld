package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
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
	backend   vulkan.T
	pass      renderpass.T
	passes    []DeferredSubpass
	completed sync.Semaphore
	fbuf      framebuffer.T
}

func NewShadowPass(backend vulkan.T, passes []DeferredSubpass) ShadowPass {
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
			Src: renderpass.ExternalSubpass,
			Dst: gpass.Name(),

			SrcStageMask:  vk.PipelineStageBottomOfPipeBit,
			DstStageMask:  vk.PipelineStageColorAttachmentOutputBit,
			SrcAccessMask: vk.AccessMemoryReadBit,
			DstAccessMask: vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit,
			Flags:         vk.DependencyByRegionBit,
		})
		dependencies = append(dependencies, renderpass.SubpassDependency{
			Src: gpass.Name(),
			Dst: renderpass.ExternalSubpass,

			SrcStageMask:  vk.PipelineStageColorAttachmentOutputBit,
			DstStageMask:  vk.PipelineStageFragmentShaderBit,
			SrcAccessMask: vk.AccessColorAttachmentWriteBit,
			DstAccessMask: vk.AccessShaderReadBit,
			Flags:         vk.DependencyByRegionBit,
		})
	}

	pass := renderpass.New(backend.Device(), renderpass.Args{
		DepthAttachment: &attachment.Depth{
			LoadOp:        vk.AttachmentLoadOpClear,
			StencilLoadOp: vk.AttachmentLoadOpClear,
			StoreOp:       vk.AttachmentStoreOpStore,
			FinalLayout:   vk.ImageLayoutShaderReadOnlyOptimal,
			Usage:         vk.ImageUsageSampledBit,
			ClearDepth:    1,
		},
		Subpasses:    subpasses,
		Dependencies: dependencies,
	})

	// todo: each light is going to need its own framebuffer
	fbuf, err := framebuffer.New(backend.Device(), size, size, pass)
	if err != nil {
		panic(err)
	}

	// instantiate geometry subpasses
	for _, gpass := range passes {
		gpass.Instantiate(pass)
	}

	return &shadowpass{
		backend:   backend,
		fbuf:      fbuf,
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
		Proj:        lightDesc.Projection,
		View:        lightDesc.View,
		ViewProj:    lightDesc.ViewProj,
		ProjInv:     lightDesc.Projection.Invert(),
		ViewInv:     lightDesc.View.Invert(),
		ViewProjInv: lightDesc.ViewProj.Invert(),
		Eye:         light.Transform().Position(),
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf)
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
	return p.fbuf.Attachment(attachment.DepthName)
}

func (p *shadowpass) Destroy() {
	for _, gpass := range p.passes {
		gpass.Destroy()
	}

	p.fbuf.Destroy()
	p.pass.Destroy()
	p.completed.Destroy()
}
