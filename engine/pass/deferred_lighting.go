package pass

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

const LightingSubpass renderpass.Name = "lighting"

type DeferredLightPass struct {
	app        engine.App
	target     engine.Target
	gbuffer    GeometryBuffer
	ssao       engine.Target
	quad       vertex.Mesh
	pass       *renderpass.Renderpass
	light      *LightShader
	fbuf       framebuffer.Array
	samplers   *cache.SamplerCache
	shadows    *ShadowCache
	lightbuf   *uniform.LightBuffer
	lightQuery *object.Query[light.T]
}

func NewDeferredLightingPass(
	app engine.App,
	target engine.Target,
	gbuffer GeometryBuffer,
	shadows *Shadowpass,
	occlusion engine.Target,
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

	maxLights := 256
	maxShadowTextures := maxLights
	samplers := cache.NewSamplerCache(app.Textures(), maxShadowTextures)
	shadowmaps := NewShadowCache(samplers, shadows.Shadowmap)
	lightbufs := uniform.NewLightBuffer(maxLights)

	return &DeferredLightPass{
		target:     target,
		gbuffer:    gbuffer,
		app:        app,
		quad:       quad,
		light:      lightsh,
		pass:       pass,
		fbuf:       fbuf,
		shadows:    shadowmaps,
		lightbuf:   lightbufs,
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *DeferredLightPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	camera := uniform.CameraFromArgs(args)

	desc := p.light.Descriptors(args.Frame)
	desc.Camera.Set(camera)

	p.lightbuf.Reset()

	// todo: perform frustum culling on light volumes
	lights := p.lightQuery.Reset().Collect(scene)
	for _, lit := range lights {
		p.lightbuf.Store(lit.LightData(p.shadows))
	}

	p.lightbuf.Flush(desc.Lights)
	p.shadows.Flush(desc.Shadow)

	quad := p.app.Meshes().Fetch(p.quad)
	cmds.Record(func(cmd *command.Buffer) {
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
	p.samplers = nil
	p.lightbuf = nil

	p.fbuf.Destroy()
	p.fbuf = nil
	p.pass.Destroy()
	p.pass = nil
	p.light.Destroy()
	p.light = nil
}
