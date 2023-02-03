package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"

	"github.com/johanhenriksson/goworld/render/vulkan"
	"github.com/vkngwrapper/core/v2/core1_0"
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
	target    vulkan.Target
	pass      renderpass.T
	fbuf      framebuffer.T
	materials *MaterialSorter
}

func NewShadowPass(target vulkan.Target) ShadowPass {
	log.Println("create shadow pass")
	size := 1024

	subpasses := make([]renderpass.Subpass, 0, 4)
	dependencies := make([]renderpass.SubpassDependency, 0, 4)
	subpasses = append(subpasses, renderpass.Subpass{
		Name:  GeometrySubpass,
		Depth: true,
	})
	dependencies = append(dependencies, renderpass.SubpassDependency{
		Src: renderpass.ExternalSubpass,
		Dst: GeometrySubpass,

		SrcStageMask:  core1_0.PipelineStageBottomOfPipe,
		DstStageMask:  core1_0.PipelineStageColorAttachmentOutput,
		SrcAccessMask: core1_0.AccessMemoryRead,
		DstAccessMask: core1_0.AccessColorAttachmentRead | core1_0.AccessColorAttachmentWrite,
		Flags:         core1_0.DependencyByRegion,
	})
	dependencies = append(dependencies, renderpass.SubpassDependency{
		Src: GeometrySubpass,
		Dst: renderpass.ExternalSubpass,

		SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
		DstStageMask:  core1_0.PipelineStageFragmentShader,
		SrcAccessMask: core1_0.AccessColorAttachmentWrite,
		DstAccessMask: core1_0.AccessShaderRead,
		Flags:         core1_0.DependencyByRegion,
	})

	pass := renderpass.New(target.Device(), renderpass.Args{
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			Usage:         core1_0.ImageUsageSampled,
			ClearDepth:    1,
		},
		Subpasses:    subpasses,
		Dependencies: dependencies,
	})

	// todo: each light is going to need its own framebuffer
	fbuf, err := framebuffer.New(target.Device(), size, size, pass)
	if err != nil {
		panic(err)
	}

	mats := NewMaterialSorter(target, pass, &material.Def{
		Shader:       "vk/shadow",
		Subpass:      GeometrySubpass,
		VertexFormat: voxel.Vertex{},
		DepthTest:    true,
		DepthWrite:   true,
	})
	mats.TransformFn = func(d *material.Def) *material.Def {
		shadowMat := *d
		shadowMat.Shader = "vk/shadow"
		return &shadowMat
	}

	return &shadowpass{
		target:    target,
		fbuf:      fbuf,
		pass:      pass,
		materials: mats,
	}
}

func (p *shadowpass) Name() string {
	return "Shadow"
}

func (p *shadowpass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	light, lightExists := object.Query[light.T]().Where(func(lit light.T) bool { return lit.Type() == light.Directional }).First(scene)
	if !lightExists {
		return
	}
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

	objects := object.Query[mesh.T]().Where(isDrawDeferred).Collect(scene)
	p.materials.DrawCamera(cmds, args, camera, objects)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *shadowpass) Shadowmap() image.View {
	return p.fbuf.Attachment(attachment.DepthName)
}

func (p *shadowpass) Destroy() {
	p.materials.Destroy()
	p.fbuf.Destroy()
	p.pass.Destroy()
}
