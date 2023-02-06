package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

const (
	LightingSubpass renderpass.Name = "lighting"
	GeometrySubpass renderpass.Name = "geometry"
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

	shadows ShadowPass

	materials *MaterialSorter
}

type DeferredSubpass interface {
	Name() renderpass.Name
	Record(command.Recorder, uniform.Camera, object.T)
	Instantiate(descriptor.Pool, renderpass.T)
	Destroy()
}

func NewGeometryPass(
	target vulkan.Target,
	shadows ShadowPass,
) Deferred {
	diffuseFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	normalFmt := core1_0.FormatR8G8B8A8UnsignedNormalized
	positionFmt := core1_0.FormatR16G16B16A16SignedFloat

	pass := renderpass.New(target.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Format:        diffuseFmt,
				Samples:       0,
				LoadOp:        core1_0.AttachmentLoadOpClear,
				StoreOp:       core1_0.AttachmentStoreOpStore,
				InitialLayout: 0,
				FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
				Clear:         color.T{},
				Usage:         core1_0.ImageUsageSampled,
				Allocator:     nil,
				Blend:         attachment.BlendAdditive,
			},
			{
				Name:        DiffuseAttachment,
				Format:      diffuseFmt,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Usage:       core1_0.ImageUsageInputAttachment | core1_0.ImageUsageTransferSrc,
			},
			{
				Name:        NormalsAttachment,
				Format:      normalFmt,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Usage:       core1_0.ImageUsageInputAttachment | core1_0.ImageUsageTransferSrc,
			},
			{
				Name:        PositionAttachment,
				Format:      positionFmt,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Usage:       core1_0.ImageUsageInputAttachment | core1_0.ImageUsageTransferSrc,
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			Usage:         core1_0.ImageUsageInputAttachment,
			ClearDepth:    1,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  GeometrySubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{DiffuseAttachment, NormalsAttachment, PositionAttachment},
			},
			{
				Name: LightingSubpass,

				ColorAttachments: []attachment.Name{OutputAttachment},
				InputAttachments: []attachment.Name{DiffuseAttachment, NormalsAttachment, PositionAttachment, attachment.DepthName},
			},
		},
		Dependencies: []renderpass.SubpassDependency{
			{
				Src: renderpass.ExternalSubpass,
				Dst: GeometrySubpass,

				SrcStageMask:  core1_0.PipelineStageBottomOfPipe,
				DstStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				SrcAccessMask: core1_0.AccessMemoryRead,
				DstAccessMask: core1_0.AccessColorAttachmentRead | core1_0.AccessColorAttachmentWrite,
				Flags:         core1_0.DependencyByRegion,
			},
			{
				Src: GeometrySubpass,
				Dst: LightingSubpass,

				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				DstStageMask:  core1_0.PipelineStageFragmentShader,
				SrcAccessMask: core1_0.AccessColorAttachmentWrite,
				DstAccessMask: core1_0.AccessShaderRead,
				Flags:         core1_0.DependencyByRegion,
			},
			{
				Src: LightingSubpass,
				Dst: renderpass.ExternalSubpass,

				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				DstStageMask:  core1_0.PipelineStageBottomOfPipe,
				SrcAccessMask: core1_0.AccessColorAttachmentRead | core1_0.AccessColorAttachmentWrite,
				DstAccessMask: core1_0.AccessMemoryRead,
				Flags:         core1_0.DependencyByRegion,
			},
		},
	})

	fbuf, err := framebuffer.New(target.Device(), target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	gbuffer := NewGbuffer(
		target,
		fbuf.Attachment(DiffuseAttachment),
		fbuf.Attachment(NormalsAttachment),
		fbuf.Attachment(PositionAttachment),
		fbuf.Attachment(OutputAttachment),
		fbuf.Attachment(attachment.DepthName),
	)

	quad := vertex.ScreenQuad("geometry-pass-quad")

	lightsh := NewLightShader(target.Device(), target.Pool(), pass)
	lightDesc := lightsh.Descriptors()

	lightDesc.Diffuse.Set(gbuffer.Diffuse())
	lightDesc.Normal.Set(gbuffer.Normal())
	lightDesc.Position.Set(gbuffer.Position())
	lightDesc.Depth.Set(gbuffer.Depth())

	shadowtex, err := texture.FromView(target.Device(), shadows.Shadowmap(), texture.Args{
		Filter: core1_0.FilterNearest,
		Wrap:   core1_0.SamplerAddressModeClampToEdge,
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
		fbuf:    fbuf,
		texture: shadowtex,

		materials: NewMaterialSorter(
			target, pass,
			&material.Def{
				Shader:       "vk/color_d",
				Subpass:      GeometrySubpass,
				VertexFormat: vertex.C{},
				DepthTest:    true,
				DepthWrite:   true,
			}),
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
	// geometry subpass
	//

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf)
	})

	frustum := shape.FrustumFromMatrix(args.VP)

	objects := object.Query[mesh.T]().
		Where(isDrawDeferred).
		Where(frustumCulled(&frustum)).
		Collect(scene)
	p.materials.Draw(cmds, args, objects)

	//
	// lighting subpass
	//

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdNextSubpass()
		p.light.Bind(cmd)
	})

	lightDesc := p.light.Descriptors()
	lightDesc.Camera.Set(camera)

	white := p.target.Textures().Fetch(texture.PathRef("textures/white.png"))
	if white != nil {
		lightDesc.Shadow.Set(0, white)

		ambient := light.NewAmbient(color.White, 0.33)
		p.DrawLight(cmds, args, ambient)
	}

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

	desc := lit.LightDescriptor(args)

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
		cmd.CmdPushConstant(core1_0.StageFragment, 0, push)

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
	p.materials.Destroy()
	p.texture.Destroy()

	p.fbuf.Destroy()
	p.pass.Destroy()
	p.gbuffer.Destroy()
	p.light.Material().Destroy()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}

func frustumCulled(frustum *shape.Frustum) func(mesh.T) bool {
	return func(m mesh.T) bool {
		bounds := m.BoundingSphere()
		return frustum.IntersectsSphere(&bounds)
	}
}
