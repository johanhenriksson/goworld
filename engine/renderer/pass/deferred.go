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
}

type GeometryDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type deferred struct {
	gbuffer GeometryBuffer
	quad    vertex.Mesh
	app     vulkan.App
	pass    renderpass.T
	light   LightShader
	fbuf    framebuffer.T
	shadows Shadow

	materials *MaterialSorter
}

func NewDeferredPass(
	app vulkan.App,
	gbuffer GeometryBuffer,
	shadows Shadow,
) Deferred {
	pass := renderpass.New(app.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Image:         attachment.FromImageArray(gbuffer.Output()),
				Samples:       0,
				LoadOp:        core1_0.AttachmentLoadOpClear,
				StoreOp:       core1_0.AttachmentStoreOpStore,
				InitialLayout: 0,
				FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
				Clear:         color.T{},
				Blend:         attachment.BlendAdditive,
			},
			{
				Name:        DiffuseAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImage(gbuffer.Diffuse()),
			},
			{
				Name:        NormalsAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImage(gbuffer.Normal()),
			},
			{
				Name:        PositionAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImage(gbuffer.Position()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			Image:         attachment.FromImageArray(gbuffer.Depth()),
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

	fbuf, err := framebuffer.New(app.Device(), app.Width(), app.Height(), pass)
	if err != nil {
		panic(err)
	}

	quad := vertex.ScreenQuad("geometry-pass-quad")

	lightsh := NewLightShader(app, pass, gbuffer)

	app.Textures().Fetch(color.White)

	return &deferred{
		gbuffer: gbuffer,
		app:     app,
		quad:    quad,
		light:   lightsh,
		pass:    pass,

		shadows: shadows,
		fbuf:    fbuf,

		materials: NewMaterialSorter(
			app, pass,
			&material.Def{
				Shader:       "color_d",
				Subpass:      GeometrySubpass,
				VertexFormat: vertex.C{},
				DepthTest:    true,
				DepthWrite:   true,
			}),
	}
}

func (p *deferred) Record(cmds command.Recorder, args render.Args, scene object.T) {
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

	// ambient lights use a plain white texture as their shadow map
	white, shadowTexReady := p.app.Textures().Fetch(color.White)
	if shadowTexReady {
		lightDesc.Shadow.Set(0, white)

		ambient := light.NewAmbient(color.White, 0.33)
		p.DrawLight(cmds, args, ambient, 0)
	}

	lights := object.Query[light.T]().Collect(scene)
	for index, lit := range lights {
		if err := p.DrawLight(cmds, args, lit, index+1); err != nil {
			fmt.Printf("light draw error in object %s: %s\n", lit.Name(), err)
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *deferred) DrawLight(cmds command.Recorder, args render.Args, lit light.T, shadowIndex int) error {
	vkmesh, meshReady := p.app.Meshes().Fetch(p.quad)
	if !meshReady {
		return nil
	}

	desc := lit.LightDescriptor(args)

	shadowtex := p.shadows.Shadowmap(lit)
	if shadowtex != nil {
		p.light.Descriptors().Shadow.Set(shadowIndex, shadowtex)
	} else {
		// no shadowmap available - disable the light until its available
		if lit.Shadows() {
			return nil
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		push := &LightConst{
			ViewProj:    desc.ViewProj,
			Color:       desc.Color,
			Position:    desc.Position,
			Type:        desc.Type,
			Shadowmap:   uint32(shadowIndex),
			Range:       desc.Range,
			Intensity:   desc.Intensity,
			Attenuation: desc.Attenuation,
		}
		cmd.CmdPushConstant(core1_0.StageFragment, 0, push)

		vkmesh.Draw(cmd, 0)
	})

	return nil
}

func (p *deferred) Name() string {
	return "Deferred"
}

func (p *deferred) Destroy() {
	// destroy subpasses
	p.materials.Destroy()

	p.fbuf.Destroy()
	p.pass.Destroy()
	p.gbuffer.Destroy()
	p.light.Destroy()
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
