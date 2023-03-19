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
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/johanhenriksson/goworld/render/vulkan"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Shadow interface {
	Pass

	Shadowmap(light.T) texture.T
}

type shadowpass struct {
	app  vulkan.App
	pass renderpass.T
	size int

	// should be replaced with a proper cache that will evict unused maps
	shadowmaps map[light.T]Shadowmap

	lightQuery *object.Query[light.T]
	meshQuery  *object.Query[mesh.T]
}

type Shadowmap struct {
	tex  texture.T
	fbuf framebuffer.T
	mats *MaterialSorter
}

func NewShadowPass(app vulkan.App) Shadow {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Shadow",
		DepthAttachment: &attachment.Depth{
			Image:         attachment.NewImage("shadowmap", core1_0.FormatD32SignedFloat, core1_0.ImageUsageDepthStencilAttachment|core1_0.ImageUsageInputAttachment|core1_0.ImageUsageSampled),
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			ClearDepth:    1,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  GeometrySubpass,
				Depth: true,
			},
		},
		Dependencies: []renderpass.SubpassDependency{
			{
				Src:   renderpass.ExternalSubpass,
				Dst:   GeometrySubpass,
				Flags: core1_0.DependencyByRegion,

				// External passes must finish reading depth textures in fragment shaders
				SrcStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				SrcAccessMask: core1_0.AccessDepthStencilAttachmentRead,

				// Before we can write to the depth buffer
				DstStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				DstAccessMask: core1_0.AccessDepthStencilAttachmentWrite,
			},
			{
				Src:   GeometrySubpass,
				Dst:   renderpass.ExternalSubpass,
				Flags: core1_0.DependencyByRegion,

				// The shadow pass must finish writing the depth attachment
				SrcStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				SrcAccessMask: core1_0.AccessDepthStencilAttachmentWrite,

				// Before it can be used as a shadow map texture in a fragment shader
				DstStageMask:  core1_0.PipelineStageFragmentShader,
				DstAccessMask: core1_0.AccessShaderRead,
			},
		},
	})

	return &shadowpass{
		app:        app,
		pass:       pass,
		shadowmaps: make(map[light.T]Shadowmap),
		size:       4096,

		meshQuery:  object.NewQuery[mesh.T](),
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *shadowpass) Name() string {
	return "Shadow"
}

func (p *shadowpass) createShadowmap(light light.T) Shadowmap {
	log.Println("creating shadowmap for", light.Name())

	fbuf, err := framebuffer.New(p.app.Device(), p.size, p.size, p.pass)
	if err != nil {
		panic(err)
	}

	// the frame buffer object will allocate a new depth image for us
	view := fbuf.Attachment(attachment.DepthName)
	key := light.Name() + "-shadow"
	tex, err := texture.FromView(p.app.Device(), key, view, texture.Args{
		Aspect: core1_0.ImageAspectDepth,
	})
	if err != nil {
		panic(err)
	}

	// each light needs its own shadow materials - or rather, their own descriptors
	// cheating a bit by creating entire materials for each light, fix it later.
	mats := NewMaterialSorter(p.app, p.pass, &material.Def{
		Shader:       "shadow",
		Subpass:      GeometrySubpass,
		VertexFormat: voxel.Vertex{},
		DepthTest:    true,
		DepthWrite:   true,
	})
	mats.TransformFn = func(d *material.Def) *material.Def {
		shadowMat := *d
		shadowMat.Shader = "shadow"
		shadowMat.CullMode = vertex.CullFront
		return &shadowMat
	}

	shadowmap := Shadowmap{
		tex:  tex,
		fbuf: fbuf,
		mats: mats,
	}
	p.shadowmaps[light] = shadowmap
	return shadowmap
}

func (p *shadowpass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	lights := p.lightQuery.
		Reset().
		Where(func(lit light.T) bool { return lit.Type() == light.Directional && lit.Shadows() }).
		Collect(scene)

	for _, light := range lights {
		shadowmap, mapExists := p.shadowmaps[light]
		if !mapExists {
			shadowmap = p.createShadowmap(light)
		}

		lightDesc := light.LightDescriptor(args)

		camera := uniform.Camera{
			Proj:        lightDesc.Projection,
			View:        lightDesc.View,
			ViewProj:    lightDesc.ViewProj,
			ProjInv:     lightDesc.Projection.Invert(),
			ViewInv:     lightDesc.View.Invert(),
			ViewProjInv: lightDesc.ViewProj.Invert(),
			Eye:         light.Transform().Position(),
			Forward:     light.Transform().Forward(),
		}

		cmds.Record(func(cmd command.Buffer) {
			cmd.CmdBeginRenderPass(p.pass, shadowmap.fbuf)
		})

		// todo: filter only meshes that cast shadows
		meshes := p.meshQuery.
			Reset().
			Where(isDrawDeferred).
			Collect(scene)
		shadowmap.mats.DrawCamera(cmds, args, camera, meshes)

		cmds.Record(func(cmd command.Buffer) {
			cmd.CmdEndRenderPass()
		})
	}
}

func (p *shadowpass) Shadowmap(light light.T) texture.T {
	if shadowmap, exists := p.shadowmaps[light]; exists {
		return shadowmap.tex
	}
	return nil
}

func (p *shadowpass) Destroy() {
	for _, shadowmap := range p.shadowmaps {
		shadowmap.fbuf.Destroy()
		shadowmap.tex.Destroy()
		shadowmap.mats.Destroy()
	}
	p.shadowmaps = nil

	p.pass.Destroy()
	p.pass = nil
}
