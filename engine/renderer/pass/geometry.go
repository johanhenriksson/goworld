package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type DeferredPass interface {
	Pass
	GeometryBuffer
}

const (
	LightingSubpass renderpass.Name = "lighting"
)

const (
	DiffuseAttachment  attachment.Name = "diffuse"
	NormalsAttachment  attachment.Name = "normals"
	PositionAttachment attachment.Name = "position"
	OutputAttachment   attachment.Name = "output"
)

type GeometryDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type GeometryPass struct {
	GeometryBuffer

	meshes    cache.MeshCache
	quad      vertex.Mesh
	backend   vulkan.T
	pass      renderpass.T
	light     LightShader
	completed sync.Semaphore
	copyReady sync.Semaphore

	gpasses []DeferredSubpass
	shadows ShadowPass
}

type DeferredSubpass interface {
	Name() renderpass.Name
	Record(command.Recorder, uniform.Camera, object.T)
	Instantiate(renderpass.T)
	Destroy()
}

func NewGeometryPass(
	backend vulkan.T,
	meshes cache.MeshCache,
	textures cache.TextureCache,
	shadows ShadowPass,
	passes []DeferredSubpass,
) *GeometryPass {
	subpasses := make([]renderpass.Subpass, 0, len(passes)+1)
	dependencies := make([]renderpass.SubpassDependency, 0, 2*len(passes)+1)

	for _, gpass := range passes {
		subpasses = append(subpasses, renderpass.Subpass{
			Name:  gpass.Name(),
			Depth: true,

			ColorAttachments: []attachment.Name{DiffuseAttachment, NormalsAttachment, PositionAttachment},
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
			Dst: LightingSubpass,

			SrcStageMask:  vk.PipelineStageColorAttachmentOutputBit,
			DstStageMask:  vk.PipelineStageFragmentShaderBit,
			SrcAccessMask: vk.AccessColorAttachmentWriteBit,
			DstAccessMask: vk.AccessShaderReadBit,
			Flags:         vk.DependencyByRegionBit,
		})
	}

	// add lighting pass
	subpasses = append(subpasses, renderpass.Subpass{
		Name: LightingSubpass,

		ColorAttachments: []attachment.Name{OutputAttachment},
		InputAttachments: []attachment.Name{DiffuseAttachment, NormalsAttachment, PositionAttachment, attachment.DepthName},
	})
	dependencies = append(dependencies, renderpass.SubpassDependency{
		Src: LightingSubpass,
		Dst: renderpass.ExternalSubpass,

		SrcStageMask:  vk.PipelineStageColorAttachmentOutputBit,
		DstStageMask:  vk.PipelineStageBottomOfPipeBit,
		SrcAccessMask: vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit,
		DstAccessMask: vk.AccessMemoryReadBit,
		Flags:         vk.DependencyByRegionBit,
	})

	diffuseFmt := vk.FormatR8g8b8a8Unorm
	normalFmt := vk.FormatR8g8b8a8Unorm
	positionFmt := vk.FormatR16g16b16a16Sfloat

	pass := renderpass.New(backend.Device(), renderpass.Args{
		Frames: 1,
		Width:  backend.Width(),
		Height: backend.Height(),

		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Format:      diffuseFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageSampledBit,
				Blend:       attachment.BlendAdditive,
			},
			{
				Name:        DiffuseAttachment,
				Format:      diffuseFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageTransferSrcBit,
			},
			{
				Name:        NormalsAttachment,
				Format:      normalFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageTransferSrcBit,
			},
			{
				Name:        PositionAttachment,
				Format:      positionFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageTransferSrcBit,
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:      vk.AttachmentLoadOpClear,
			StoreOp:     vk.AttachmentStoreOpStore,
			FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			Usage:       vk.ImageUsageInputAttachmentBit,
			ClearDepth:  1,
		},
		Subpasses:    subpasses,
		Dependencies: dependencies,
	})

	// instantiate geometry subpasses
	for _, subpass := range passes {
		subpass.Instantiate(pass)
	}

	gbuffer := NewGbuffer(
		backend,
		pass.Attachment(DiffuseAttachment).View(0),
		pass.Attachment(NormalsAttachment).View(0),
		pass.Attachment(PositionAttachment).View(0),
		pass.Attachment(OutputAttachment).View(0),
		pass.Depth().View(0),
	)

	quad := vertex.ScreenQuad()

	lightsh := NewLightShader(backend.Device(), pass)
	lightDesc := lightsh.Descriptors()

	lightDesc.Diffuse.Set(gbuffer.Diffuse())
	lightDesc.Normal.Set(gbuffer.Normal())
	lightDesc.Position.Set(gbuffer.Position())
	lightDesc.Depth.Set(gbuffer.Depth())

	white := textures.Fetch(texture.PathRef("textures/white.png"))
	lightDesc.Shadow.Set(0, white)

	shadowtex := texture.FromView(backend.Device(), shadows.Shadowmap(), texture.Args{
		Filter: vk.FilterNearest,
		Wrap:   vk.SamplerAddressModeClampToEdge,
	})
	lightDesc.Shadow.Set(1, shadowtex)

	return &GeometryPass{
		GeometryBuffer: gbuffer,

		backend:   backend,
		meshes:    meshes,
		quad:      quad,
		light:     lightsh,
		pass:      pass,
		completed: sync.NewSemaphore(backend.Device()),
		copyReady: sync.NewSemaphore(backend.Device()),

		shadows: shadows,
		gpasses: passes,
	}
}

func (p *GeometryPass) Completed() sync.Semaphore {
	return p.completed
}

func (p *GeometryPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	cmds := command.NewRecorder()

	camera := uniform.Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         args.Position,
	}

	//
	// geometry subpasses
	//

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
	})

	for _, gpass := range p.gpasses {
		gpass.Record(cmds, camera, scene)

		cmds.Record(func(cmd command.Buffer) {
			cmd.CmdNextSubpass()
		})
	}

	//
	// lighting pass
	//

	cmds.Record(func(cmd command.Buffer) {
		p.light.Bind(cmd)
	})

	lightDesc := p.light.Descriptors()
	lightDesc.Camera.Set(camera)

	ambient := light.Descriptor{
		Type:      light.Ambient,
		Color:     color.White,
		Intensity: 0.33,
	}
	p.DrawLight(cmds, args, ambient)

	lights := query.New[light.T]().Collect(scene)
	for _, lit := range lights {
		if err := p.DrawLight(cmds, args, lit.LightDescriptor()); err != nil {
			fmt.Printf("light draw error in object %s: %s\n", lit.Name(), err)
		}
	}

	//
	// todo: forward subpasses
	//

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})

	worker := p.backend.Worker(ctx.Index)
	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Signal: []sync.Semaphore{p.completed, p.copyReady},
		Wait: []command.Wait{
			{
				Semaphore: p.shadows.Completed(),
				Mask:      vk.PipelineStageFragmentShaderBit,
			},
		},
	})

	// issue Geometry Buffer copy, so that gbuffers may be read back.
	// if more data gbuffer is to be drawn later, we need to move this to a later stage
	worker.Queue(p.GeometryBuffer.CopyBuffers())
	worker.Submit(command.SubmitInfo{
		Wait: []command.Wait{
			{
				Semaphore: p.copyReady,
				Mask:      vk.PipelineStageTopOfPipeBit,
			},
		},
	})
}

func (p *GeometryPass) DrawLight(cmds command.Recorder, args render.Args, lit light.Descriptor) error {
	vkmesh := p.meshes.Fetch(p.quad)
	cmds.Record(func(cmd command.Buffer) {

		push := LightConst{
			ViewProj:    lit.ViewProj,
			Color:       lit.Color,
			Position:    lit.Position,
			Type:        lit.Type,
			Shadowmap:   uint32(1),
			Range:       lit.Range,
			Intensity:   lit.Intensity,
			Attenuation: lit.Attenuation,
		}
		cmd.CmdPushConstant(vk.ShaderStageFragmentBit, 0, &push)

		vkmesh.Draw(cmd, 0)
	})

	return nil
}

func (p *GeometryPass) Destroy() {
	// destroy subpasses
	for _, gpass := range p.gpasses {
		gpass.Destroy()
	}

	p.pass.Destroy()
	p.GeometryBuffer.Destroy()
	p.light.Material().Destroy()
	p.completed.Destroy()
	p.copyReady.Destroy()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}