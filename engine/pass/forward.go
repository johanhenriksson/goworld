package pass

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ForwardPass struct {
	target engine.Target
	app    engine.App
	pass   *renderpass.Renderpass
	fbuf   framebuffer.Array

	textures cache.SamplerCache
	objects  *ObjectBuffer
	lights   *LightBuffer
	shadows  *ShadowCache

	materials  MaterialCache
	meshQuery  *object.Query[mesh.Mesh]
	lightQuery *object.Query[light.T]
}

var _ draw.Pass = &ForwardPass{}

func NewForwardPass(
	app engine.App,
	target engine.Target,
	depth engine.Target,
	shadowPass *Shadowpass,
) *ForwardPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Forward",
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				LoadOp:        core1_0.AttachmentLoadOpLoad,
				StoreOp:       core1_0.AttachmentStoreOpStore,
				InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
				Blend:         attachment.BlendMultiply,

				Image: attachment.FromImageArray(target.Surfaces()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpLoad,
			StencilLoadOp: core1_0.AttachmentLoadOpLoad,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,

			Image: attachment.FromImageArray(depth.Surfaces()),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
	})

	fbuf, err := framebuffer.NewArray(target.Frames(), app.Device(), "forward", target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	// todo: arguments
	maxLights := 256
	maxTextures := 100
	maxObjects := 1000

	textures := cache.NewSamplerCache(app.Textures(), maxTextures)
	objects := NewObjectBuffer(maxObjects)
	lights := NewLightBuffer(maxLights)
	shadows := NewShadowCache(textures, shadowPass.Shadowmap)

	return &ForwardPass{
		target: target,
		app:    app,
		pass:   pass,
		fbuf:   fbuf,

		objects:    objects,
		lights:     lights,
		textures:   textures,
		shadows:    shadows,
		materials:  NewForwardMaterialCache(app, pass, target.Frames(), textures, objects, lights, shadows),
		meshQuery:  object.NewQuery[mesh.Mesh](),
		lightQuery: object.NewQuery[light.T](),
	}
}

// premise:
// - its reasonable to have a separate object buffer for each render pass.
//   this is because there is no overlap between the objects that are rendered in each pass.
//   however, it might be faster to do a single descriptor set write
// - two phases:
//   - collect: gather objects that will be rendered. synchronized between update/render
//   - record: record commands for each object. runs on the render thread

type Context struct{}

func (p *ForwardPass) Collect(ctx *Context, scene object.Component) {
	//
	// opaque := p.meshQuery.
	// 	Reset().
	// 	Where(isDrawForward).
	// 	Where(isTransparent(false)).
	// 	Collect(scene)
	//
	// for _, mesh := range opaque {
	// 	gpuMesh, ok := p.app.Meshes().TryFetch(mesh.Mesh())
	// 	if !ok {
	// 		continue
	// 	}
	//
	// 	// material
	// 	p.materials.TryFetch(mesh.Material())
	//
	// }
}

// Record2 records commands for each object in the render context
func (p *ForwardPass) Record2(cmds command.Recorder, ctx *Context) {
}

func (p *ForwardPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	cam := uniform.CameraFromArgs(args)
	lights := p.lightQuery.Reset().Collect(scene)

	// clear object buffer
	p.lights.Reset()
	p.objects.Reset()

	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Frame])
	})

	// opaque pass
	opaque := p.meshQuery.
		Reset().
		Where(isDrawForward).
		Where(isTransparent(false)).
		Collect(scene)
	groups := MaterialGroups(p.materials, args.Frame, opaque)
	groups.Draw(cmds, cam, lights)

	// transparent pass
	transparent := p.meshQuery.
		Reset().
		Where(isDrawForward).
		Where(isTransparent(true)).
		Collect(scene)
	groups = DepthSortGroups(p.materials, args.Frame, cam, transparent)
	groups.Draw(cmds, cam, lights)

	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *ForwardPass) Name() string {
	return "Forward"
}

func (p *ForwardPass) Destroy() {
	p.textures.Destroy()
	p.fbuf.Destroy()
	p.pass.Destroy()
	p.materials.Destroy()
}

func isDrawForward(m mesh.Mesh) bool {
	if mat := m.Material(); mat != nil {
		return mat.Pass == material.Forward
	}
	return false
}

func isTransparent(transparent bool) func(m mesh.Mesh) bool {
	return func(m mesh.Mesh) bool {
		if mat := m.Material(); mat != nil {
			return m.Material().Transparent == transparent
		}
		return false
	}
}
