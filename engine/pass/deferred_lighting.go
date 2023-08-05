package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/uniform"
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
	target     vulkan.Target
	gbuffer    GeometryBuffer
	quad       vertex.Mesh
	app        vulkan.App
	pass       renderpass.T
	light      LightShader
	fbuf       framebuffer.Array
	shadows    Shadow
	shadowmaps []cache.SamplerCache
	lightbufs  []*LightBuffer
	lightQuery *object.Query[light.T]
}

func NewDeferredLightingPass(
	app vulkan.App,
	target vulkan.Target,
	gbuffer GeometryBuffer,
	shadows Shadow,
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
	})

	fbuf, err := framebuffer.NewArray(target.Frames(), app.Device(), "deferred-lighting", target.Width(), target.Height(), pass)
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

	return &DeferredLightPass{
		target:     target,
		gbuffer:    gbuffer,
		app:        app,
		quad:       quad,
		light:      lightsh,
		pass:       pass,
		shadows:    shadows,
		fbuf:       fbuf,
		shadowmaps: shadowmaps,
		lightbufs:  lightbufs,
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *DeferredLightPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
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

	p.light.Descriptors(args.Context.Index).Camera.Set(camera)

	lightbuf := p.lightbufs[args.Context.Index]
	lightbuf.Reset()

	// todo: perform frustum culling on light volumes
	lights := p.lightQuery.Reset().Collect(scene)
	for _, lit := range lights {
		lightbuf.Store(args, lit)
	}

	lightbuf.Flush()

	quad := p.app.Meshes().Fetch(p.quad)
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Context.Index])

		p.light.Bind(cmd, args.Context.Index)

		quad.Draw(cmd, 0)

		cmd.CmdEndRenderPass()
	})
}

func (p *DeferredLightPass) Name() string {
	return "Deferred Lighting"
}

func (p *DeferredLightPass) Destroy() {
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
