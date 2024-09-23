package pass

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Shadowpass struct {
	app    engine.App
	target engine.Target
	pass   *renderpass.Renderpass
	size   int

	// should be replaced with a proper cache that will evict unused maps
	shadowmaps map[light.T]Shadowmap

	lightQuery *object.Query[light.T]
	meshQuery  *object.Query[mesh.Mesh]
}

type Shadowmap struct {
	Cascades []Cascade
}

type Cascade struct {
	Texture *texture.Texture
	Frame   *framebuffer.Framebuffer
	Mats    MaterialCache
}

func NewShadowPass(app engine.App, target engine.Target) *Shadowpass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Shadow",
		DepthAttachment: &attachment.Depth{
			Image:         attachment.NewImage("shadowmap", core1_0.FormatD32SignedFloat, core1_0.ImageUsageDepthStencilAttachment|core1_0.ImageUsageInputAttachment|core1_0.ImageUsageSampled),
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			ClearDepth:    1,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,
			},
		},
		Dependencies: []renderpass.SubpassDependency{
			{
				Src:   renderpass.ExternalSubpass,
				Dst:   MainSubpass,
				Flags: core1_0.DependencyByRegion,

				// External passes must finish reading depth textures in fragment shaders
				SrcStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				SrcAccessMask: core1_0.AccessDepthStencilAttachmentRead,

				// Before we can write to the depth buffer
				DstStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				DstAccessMask: core1_0.AccessDepthStencilAttachmentWrite,
			},
			{
				Src:   MainSubpass,
				Dst:   renderpass.ExternalSubpass,
				Flags: core1_0.DependencyByRegion,

				// The shadow pass must finish writing the depth attachment
				SrcStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				SrcAccessMask: core1_0.AccessDepthStencilAttachmentWrite,

				// Before it can be used as a shadow map texture in a fragment shader
				DstStageMask:  core1_0.PipelineStageFragmentShader,
				DstAccessMask: core1_0.AccessShaderRead,
			},
		},
	})

	return &Shadowpass{
		app:        app,
		target:     target,
		pass:       pass,
		shadowmaps: make(map[light.T]Shadowmap),
		size:       2048,

		meshQuery:  object.NewQuery[mesh.Mesh](),
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *Shadowpass) Name() string {
	return "Shadow"
}

func (p *Shadowpass) createShadowmap(light light.T) Shadowmap {
	log.Println("creating shadowmap for", light.Name())

	cascades := make([]Cascade, light.Shadowmaps())
	for i := range cascades {
		key := fmt.Sprintf("%s-%d", object.Key("light", light), i)
		fbuf, err := framebuffer.New(p.app.Device(), key, p.size, p.size, p.pass)
		if err != nil {
			panic(err)
		}

		// the frame buffer object will allocate a new depth image for us
		view := fbuf.Attachment(attachment.DepthName)
		tex, err := texture.FromView(p.app.Device(), key, view, texture.Args{
			Aspect: core1_0.ImageAspectDepth,
		})
		if err != nil {
			panic(err)
		}

		cascades[i].Texture = tex
		cascades[i].Frame = fbuf

		// each light cascade needs its own shadow materials - or rather, their own descriptors
		// cheating a bit by creating entire materials for each light
		// todo: optimize this later.
		mats := NewShadowMaterialMaker(p.app, p.pass, p.target.Frames())
		cascades[i].Mats = mats
	}

	shadowmap := Shadowmap{
		Cascades: cascades,
	}
	p.shadowmaps[light] = shadowmap
	return shadowmap
}

func (p *Shadowpass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	lights := p.lightQuery.
		Reset().
		Where(func(lit light.T) bool { return lit.Type() == light.TypeDirectional && lit.CastShadows() }).
		Collect(scene)

	meshes := p.meshQuery.
		Reset().
		Where(castsShadows).
		Collect(scene)

	for _, light := range lights {
		shadowmap, mapExists := p.shadowmaps[light]
		if !mapExists {
			shadowmap = p.createShadowmap(light)
		}

		for index, cascade := range shadowmap.Cascades {
			camera := light.ShadowProjection(index)
			frame := cascade.Frame
			cmds.Record(func(cmd *command.Buffer) {
				cmd.CmdBeginRenderPass(p.pass, frame)
			})

			// todo: frustum cull meshes using light frustum

			groups := MaterialGroups(cascade.Mats, args.Frame, meshes)
			groups.Draw(cmds, camera, nil)

			cmds.Record(func(cmd *command.Buffer) {
				cmd.CmdEndRenderPass()
			})
		}
	}
}

func castsShadows(m mesh.Mesh) bool {
	return m.CastShadows()
}

func (p *Shadowpass) Shadowmap(light light.T, cascade int) *texture.Texture {
	if shadowmap, exists := p.shadowmaps[light]; exists {
		return shadowmap.Cascades[cascade].Texture
	}
	return nil
}

func (p *Shadowpass) Destroy() {
	for _, shadowmap := range p.shadowmaps {
		for _, cascade := range shadowmap.Cascades {
			cascade.Frame.Destroy()
			cascade.Texture.Destroy()
			cascade.Mats.Destroy()
		}
	}
	p.shadowmaps = nil

	p.pass.Destroy()
	p.pass = nil
}
