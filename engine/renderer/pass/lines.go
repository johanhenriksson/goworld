package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type LinePass struct {
	target    vulkan.Target
	material  material.Instance[*LineDescriptors]
	pass      renderpass.T
	completed sync.Semaphore
	fbufs     framebuffer.Array
	wait      sync.Semaphore
}

type LineDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[uniform.Camera]
	Objects *descriptor.Storage[mat4.T]
}

func NewLinePass(target vulkan.Target, pool descriptor.Pool, output Pass, geometry DeferredPass, wait sync.Semaphore) *LinePass {
	log.Println("create line pass")

	p := &LinePass{
		target:    target,
		wait:      wait,
		completed: sync.NewSemaphore(target.Device()),
	}

	depth := make([]image.T, target.Frames())
	for i := range depth {
		depth[i] = geometry.Depth().Image()
	}

	p.pass = renderpass.New(target.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:          "color",
				Allocator:     attachment.FromSwapchain(target.Swapchain()),
				Format:        target.Swapchain().SurfaceFormat(),
				LoadOp:        vk.AttachmentLoadOpLoad,
				StoreOp:       vk.AttachmentStoreOpStore,
				InitialLayout: vk.ImageLayoutPresentSrc,
				FinalLayout:   vk.ImageLayoutPresentSrc,
			},
		},
		DepthAttachment: &attachment.Depth{
			Allocator:     attachment.FromImageArray(depth),
			LoadOp:        vk.AttachmentLoadOpLoad,
			InitialLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			FinalLayout:   vk.ImageLayoutDepthStencilAttachmentOptimal,
			Usage:         vk.ImageUsageInputAttachmentBit,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  "output",
				Depth: true,

				ColorAttachments: []attachment.Name{"color"},
			},
		},
	})

	var err error
	p.fbufs, err = framebuffer.NewArray(target.Frames(), target.Device(), target.Width(), target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.material = material.New(
		target.Device(),
		material.Args{
			Shader:    shader.New(target.Device(), "vk/lines"),
			Pass:      p.pass,
			Pointers:  vertex.ParsePointers(vertex.C{}),
			Primitive: vertex.Lines,
			DepthTest: true,
		},
		&LineDescriptors{
			Camera: &descriptor.Uniform[uniform.Camera]{
				Stages: vk.ShaderStageVertexBit,
			},
			Objects: &descriptor.Storage[mat4.T]{
				Size:   100,
				Stages: vk.ShaderStageVertexBit,
			},
		}).Instantiate(pool)

	return p
}

func (p *LinePass) Draw(args render.Args, scene object.T) {
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
	p.material.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[ctx.Index])
		p.material.Bind(cmd)
	})

	objects := query.New[mesh.T]().Where(isDrawLines).Collect(scene)
	for i, mesh := range objects {
		p.DrawLines(cmds, i, args, mesh)
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})

	worker := p.target.Worker(ctx.Index)
	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Marker: "LinePass",
		Signal: []sync.Semaphore{p.completed},
		Wait: []command.Wait{
			{
				Semaphore: p.wait,
				Mask:      vk.PipelineStageColorAttachmentOutputBit,
			},
		},
	})
}

func (p *LinePass) DrawLines(cmds command.Recorder, index int, args render.Args, mesh mesh.T) error {
	vkmesh := p.target.Meshes().Fetch(mesh.Mesh())
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

func (p *LinePass) Name() string {
	return "Lines"
}

func (p *LinePass) Completed() sync.Semaphore {
	return p.completed
}

func (p *LinePass) Destroy() {
	p.completed.Destroy()
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.material.Material().Destroy()
}

func isDrawLines(m mesh.T) bool {
	return m.Mode() == mesh.Lines
}
