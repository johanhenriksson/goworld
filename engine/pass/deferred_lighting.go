package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

const LightingSubpass renderpass.Name = "lighting"

type DeferredLightPass struct {
	app        vulkan.App
	target     vulkan.Target
	gbuffer    GeometryBuffer
	ssao       vulkan.Target
	quad       vertex.Mesh
	pass       renderpass.T
	light      LightShader
	fbuf       framebuffer.Array
	samplers   []cache.SamplerCache
	shadows    []*ShadowCache
	lightbufs  []*LightBuffer
	lightQuery *object.Query[light.T]
}

func NewDeferredLightingPass(
	app vulkan.App,
	target vulkan.Target,
	gbuffer GeometryBuffer,
	shadows Shadow,
	occlusion vulkan.Target,
) *DeferredLightPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Deferred Lighting",
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
		},
		Subpasses: []renderpass.Subpass{
			{
				Name: LightingSubpass,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
		Dependencies: []renderpass.SubpassDependency{
			{
				// For color attachment operations
				Src:           renderpass.ExternalSubpass,
				Dst:           MainSubpass,
				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				DstStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				SrcAccessMask: core1_0.AccessColorAttachmentWrite,
				DstAccessMask: core1_0.AccessColorAttachmentWrite | core1_0.AccessColorAttachmentRead,
				Flags:         core1_0.DependencyByRegion,
			},
			{
				// For fragment shader reads
				Src:           renderpass.ExternalSubpass,
				Dst:           MainSubpass,
				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				DstStageMask:  core1_0.PipelineStageFragmentShader,
				SrcAccessMask: core1_0.AccessColorAttachmentWrite,
				DstAccessMask: core1_0.AccessShaderRead,
				Flags:         core1_0.DependencyByRegion,
			},
		},
	})

	fbuf, err := framebuffer.NewArray(target.Frames(), app.Device(), "deferred-lighting", target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	quad := vertex.ScreenQuad("geometry-pass-quad")

	lightsh := NewLightShader(app, pass, gbuffer, occlusion)

	samplers := make([]cache.SamplerCache, target.Frames())
	lightbufs := make([]*LightBuffer, target.Frames())
	shadowmaps := make([]*ShadowCache, target.Frames())
	for i := range lightbufs {
		samplers[i] = cache.NewSamplerCache(app.Textures(), lightsh.Descriptors(i).Shadow)
		shadowmaps[i] = NewShadowCache(samplers[i], shadows.Shadowmap)
		lightbufs[i] = NewLightBuffer(256)
	}

	return &DeferredLightPass{
		target:     target,
		gbuffer:    gbuffer,
		app:        app,
		quad:       quad,
		light:      lightsh,
		pass:       pass,
		fbuf:       fbuf,
		shadows:    shadowmaps,
		lightbufs:  lightbufs,
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *DeferredLightPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	camera := CameraFromArgs(args)

	desc := p.light.Descriptors(args.Frame)
	desc.Camera.Set(camera)

	lightbuf := p.lightbufs[args.Frame]
	shadows := p.shadows[args.Frame]
	lightbuf.Reset()

	// todo: perform frustum culling on light volumes
	lights := p.lightQuery.Reset().Collect(scene)
	for _, lit := range lights {
		lightbuf.Store(lit.LightData(shadows))
	}

	lightbuf.Flush(desc.Lights)
	shadows.Flush()

	quad := p.app.Meshes().Fetch(p.quad)
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Frame])

		p.light.Bind(cmd, args.Frame)

		quad.Bind(cmd)
		quad.Draw(cmd, 0)

		cmd.CmdEndRenderPass()
	})
}

func (p *DeferredLightPass) Name() string {
	return "Deferred Lighting"
}

func (p *DeferredLightPass) Destroy() {
	for _, cache := range p.samplers {
		cache.Destroy()
	}
	p.samplers = nil
	p.lightbufs = nil

	p.fbuf.Destroy()
	p.fbuf = nil
	p.pass.Destroy()
	p.pass = nil
	p.light.Destroy()
	p.light = nil
}
