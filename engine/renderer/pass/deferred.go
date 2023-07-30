package pass

import (
	"fmt"
	"log"

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
	target  RenderTarget
	gbuffer GeometryBuffer
	quad    vertex.Mesh
	app     vulkan.App
	pass    renderpass.T
	light   LightShader
	fbuf    framebuffer.Array
	shadows Shadow

	materials *MaterialSorter

	meshQuery  *object.Query[mesh.Mesh]
	lightQuery *object.Query[light.T]
}

func NewDeferredPass(
	app vulkan.App,
	target RenderTarget,
	gbuffer GeometryBuffer,
	shadows Shadow,
) Deferred {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Deferred",
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Image:         attachment.FromImageArray(target.Output()),
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
				Image:       attachment.FromImageArray(gbuffer.Diffuse()),
			},
			{
				Name:        NormalsAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImageArray(gbuffer.Normal()),
			},
			{
				Name:        PositionAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImageArray(gbuffer.Position()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			Image:         attachment.FromImageArray(target.Depth()),
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
				Src:   renderpass.ExternalSubpass,
				Dst:   GeometrySubpass,
				Flags: core1_0.DependencyByRegion,

				// Color & depth attachments must be unused
				SrcStageMask: core1_0.PipelineStageEarlyFragmentTests |
					core1_0.PipelineStageLateFragmentTests |
					core1_0.PipelineStageFragmentShader,
				SrcAccessMask: core1_0.AccessDepthStencilAttachmentRead |
					core1_0.AccessInputAttachmentRead,

				// Before we can write to the color & depth attachments
				DstStageMask: core1_0.PipelineStageEarlyFragmentTests |
					core1_0.PipelineStageLateFragmentTests |
					core1_0.PipelineStageColorAttachmentOutput,
				DstAccessMask: core1_0.AccessColorAttachmentWrite | core1_0.AccessDepthStencilAttachmentWrite,
			},
			{
				Src:   GeometrySubpass,
				Dst:   LightingSubpass,
				Flags: core1_0.DependencyByRegion,
				// todo: consider that shadow maps should be ready before we read them in the fragment shader

				// Color attachments must be written by the geometry pass
				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				SrcAccessMask: core1_0.AccessColorAttachmentWrite,

				// Before we can read them in the lighting fragment shader
				DstAccessMask: core1_0.AccessInputAttachmentRead,
				DstStageMask:  core1_0.PipelineStageFragmentShader,
			},
			{
				Src:   LightingSubpass,
				Dst:   renderpass.ExternalSubpass,
				Flags: core1_0.DependencyByRegion,

				// Lighting subpass must finish writing the color attachments
				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				SrcAccessMask: core1_0.AccessColorAttachmentWrite,

				// Before they can be read in later fragment shaders
				DstStageMask:  core1_0.PipelineStageFragmentShader,
				DstAccessMask: core1_0.AccessInputAttachmentRead | core1_0.AccessShaderRead,
			},
		},
	})

	fbuf, err := framebuffer.NewArray(app.Frames(), app.Device(), "deferred", app.Width(), app.Height(), pass)
	if err != nil {
		panic(err)
	}

	quad := vertex.ScreenQuad("geometry-pass-quad")

	lightsh := NewLightShader(app, pass, target, gbuffer)

	app.Textures().Fetch(color.White)

	return &deferred{
		target:  target,
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
				Shader:       "deferred/textured",
				Subpass:      GeometrySubpass,
				VertexFormat: vertex.T{},
				DepthTest:    true,
				DepthWrite:   true,
			}),

		meshQuery:  object.NewQuery[mesh.Mesh](),
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *deferred) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	camera := uniform.Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         args.Position,
		Forward:     args.Forward,
	}

	//
	// geometry subpass
	//

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Context.Index])
	})

	frustum := shape.FrustumFromMatrix(args.VP)

	objects := p.meshQuery.
		Reset().
		Where(isDrawDeferred).
		Where(frustumCulled(&frustum)).
		Collect(scene)
	p.materials.Draw(cmds, args, objects)

	//
	// lighting subpass
	//

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdNextSubpass()
		p.light.Bind(cmd, args.Context.Index)
	})

	lightDesc := p.light.Descriptors(args.Context.Index)
	lightDesc.Camera.Set(camera)

	// ambient lights use a plain white texture as their shadow map
	white := p.app.Textures().Fetch(color.White)
	lightDesc.Shadow.Set(0, white)
	ambient := light.NewAmbient(color.White, 0.33)
	p.DrawLight(cmds, args, ambient, 0, 0)

	lights := p.lightQuery.
		Reset().
		Collect(scene)
	for index, lit := range lights {
		lightIndex := index + 1
		if err := p.DrawLight(cmds, args, lit, lightIndex, 5*lightIndex); err != nil {
			fmt.Printf("light draw error in object %s: %s\n", lit.Name(), err)
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *deferred) DrawLight(cmds command.Recorder, args render.Args, lit light.T, lightIndex, textureIndexOffset int) error {
	quad := p.app.Meshes().Fetch(p.quad)
	desc := lit.LightDescriptor(args, 0)

	if lit.CastShadows() {
		dirlight := uniform.Light{}
		for cascadeIndex, cascade := range lit.Cascades() {
			textureIndex := textureIndexOffset + cascadeIndex

			shadowtex := p.shadows.Shadowmap(lit, cascadeIndex)
			if shadowtex == nil {
				// no shadowmap available - disable the light until its available
				log.Println("missing cascade shadowmap", cascadeIndex)
				return nil
			}
			p.light.Descriptors(args.Context.Index).Shadow.Set(textureIndex, shadowtex)

			dirlight.ViewProj[cascadeIndex] = cascade.ViewProj
			dirlight.Distance[cascadeIndex] = cascade.FarSplit
			dirlight.Shadowmap[cascadeIndex] = uint32(textureIndex)
		}

		p.light.Descriptors(args.Context.Index).Lights.Set(lightIndex, dirlight)
	} else {
		// shadows are disabled - use a blank white texture as shadowmap
		blank := p.app.Textures().Fetch(color.White)
		p.light.Descriptors(args.Context.Index).Shadow.Set(lightIndex, blank)
	}

	cmds.Record(func(cmd command.Buffer) {
		push := &LightConst{
			ViewProj:    desc.ViewProj,
			Color:       desc.Color,
			Position:    desc.Position,
			Type:        desc.Type,
			Index:       uint32(lightIndex),
			Range:       desc.Range,
			Intensity:   desc.Intensity,
			Attenuation: desc.Attenuation,
		}
		cmd.CmdPushConstant(core1_0.StageFragment, 0, push)

		quad.Draw(cmd, 0)
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
	p.light.Destroy()
}

func isDrawDeferred(m mesh.Mesh) bool {
	return m.Mode() == mesh.Deferred
}

func frustumCulled(frustum *shape.Frustum) func(mesh.Mesh) bool {
	return func(m mesh.Mesh) bool {
		bounds := m.BoundingSphere()
		return frustum.IntersectsSphere(&bounds)
	}
}
