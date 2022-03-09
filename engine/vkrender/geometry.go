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
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_shader"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/shader"

	vk "github.com/vulkan-go/vulkan"
)

type DeferredPass interface {
	Pass
	GeometryBuffer
}

type CameraData struct {
	Projection mat4.T
	View       mat4.T
}

type ObjectStorage struct {
	Model mat4.T
}

type GeometryDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[CameraData]
	Objects *descriptor.Storage[ObjectStorage]
}

type GeometryPass struct {
	GeometryBuffer

	meshes    cache.Meshes
	backend   vulkan.T
	pass      renderpass.T
	shader    vk_shader.T[game.VoxelVertex, *GeometryDescriptors]
	completed sync.Semaphore
}

func NewGeometryPass(backend vulkan.T, meshes cache.Meshes) *GeometryPass {
	diffuseFmt := vk.FormatR16g16b16a16Sfloat
	normalFmt := vk.FormatR8g8b8a8Unorm
	positionFmt := vk.FormatR16g16b16a16Sfloat

	pass := renderpass.New(backend.Device(), renderpass.Args{
		Frames: backend.Frames(),
		Width:  backend.Width(),
		Height: backend.Height(),

		ColorAttachments: []renderpass.ColorAttachment{
			{
				Name:        "diffuse",
				Format:      diffuseFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Clear:       color.RGB(0.1, 0.1, 0.16),
			},
			{
				Name:        "normal",
				Format:      normalFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			},
			{
				Name:        "position",
				Format:      positionFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			},
		},
		DepthAttachment: &renderpass.DepthAttachment{
			LoadOp:       vk.AttachmentLoadOpClear,
			FinalLayout:  vk.ImageLayoutDepthStencilAttachmentOptimal,
			ClearDepth:   1,
			ClearStencil: 0,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  "geometry",
				Depth: true,

				ColorAttachments: []string{"diffuse", "normal", "position"},
			},
		},
		Dependencies: []renderpass.SubpassDependency{},
	})

	gbuffer := NewGbuffer(backend, pass)

	sh := vk_shader.New[game.VoxelVertex](
		backend,
		&GeometryDescriptors{
			Camera: &descriptor.Uniform[CameraData]{
				Binding: 0,
				Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
			},
			Objects: &descriptor.Storage[ObjectStorage]{
				Binding: 1,
				Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
				Size:    100,
			},
		},
		vk_shader.Args{
			Path: "vk/color_f",
			Pass: pass,
			Attributes: shader.AttributeMap{
				"position": {
					Bind: 0,
					Type: types.Float,
				},
				"normal_id": {
					Bind: 1,
					Type: types.UInt8,
				},
				"color_0": {
					Bind: 2,
					Type: types.Float,
				},
			},
		})

	return &GeometryPass{
		GeometryBuffer: gbuffer,

		backend:   backend,
		meshes:    meshes,
		shader:    sh,
		pass:      pass,
		completed: sync.NewSemaphore(backend.Device()),
	}
}

func (p *GeometryPass) Completed() sync.Semaphore {
	return p.completed
}

func (p *GeometryPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	cmds := command.NewRecorder()

	descriptors := p.shader.Descriptors(ctx.Index)

	descriptors.Camera.Set(CameraData{
		Projection: args.Projection,
		View:       args.View,
	})

	descriptors.Objects.Set(0, ObjectStorage{
		Model: mat4.Ident(),
	})

	descriptors.Objects.Set(1, ObjectStorage{
		Model: mat4.Translate(vec3.New(-16, 0, 0)),
	})

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
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

		// index of the object properties in the ssbo
		idx := 0
		count := 2
		cmd.CmdDrawIndexed(vkmesh.Mesh.Elements(), count, 0, 0, idx)
	})

	return nil
}

func (p *GeometryPass) Destroy() {
	p.pass.Destroy()
	p.GeometryBuffer.Destroy()
	p.shader.Destroy()
	p.completed.Destroy()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}
