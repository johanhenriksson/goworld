package pass

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/uniform"
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

const LightingSubpass renderpass.Name = "lighting"

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
	target  vulkan.Target
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
	target vulkan.Target,
	depth vulkan.Target,
	gbuffer GeometryBuffer,
	shadows Shadow,
) Deferred {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Deferred",
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Image:         attachment.FromImageArray(target.Surfaces()),
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
			Image:         attachment.FromImageArray(depth.Surfaces()),
			ClearDepth:    1,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
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
				Dst:   MainSubpass,
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
				Src:   MainSubpass,
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

	fbuf, err := framebuffer.NewArray(target.Frames(), app.Device(), "deferred", target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	quad := vertex.ScreenQuad("geometry-pass-quad")

	lightsh := NewLightShader(app, pass, depth, gbuffer)

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

		materials: NewMaterialSorter(app, target, pass, material.StandardDeferred()),

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

	lightDesc := p.light.Descriptors(args.Context.Index)
	lightDesc.Camera.Set(camera)

	// ambient lights use a plain white texture as their shadow map
	white := p.app.Textures().Fetch(color.White)
	lightDesc.Shadow.Set(0, white)
	ambient := light.NewAmbient(color.White, 0.33)
	p.UpdateLight(cmds, args, ambient, 0, 0)

	// todo: perform frustum culling on light volumes
	lights := p.lightQuery.
		Reset().
		Collect(scene)
	for index, lit := range lights {
		lightIndex := index + 1
		if err := p.UpdateLight(cmds, args, lit, lightIndex, 5*lightIndex); err != nil {
			fmt.Printf("light draw error in object %s: %s\n", lit.Name(), err)
		}
	}

	quad := p.app.Meshes().Fetch(p.quad)
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdNextSubpass()
		p.light.Bind(cmd, args.Context.Index)

		cmd.CmdPushConstant(core1_0.StageFragment, 0, &LightConst{
			Count: uint32(len(lights) + 1),
		})
		quad.Draw(cmd, 0)

		cmd.CmdEndRenderPass()
	})
}

func (p *deferred) UpdateLight(cmds command.Recorder, args render.Args, lit light.T, lightIndex, textureIndexOffset int) error {
	desc := lit.LightDescriptor(args, 0)

	entry := uniform.Light{
		Type:      desc.Type,
		Color:     desc.Color,
		Position:  desc.Position,
		Intensity: desc.Intensity,
	}

	switch lit.(type) {
	case *light.Point:
		entry.Attenuation = desc.Attenuation
		entry.Range = desc.Range

	case *light.Directional:
		for cascadeIndex, cascade := range lit.Cascades() {
			textureIndex := textureIndexOffset + cascadeIndex

			if shadowtex := p.shadows.Shadowmap(lit, cascadeIndex); shadowtex != nil {
				p.light.Descriptors(args.Context.Index).Shadow.Set(textureIndex, shadowtex)
			} else {
				// no shadowmap available - disable shadows until its available
				log.Println("missing cascade shadowmap", cascadeIndex)
				textureIndex = 0
			}

			entry.ViewProj[cascadeIndex] = cascade.ViewProj
			entry.Distance[cascadeIndex] = cascade.FarSplit
			entry.Shadowmap[cascadeIndex] = uint32(textureIndex)
		}

	default:
		// shadows are disabled - use a blank white texture as shadowmap
		blank := p.app.Textures().Fetch(color.White)
		p.light.Descriptors(args.Context.Index).Shadow.Set(lightIndex, blank)
	}

	p.light.Descriptors(args.Context.Index).Lights.Set(lightIndex, entry)

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
	if mat := m.Material(); mat != nil {
		return mat.Pass == material.Deferred
	}
	return false
}

func frustumCulled(frustum *shape.Frustum) func(mesh.Mesh) bool {
	return func(m mesh.Mesh) bool {
		bounds := m.BoundingSphere()
		return frustum.IntersectsSphere(&bounds)
	}
}
