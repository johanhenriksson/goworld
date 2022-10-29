package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type DeferredPass interface {
	Pass
	GeometryBuffer
}

type GeometryDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type GeometryPass struct {
	GeometryBuffer

	meshes    cache.MeshCache
	quad      vertex.Mesh
	backend   vulkan.T
	pass      renderpass.T
	light     material.Instance[*LightDescriptors]
	completed sync.Semaphore
	copyReady sync.Semaphore

	gpasses []DeferredSubpass
	shadows ShadowPass
}

type DeferredSubpass interface {
	Name() string
	Record(command.Recorder, uniform.Camera, object.T)
	Instantiate(renderpass.T)
	Destroy()
}

func NewGeometryPass(
	backend vulkan.T,
	meshes cache.MeshCache,
	textures cache.TextureCache,
	shadows ShadowPass,
	passes []DeferredSubpass,
) *GeometryPass {
	subpasses := make([]renderpass.Subpass, 0, len(passes)+1)
	dependencies := make([]renderpass.SubpassDependency, 0, 2*len(passes)+1)

	for _, gpass := range passes {
		subpasses = append(subpasses, renderpass.Subpass{
			Name:  gpass.Name(),
			Depth: true,

			ColorAttachments: []string{"diffuse", "normal", "position"},
		})
		dependencies = append(dependencies, renderpass.SubpassDependency{
			Src: "external",
			Dst: gpass.Name(),

			SrcStageMask:  vk.PipelineStageBottomOfPipeBit,
			DstStageMask:  vk.PipelineStageColorAttachmentOutputBit,
			SrcAccessMask: vk.AccessMemoryReadBit,
			DstAccessMask: vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit,
			Flags:         vk.DependencyByRegionBit,
		})
		dependencies = append(dependencies, renderpass.SubpassDependency{
			Src: gpass.Name(),
			Dst: "lighting",

			SrcStageMask:  vk.PipelineStageColorAttachmentOutputBit,
			DstStageMask:  vk.PipelineStageFragmentShaderBit,
			SrcAccessMask: vk.AccessColorAttachmentWriteBit,
			DstAccessMask: vk.AccessShaderReadBit,
			Flags:         vk.DependencyByRegionBit,
		})
	}

	// add lighting pass
	subpasses = append(subpasses, renderpass.Subpass{
		Name: "lighting",

		ColorAttachments: []string{"output"},
		InputAttachments: []string{"diffuse", "normal", "position", "depth"},
	})
	dependencies = append(dependencies, renderpass.SubpassDependency{
		Src: "lighting",
		Dst: "external",

		SrcStageMask:  vk.PipelineStageColorAttachmentOutputBit,
		DstStageMask:  vk.PipelineStageBottomOfPipeBit,
		SrcAccessMask: vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit,
		DstAccessMask: vk.AccessMemoryReadBit,
		Flags:         vk.DependencyByRegionBit,
	})

	diffuseFmt := vk.FormatR8g8b8a8Unorm
	normalFmt := vk.FormatR8g8b8a8Unorm
	positionFmt := vk.FormatR16g16b16a16Sfloat

	pass := renderpass.New(backend.Device(), renderpass.Args{
		Frames: 1, // backend.Frames(),
		Width:  backend.Width(),
		Height: backend.Height(),

		ColorAttachments: []renderpass.ColorAttachment{
			{
				Name:        "output",
				Format:      diffuseFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageSampledBit,
				Blend:       renderpass.BlendAdditive,
			},
			{
				Name:        "diffuse",
				Format:      diffuseFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageTransferSrcBit,
			},
			{
				Name:        "normal",
				Format:      normalFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageTransferSrcBit,
			},
			{
				Name:        "position",
				Format:      positionFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageTransferSrcBit,
			},
		},
		DepthAttachment: &renderpass.DepthAttachment{
			LoadOp:      vk.AttachmentLoadOpClear,
			StoreOp:     vk.AttachmentStoreOpStore,
			FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			Usage:       vk.ImageUsageInputAttachmentBit,
			ClearDepth:  1,
		},
		Subpasses:    subpasses,
		Dependencies: dependencies,
	})

	// instantiate geometry subpasses
	for _, gpass := range passes {
		gpass.Instantiate(pass)
	}

	gbuffer := NewGbuffer(backend, pass)

	quad := vertex.NewTriangles("screen_quad", []vertex.T{
		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
		{P: vec3.New(-1, 1, 0), T: vec2.New(0, 1)},
		{P: vec3.New(1, -1, 0), T: vec2.New(1, 0)},
	}, []uint16{
		0, 1, 2,
		0, 3, 1,
	})

	lightsh := NewLightShader(backend.Device(), pass)
	lightDesc := lightsh.Descriptors()

	lightDesc.Diffuse.Set(gbuffer.Diffuse(0))
	lightDesc.Normal.Set(gbuffer.Normal(0))
	lightDesc.Position.Set(gbuffer.Position(0))
	lightDesc.Depth.Set(gbuffer.Depth(0))

	white := textures.Fetch(texture.PathRef("assets/textures/white.png"))
	lightDesc.Shadow.Set(0, white)

	shadowtex := texture.FromView(backend.Device(), shadows.Shadowmap(), texture.Args{
		Filter: vk.FilterNearest,
		Wrap:   vk.SamplerAddressModeClampToEdge,
	})
	lightDesc.Shadow.Set(1, shadowtex)

	return &GeometryPass{
		GeometryBuffer: gbuffer,

		backend:   backend,
		meshes:    meshes,
		quad:      quad,
		light:     lightsh,
		pass:      pass,
		completed: sync.NewSemaphore(backend.Device()),
		copyReady: sync.NewSemaphore(backend.Device()),

		shadows: shadows,
		gpasses: passes,
	}
}

func (p *GeometryPass) Completed() sync.Semaphore {
	return p.completed
}

func (p *GeometryPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	cmds := command.NewRecorder()

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
	// geometry subpasses
	//

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
	})

	for _, gpass := range p.gpasses {
		gpass.Record(cmds, camera, scene)

		cmds.Record(func(cmd command.Buffer) {
			cmd.CmdNextSubpass()
		})
	}

	//
	// lighting pass
	//

	cmds.Record(func(cmd command.Buffer) {
		p.light.Bind(cmd)
	})

	lightDesc := p.light.Descriptors()
	lightDesc.Camera.Set(camera)

	ambient := light.Descriptor{
		Type:      light.Ambient,
		Color:     color.White,
		Intensity: 0.33,
	}
	p.DrawLight(cmds, args, ambient)

	lights := query.New[light.T]().Collect(scene)
	for _, lit := range lights {
		if err := p.DrawLight(cmds, args, lit.LightDescriptor()); err != nil {
			fmt.Printf("light draw error in object %s: %s\n", lit.Name(), err)
		}
	}

	//
	// todo: forward subpasses
	//

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})

	worker := p.backend.Worker(ctx.Index)
	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Signal: []sync.Semaphore{p.completed, p.copyReady},
		Wait: []command.Wait{
			{
				Semaphore: p.shadows.Completed(),
				Mask:      vk.PipelineStageFragmentShaderBit,
			},
		},
	})

	p.GeometryBuffer.CopyBuffers(worker, ctx.Index, command.Wait{
		Semaphore: p.copyReady,
		Mask:      vk.PipelineStageTopOfPipeBit,
	})
}

func (p *GeometryPass) DrawLight(cmds command.Recorder, args render.Args, lit light.Descriptor) error {
	vkmesh := p.meshes.Fetch(p.quad)
	cmds.Record(func(cmd command.Buffer) {

		push := LightConst{
			ViewProj:    lit.ViewProj,
			Color:       lit.Color,
			Position:    lit.Position,
			Type:        lit.Type,
			Shadowmap:   uint32(1),
			Range:       lit.Range,
			Intensity:   lit.Intensity,
			Attenuation: lit.Attenuation,
		}
		cmd.CmdPushConstant(vk.ShaderStageFragmentBit, 0, &push)

		vkmesh.Draw(cmd, 0)
	})

	return nil
}

func (p *GeometryPass) Destroy() {
	// destroy subpasses
	for _, gpass := range p.gpasses {
		gpass.Destroy()
	}

	p.pass.Destroy()
	p.GeometryBuffer.Destroy()
	p.light.Material().Destroy()
	p.completed.Destroy()
	p.copyReady.Destroy()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}
