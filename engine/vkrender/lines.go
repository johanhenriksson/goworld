package vkrender

import (
	"log"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/types"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type LinePass struct {
	backend   vulkan.T
	meshes    cache.MeshCache
	material  material.Instance[*LineDescriptors]
	geometry  DeferredPass
	pass      renderpass.T
	output    Pass
	completed sync.Semaphore

	shadows ShadowPass
}

type LineDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[Camera]
	Objects *descriptor.Storage[mat4.T]
}

func NewLinePass(backend vulkan.T, meshes cache.MeshCache, output Pass, geometry DeferredPass) *LinePass {
	log.Println("create line pass")

	p := &LinePass{
		backend:   backend,
		meshes:    meshes,
		geometry:  geometry,
		output:    output,
		completed: sync.NewSemaphore(backend.Device()),
	}

	depth := make([]image.T, backend.Frames())
	for i := range depth {
		depth[i] = geometry.Depth(i).Image()
	}

	p.pass = renderpass.New(backend.Device(), renderpass.Args{
		Frames: backend.Frames(),
		Width:  backend.Width(),
		Height: backend.Height(),
		ColorAttachments: []renderpass.ColorAttachment{
			{
				Name:          "color",
				Images:        backend.Swapchain().Images(),
				Format:        backend.Swapchain().SurfaceFormat(),
				LoadOp:        vk.AttachmentLoadOpLoad,
				StoreOp:       vk.AttachmentStoreOpStore,
				InitialLayout: vk.ImageLayoutPresentSrc,
				FinalLayout:   vk.ImageLayoutPresentSrc,
			},
		},
		DepthAttachment: &renderpass.DepthAttachment{
			Images:        depth,
			LoadOp:        vk.AttachmentLoadOpLoad,
			InitialLayout: vk.ImageLayoutUndefined,
			FinalLayout:   vk.ImageLayoutDepthStencilAttachmentOptimal,
			Usage:         vk.ImageUsageInputAttachmentBit,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  "output",
				Depth: true,

				ColorAttachments: []string{"color"},
			},
		},
	})

	p.material = material.New(
		backend.Device(),
		material.Args{
			Shader: shader.New(
				backend.Device(),
				"vk/lines",
				shader.Inputs{
					"position": {
						Index: 0,
						Type:  types.Float,
					},
					"color_0": {
						Index: 1,
						Type:  types.Float,
					},
				},
				shader.Descriptors{
					"Camera":  0,
					"Objects": 1,
				},
			),
			Pass:      p.pass,
			Pointers:  vertex.ParsePointers(vertex.C{}),
			Primitive: vertex.Lines,
			DepthTest: true,
		},
		&LineDescriptors{
			Camera: &descriptor.Uniform[Camera]{
				Stages: vk.ShaderStageVertexBit,
			},
			Objects: &descriptor.Storage[mat4.T]{
				Size:   100,
				Stages: vk.ShaderStageVertexBit,
			},
		}).Instantiate()

	return p
}

func (p *LinePass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	cmds := command.NewRecorder()

	camera := Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         args.Position,
	}
	p.material.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
		p.material.Bind(cmd)
	})

	objects := query.New[mesh.T]().Where(isDrawLines).Collect(scene)
	for i, mesh := range objects {
		p.DrawLines(cmds, i, args, mesh)
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})

	worker := p.backend.Worker(ctx.Index)
	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Signal: []sync.Semaphore{p.completed},
		Wait: []command.Wait{
			{
				Semaphore: p.output.Completed(),
				Mask:      vk.PipelineStageColorAttachmentOutputBit,
			},
		},
	})
}

func (p *LinePass) DrawLines(cmds command.Recorder, index int, args render.Args, mesh mesh.T) error {
	args = args.Apply(mesh.Transform().World())

	vkmesh := p.meshes.Fetch(mesh.Mesh())
	if vkmesh == nil {
		log.Println("line mesh", mesh.Mesh().Id(), "is nil")
		return nil
	}

	cmds.Record(func(cmd command.Buffer) {
		p.material.Descriptors().Objects.Set(index, mesh.Transform().World())

		vkmesh.Draw(cmd, index)
	})

	return nil
}

func (p *LinePass) Completed() sync.Semaphore {
	return p.completed
}

func (p *LinePass) Destroy() {
	p.completed.Destroy()
	p.pass.Destroy()
	p.material.Material().Destroy()
}

func isDrawLines(m mesh.T) bool {
	return m.Mode() == mesh.Lines
}
