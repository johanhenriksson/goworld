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
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ForwardDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Lights   *descriptor.Storage[uniform.Light]
	Textures *descriptor.SamplerArray
}

type ForwardPass struct {
	target engine.Target
	app    engine.App
	pass   *renderpass.Renderpass
	fbuf   framebuffer.Array

	layout      *pipeline.Layout
	descLayout  *descriptor.Layout[*ForwardDescriptors]
	descriptors []*ForwardDescriptors
	textures    cache.SamplerCache
	objects     *ObjectBuffer
	lights      *LightBuffer
	shadows     *ShadowCache
	groups      map[material.ID]*MatGroup
	commands    []*command.IndirectDrawBuffer

	materials  cache.T[*material.Def, *Pipeline]
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

	// todo: these could probably be global descriptors
	// pass descriptor layout
	descLayout := descriptor.NewLayout(app.Device(), "Forward", &ForwardDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   objects.Size(),
		},
		Lights: &descriptor.Storage[uniform.Light]{
			Stages: core1_0.StageAll,
			Size:   lights.Size(),
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  textures.Size(),
		},
	})
	descriptors := descLayout.InstantiateMany(app.Pool(), target.Frames())
	layout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{descLayout}, nil)

	commands := make([]*command.IndirectDrawBuffer, target.Frames())
	for i := range commands {
		commands[i] = command.NewIndirectDrawBuffer(app.Device(), "Forward", objects.Size())
	}

	return &ForwardPass{
		target: target,
		app:    app,
		pass:   pass,
		fbuf:   fbuf,

		layout:      layout,
		descLayout:  descLayout,
		descriptors: descriptors,
		objects:     objects,
		lights:      lights,
		textures:    textures,
		shadows:     shadows,
		groups:      make(map[material.ID]*MatGroup, 32),
		commands:    commands,

		materials:  NewPipelineCache(app, pass, target.Frames(), layout),
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

type MatObject struct {
	Handle  int
	Indices int
}

type MatGroup struct {
	Material *Pipeline
	Objects  []MatObject
}

func (m *MatGroup) Clear() {
	m.Objects = m.Objects[:0]
}

func (m *MatGroup) Add(mat *Pipeline, handle, indices int) {
	m.Material = mat
	m.Objects = append(m.Objects, MatObject{
		Handle:  handle,
		Indices: indices,
	})
}

func (p *ForwardPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	p.Collect(args, scene)
	p.Record2(cmds, args)
}

// Within Collect, the render is allowed to query the scene.
// While Collect is executing, the scene is guaranteed to be in a consistent state.
// References to scene objects are valid until the end of the Collect function.
func (p *ForwardPass) Collect(args draw.Args, scene object.Component) {
	descriptors := p.descriptors[args.Frame]

	lights := p.lightQuery.Reset().Collect(scene)

	// update camera descriptor
	cam := uniform.CameraFromArgs(args)
	descriptors.Camera.Set(cam)

	// clear object buffer
	p.objects.Reset()

	// fill light buffer
	p.lights.Reset()
	for _, lit := range lights {
		p.lights.Store(lit.LightData(p.shadows))
	}

	// opaque pass
	opaqueQuery := p.meshQuery.
		Reset().
		Where(isDrawForward).
		Where(isTransparent(false)).
		Collect(scene)

	// empty material groups
	for _, group := range p.groups {
		group.Clear()
	}

	for _, msh := range opaqueQuery {
		gpuMesh, meshReady := p.app.Meshes().TryFetch(msh.Mesh())
		if !meshReady {
			continue
		}

		mat, matReady := p.materials.TryFetch(msh.Material())
		if !matReady {
			continue
		}

		// this could happen inside the mesh cache!
		// basically *GpuMesh could be the entire uniform object
		// or even the entire object buffer similar to the sampler cache?
		textureIds := AssignMeshTextures(p.textures, msh, mat.slots)

		objectId := p.objects.Store(uniform.Object{
			Model:    msh.Transform().Matrix(),
			Textures: textureIds,
			Vertices: gpuMesh.Vertices.Address(),
			Indices:  gpuMesh.Indices.Address(),
		})

		if _, ok := p.groups[mat.id]; !ok {
			p.groups[mat.id] = &MatGroup{
				Objects: make([]MatObject, 0, 32),
			}
		}
		p.groups[mat.id].Add(mat, objectId, gpuMesh.IndexCount)
	}

	// opaque := MaterialGroups(p.materials, args.Frame, opaqueQuery)

	// transparent pass
	// transparentQuery := p.meshQuery.
	// 	Reset().
	// 	Where(isDrawForward).
	// 	Where(isTransparent(true)).
	// 	Collect(scene)
	// transparent := DepthSortGroups(p.materials, args.Frame, cam, transparentQuery)

	// flush descriptors
	p.lights.Flush(descriptors.Lights)
	p.objects.Flush(descriptors.Objects)
	p.textures.Flush(descriptors.Textures)
}

// Record2 records commands for each object in the render context
// Record2 runs in parallel with scene updates
func (p *ForwardPass) Record2(cmds command.Recorder, args draw.Args) {
	descriptors := p.descriptors[args.Frame]
	indirect := p.commands[args.Frame]
	framebuf := p.fbuf[args.Frame]

	cmds.Record(func(cmd *command.Buffer) {
		indirect.Reset()
		cmd.CmdBeginRenderPass(p.pass, framebuf)
		cmd.CmdBindGraphicsDescriptor(p.layout, 0, descriptors)
		for _, group := range p.groups {
			group.Material.Bind(cmd)
			indirect.BeginDrawIndirect()
			for _, obj := range group.Objects {
				indirect.CmdDraw(command.Draw{
					InstanceCount: 1,

					// InstanceOffset is the index of the object properties in the object buffer
					InstanceOffset: uint32(obj.Handle),

					// Vertex count is actually the number of indices, since indexing is implemented in the shader
					VertexCount: uint32(obj.Indices),
				})
			}
			indirect.EndDrawIndirect(cmd)
		}
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
	for _, desc := range p.descriptors {
		desc.Destroy()
	}
	for _, commands := range p.commands {
		commands.Destroy()
	}
	p.layout.Destroy()
	p.descLayout.Destroy()
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
