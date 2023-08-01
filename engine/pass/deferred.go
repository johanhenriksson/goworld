package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
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

	materials  *MaterialSorter
	shadowmaps []cache.SamplerCache
	lightbufs  []*LightBuffer

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

	lightsh := NewLightShader(app, pass, gbuffer)

	lightbufs := make([]*LightBuffer, target.Frames())
	shadowmaps := make([]cache.SamplerCache, target.Frames())
	for i := range lightbufs {
		shadowmaps[i] = cache.NewSamplerCache(app.Textures(), lightsh.Descriptors(i).Shadow)
		lightbufs[i] = NewLightBuffer(lightsh.Descriptors(i).Lights, shadowmaps[i], shadows.Shadowmap)
	}

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

		materials:  NewMaterialSorter(app, target, pass, material.StandardDeferred()),
		shadowmaps: shadowmaps,
		lightbufs:  lightbufs,

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

	lightbuf := p.lightbufs[args.Context.Index]
	lightbuf.Reset()

	lightDesc := p.light.Descriptors(args.Context.Index)
	lightDesc.Camera.Set(camera)

	// ambient lights use a plain white texture as their shadow map
	ambient := light.NewAmbient(color.White, 0.33)
	lightbuf.Store(args, ambient)

	// todo: perform frustum culling on light volumes
	lights := p.lightQuery.
		Reset().
		Collect(scene)
	for _, lit := range lights {
		lightbuf.Store(args, lit)
	}

	lightbuf.Flush()

	quad := p.app.Meshes().Fetch(p.quad)
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdNextSubpass()
		p.light.Bind(cmd, args.Context.Index)

		cmd.CmdPushConstant(core1_0.StageFragment, 0, &LightConst{
			Count: uint32(lightbuf.Count()),
		})
		quad.Draw(cmd, 0)

		cmd.CmdEndRenderPass()
	})
}

func (p *deferred) Name() string {
	return "Deferred"
}

func (p *deferred) Destroy() {
	// destroy subpasses
	p.materials.Destroy()
	p.materials = nil

	for _, cache := range p.shadowmaps {
		cache.Destroy()
	}
	p.shadowmaps = nil
	p.lightbufs = nil

	p.fbuf.Destroy()
	p.fbuf = nil
	p.pass.Destroy()
	p.pass = nil
	p.light.Destroy()
	p.light = nil
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
