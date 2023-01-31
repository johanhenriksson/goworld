package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

const (
	LightingSubpass renderpass.Name = "lighting"
)

const (
	DiffuseAttachment  attachment.Name = "diffuse"
	NormalsAttachment  attachment.Name = "normals"
	PositionAttachment attachment.Name = "position"
	OutputAttachment   attachment.Name = "output"
)

type Deferred interface {
	Pass
	GBuffer() GeometryBuffer
}

type GeometryDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type GeometryPass struct {
	gbuffer GeometryBuffer
	quad    vertex.Mesh
	target  vulkan.Target
	pass    renderpass.T
	light   LightShader
	fbuf    framebuffer.T
	texture texture.T

	gpasses []DeferredSubpass
	shadows ShadowPass
}

type DeferredSubpass interface {
	Name() renderpass.Name
	Record(command.Recorder, uniform.Camera, object.T)
	Instantiate(descriptor.Pool, renderpass.T)
	Destroy()
}

func NewGeometryPass(
	target vulkan.Target,
	pool descriptor.Pool,
	shadows ShadowPass,
	passes []DeferredSubpass,
) Deferred {
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

	pass := renderpass.New(target.Device(), renderpass.Args{
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
			LoadOp:        vk.AttachmentLoadOpClear,
			StencilLoadOp: vk.AttachmentLoadOpClear,
			StoreOp:       vk.AttachmentStoreOpStore,
			FinalLayout:   vk.ImageLayoutShaderReadOnlyOptimal,
			Usage:         vk.ImageUsageInputAttachmentBit,
			ClearDepth:    1,
		},
		Subpasses:    subpasses,
		Dependencies: dependencies,
	})

	fbuf, err := framebuffer.New(target.Device(), target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	// instantiate geometry subpasses
	for _, subpass := range passes {
		subpass.Instantiate(pool, pass)
	}

	gbuffer := NewGbuffer(
		target,
		fbuf.Attachment(DiffuseAttachment),
		fbuf.Attachment(NormalsAttachment),
		fbuf.Attachment(PositionAttachment),
		fbuf.Attachment(OutputAttachment),
		fbuf.Attachment(attachment.DepthName),
	)

	quad := vertex.ScreenQuad()

	lightsh := NewLightShader(target.Device(), pool, pass)
	lightDesc := lightsh.Descriptors()

	lightDesc.Diffuse.Set(gbuffer.Diffuse())
	lightDesc.Normal.Set(gbuffer.Normal())
	lightDesc.Position.Set(gbuffer.Position())
	lightDesc.Depth.Set(gbuffer.Depth())

	shadowtex, err := texture.FromView(target.Device(), shadows.Shadowmap(), texture.Args{
		Filter: vk.FilterNearest,
		Wrap:   vk.SamplerAddressModeClampToEdge,
	})
	if err != nil {
		panic(err)
	}
	lightDesc.Shadow.Set(1, shadowtex)
	target.Textures().Fetch(texture.PathRef("textures/white.png")) // warmup texture

	return &GeometryPass{
		gbuffer: gbuffer,
		target:  target,
		quad:    quad,
		light:   lightsh,
		pass:    pass,

		shadows: shadows,
		gpasses: passes,
		fbuf:    fbuf,
		texture: shadowtex,
	}
}

func (p *GeometryPass) Record(cmds command.Recorder, args render.Args, scene object.T) {
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
		cmd.CmdBeginRenderPass(p.pass, p.fbuf)
	})

	for _, gpass := range p.gpasses {
		gpass.Record(cmds, camera, scene)

		cmds.Record(func(cmd command.Buffer) {
			cmd.CmdNextSubpass()
		})
	}

	//
	// lighting subpass
	//

	cmds.Record(func(cmd command.Buffer) {
		p.light.Bind(cmd)
	})

	lightDesc := p.light.Descriptors()
	lightDesc.Camera.Set(camera)

	white := p.target.Textures().Fetch(texture.PathRef("textures/white.png"))
	lightDesc.Shadow.Set(0, white)

	ambient := light.NewAmbient(color.White, 0.33)
	p.DrawLight(cmds, args, ambient)

	lights := object.Query[light.T]().Collect(scene)
	for _, lit := range lights {
		if err := p.DrawLight(cmds, args, lit); err != nil {
			fmt.Printf("light draw error in object %s: %s\n", lit.Name(), err)
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *GeometryPass) DrawLight(cmds command.Recorder, args render.Args, lit light.T) error {
	vkmesh := p.target.Meshes().Fetch(p.quad)
	if vkmesh == nil {
		return nil
	}

	desc := lit.LightDescriptor()

	cmds.Record(func(cmd command.Buffer) {
		push := &LightConst{
			ViewProj:    desc.ViewProj,
			Color:       desc.Color,
			Position:    desc.Position,
			Type:        desc.Type,
			Shadowmap:   uint32(1),
			Range:       desc.Range,
			Intensity:   desc.Intensity,
			Attenuation: desc.Attenuation,
		}
		cmd.CmdPushConstant(vk.ShaderStageFragmentBit, 0, push)

		vkmesh.Draw(cmd, 0)
	})

	return nil
}

func (p *GeometryPass) Name() string {
	return "Geometry"
}

func (d *GeometryPass) GBuffer() GeometryBuffer {
	return d.gbuffer
}

func (p *GeometryPass) Destroy() {
	// destroy subpasses
	for _, gpass := range p.gpasses {
		gpass.Destroy()
	}

	p.texture.Destroy()

	p.fbuf.Destroy()
	p.pass.Destroy()
	p.gbuffer.Destroy()
	p.light.Material().Destroy()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}
