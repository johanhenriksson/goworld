package vkrender

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_shader"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_texture"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/shader"

	vk "github.com/vulkan-go/vulkan"
)

type DeferredPass interface {
	Pass
	Diffuse(frame int) vk_texture.T
}

type Uniforms struct {
	Projection mat4.T
	View       mat4.T
}

type Storage struct {
	Model mat4.T
}

type GeometryPass struct {
	meshes  cache.Meshes
	backend vulkan.T
	shader  vk_shader.T[game.VoxelVertex, Uniforms, Storage]

	diffuse     []image.T
	diffuseTex  []vk_texture.T
	depth       []image.T
	depthView   []image.View
	framebuffer []framebuffer.T
	pass        pipeline.Pass
	completed   sync.Semaphore
}

func NewGeometryPass(backend vulkan.T, meshes cache.Meshes) *GeometryPass {
	diffuseFmt := vk.FormatR16g16b16a16Sfloat

	passInfo := vk.RenderPassCreateInfo{
		SType:           vk.StructureTypeRenderPassCreateInfo,
		AttachmentCount: 2,
		PAttachments: []vk.AttachmentDescription{
			{
				Format:         diffuseFmt,
				Samples:        vk.SampleCount1Bit,
				LoadOp:         vk.AttachmentLoadOpClear,
				StoreOp:        vk.AttachmentStoreOpStore,
				StencilLoadOp:  vk.AttachmentLoadOpDontCare,
				StencilStoreOp: vk.AttachmentStoreOpDontCare,
				InitialLayout:  vk.ImageLayoutUndefined,
				FinalLayout:    vk.ImageLayoutShaderReadOnlyOptimal,
			},
			{
				Format:         backend.Device().GetDepthFormat(),
				Samples:        vk.SampleCount1Bit,
				LoadOp:         vk.AttachmentLoadOpClear,
				StoreOp:        vk.AttachmentStoreOpStore,
				StencilLoadOp:  vk.AttachmentLoadOpDontCare,
				StencilStoreOp: vk.AttachmentStoreOpStore,
				InitialLayout:  vk.ImageLayoutUndefined,
				FinalLayout:    vk.ImageLayoutDepthStencilAttachmentOptimal,
			},
		},
		SubpassCount: 1,
		PSubpasses: []vk.SubpassDescription{
			{
				PipelineBindPoint:    vk.PipelineBindPointGraphics,
				InputAttachmentCount: 0,
				ColorAttachmentCount: 1,
				PColorAttachments: []vk.AttachmentReference{
					{
						Attachment: 0,
						Layout:     vk.ImageLayoutColorAttachmentOptimal,
					},
				},
				PDepthStencilAttachment: &vk.AttachmentReference{
					Attachment: 1,
					Layout:     vk.ImageLayoutDepthStencilAttachmentOptimal,
				},
			},
		},
		DependencyCount: 0,
	}
	pass := pipeline.NewPass(backend.Device(), &passInfo)

	framebuffers := make([]framebuffer.T, backend.Frames())
	diffuse := make([]image.T, backend.Frames())
	diffuseTex := make([]vk_texture.T, backend.Frames())
	depth := make([]image.T, backend.Frames())
	depthView := make([]image.View, backend.Frames())

	depthFormat := backend.Device().GetDepthFormat()
	depthUsage := vk.ImageUsageFlags(vk.ImageUsageDepthStencilAttachmentBit | vk.ImageUsageSampledBit)

	for i := 0; i < backend.Frames(); i++ {
		diffuse[i] = image.New2D(backend.Device(), backend.Width(), backend.Height(), diffuseFmt, vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit|vk.ImageUsageSampledBit))
		diffuseTex[i] = vk_texture.FromImage(backend.Device(), diffuse[i], vk_texture.Args{
			Format: diffuseFmt,
			Filter: vk.FilterLinear,
			Wrap:   vk.SamplerAddressModeRepeat,
		})
		depth[i] = image.New2D(backend.Device(), backend.Width(), backend.Height(), depthFormat, depthUsage)
		depthView[i] = depth[i].View(depthFormat, vk.ImageAspectFlags(vk.ImageAspectDepthBit|vk.ImageAspectStencilBit))
		framebuffers[i] = framebuffer.New(backend.Device(), backend.Width(), backend.Height(), pass.Ptr(), []image.View{
			diffuseTex[i].View(),
			depthView[i],
		})
	}

	sh := vk_shader.New[game.VoxelVertex, Uniforms, Storage](backend, vk_shader.Args{
		Path: "vk/color_f",
		Pass: pass,
		Attributes: shader.AttributeMap{
			"position": {
				Bind: 0,
				Type: types.Float,
			},
			"color_0": {
				Bind: 1,
				Type: types.Float,
			},
		},
	})
	return &GeometryPass{
		backend:   backend,
		meshes:    meshes,
		shader:    sh,
		completed: sync.NewSemaphore(backend.Device()),

		pass:        pass,
		diffuse:     diffuse,
		diffuseTex:  diffuseTex,
		depth:       depth,
		depthView:   depthView,
		framebuffer: framebuffers,
	}
}

func (p *GeometryPass) Completed() sync.Semaphore {
	return p.completed
}

func (p *GeometryPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	cmds := command.NewRecorder()

	p.shader.SetUniforms(ctx.Index, []Uniforms{
		{
			Projection: args.Projection,
			View:       args.View,
		},
	})

	p.shader.SetStorage(ctx.Index, []Storage{
		{
			Model: mat4.Ident(),
		},
		{
			Model: mat4.Translate(vec3.New(-16, 0, 0)),
		},
	})

	cmds.Record(func(cmd command.Buffer) {
		clear := color.RGB(0.1, 0.1, 0.3)

		cmd.CmdBeginRenderPass(p.pass, p.framebuffer[ctx.Index], clear)
		cmd.CmdSetViewport(0, 0, ctx.Width, ctx.Height)
		cmd.CmdSetScissor(0, 0, ctx.Width, ctx.Height)

		p.shader.Bind(ctx.Index, cmd)
	})

	objects := query.New[mesh.T]().Where(isDrawDeferred).Collect(scene)
	for _, mesh := range objects {
		if err := p.DrawDeferred(cmds, args, mesh); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", mesh.Name(), err)
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})

	worker := ctx.Workers[0]
	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Wait:   []sync.Semaphore{ctx.ImageAvailable},
		Signal: []sync.Semaphore{p.completed},
		WaitMask: []vk.PipelineStageFlags{
			vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		},
	})
}

func (p *GeometryPass) DrawDeferred(cmds command.Recorder, args render.Args, mesh mesh.T) error {
	args = args.Apply(mesh.Transform().World())

	vkmesh, ok := p.meshes.Fetch(mesh.Mesh(), nil).(*cache.VkMesh)
	if !ok {
		fmt.Println("mesh is nil")
		return nil
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBindVertexBuffer(vkmesh.Vertices, 0)
		cmd.CmdBindIndexBuffers(vkmesh.Indices, 0, vk.IndexTypeUint16)
		// cmd.CmdDraw(vkmesh.Mesh.Elements(), vkmesh.Mesh.Elements()/3, 0, 0)
		// cmd.CmdPushConstant(p.shader.Layout(), vk.ShaderStageFlags(vk.ShaderStageAll), 0, &push)

		// index of the object properties in the ssbo
		idx := 0
		cmd.CmdDrawIndexed(vkmesh.Mesh.Elements(), 1, 0, 0, idx)

		idx = 1
		cmd.CmdDrawIndexed(vkmesh.Mesh.Elements(), 1, 0, 0, idx)
	})

	return nil
}

func (p *GeometryPass) Diffuse(frame int) vk_texture.T {
	return p.diffuseTex[frame]
}

func (p *GeometryPass) Destroy() {
	p.pass.Destroy()
	for i := 0; i < p.backend.Frames(); i++ {
		p.framebuffer[i].Destroy()
		p.diffuseTex[i].Destroy()
		p.depth[i].Destroy()
		p.depthView[i].Destroy()
	}
	p.shader.Destroy()
	p.completed.Destroy()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}
