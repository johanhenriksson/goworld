package vkrender

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_shader"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type DeferredPass interface {
	Pass
	GeometryBuffer
}

type CameraData struct {
	Proj        mat4.T
	View        mat4.T
	ViewProj    mat4.T
	ProjInv     mat4.T
	ViewInv     mat4.T
	ViewProjInv mat4.T
	Eye         vec3.T
}

type ObjectStorage struct {
	Model mat4.T
}

type GeometryDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[CameraData]
	Objects *descriptor.Storage[ObjectStorage]
}

type LightConst struct {
	LightIndex uint32
}

type GeometryPass struct {
	GeometryBuffer

	meshes    cache.Meshes
	quad      vertex.Mesh
	backend   vulkan.T
	pass      renderpass.T
	geom      vk_shader.T[*GeometryDescriptors]
	light     vk_shader.T[*LightDescriptors]
	completed sync.Semaphore
}

func NewGeometryPass(backend vulkan.T, meshes cache.Meshes) *GeometryPass {
	diffuseFmt := vk.FormatR8g8b8a8Unorm
	normalFmt := vk.FormatR8g8b8a8Unorm
	positionFmt := vk.FormatR16g16b16a16Sfloat

	pass := renderpass.New(backend.Device(), renderpass.Args{
		Frames: backend.Frames(),
		Width:  backend.Width(),
		Height: backend.Height(),

		ColorAttachments: []renderpass.ColorAttachment{
			{
				Name:        "output",
				Format:      diffuseFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageSampledBit,
				Blend:       true,
			},
			{
				Name:        "diffuse",
				Format:      diffuseFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageSampledBit,
			},
			{
				Name:        "normal",
				Format:      normalFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpDontCare,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit,
			},
			{
				Name:        "position",
				Format:      positionFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpDontCare,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit,
			},
		},
		DepthAttachment: &renderpass.DepthAttachment{
			LoadOp:      vk.AttachmentLoadOpClear,
			StoreOp:     vk.AttachmentStoreOpDontCare,
			FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			Usage:       vk.ImageUsageInputAttachmentBit,
			ClearDepth:  1,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  "geometry",
				Depth: true,

				ColorAttachments: []string{"diffuse", "normal", "position"},
			},
			{
				Name:  "lighting",
				Depth: false,

				ColorAttachments: []string{"output"},
				InputAttachments: []string{"diffuse", "normal", "position", "depth"},
			},
		},
		Dependencies: []renderpass.SubpassDependency{
			{
				Src: "external",
				Dst: "geometry",

				SrcStageMask:  vk.PipelineStageBottomOfPipeBit,
				DstStageMask:  vk.PipelineStageColorAttachmentOutputBit,
				SrcAccessMask: vk.AccessMemoryReadBit,
				DstAccessMask: vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit,
				Flags:         vk.DependencyByRegionBit,
			},
			{
				Src: "geometry",
				Dst: "lighting",

				SrcStageMask:  vk.PipelineStageColorAttachmentOutputBit,
				DstStageMask:  vk.PipelineStageFragmentShaderBit,
				SrcAccessMask: vk.AccessColorAttachmentWriteBit,
				DstAccessMask: vk.AccessShaderReadBit,
				Flags:         vk.DependencyByRegionBit,
			},
			{
				Src: "geometry",
				Dst: "external",

				SrcStageMask:  vk.PipelineStageColorAttachmentOutputBit,
				DstStageMask:  vk.PipelineStageBottomOfPipeBit,
				SrcAccessMask: vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit,
				DstAccessMask: vk.AccessMemoryReadBit,
				Flags:         vk.DependencyByRegionBit,
			},
		},
	})

	gbuffer := NewGbuffer(backend, pass, backend.Frames())

	geomsh := vk_shader.New(
		backend,
		vk_shader.Args{
			Path:     "vk/color_f",
			Frames:   1,
			Pass:     pass,
			Subpass:  "geometry",
			Pointers: vertex.ParsePointers(game.VoxelVertex{}),
			Attributes: shader.AttributeMap{
				"position": {
					Loc:  0,
					Type: types.Float,
				},
				"normal_id": {
					Loc:  1,
					Type: types.UInt8,
				},
				"color_0": {
					Loc:  2,
					Type: types.Float,
				},
				"occlusion": {
					Loc:  3,
					Type: types.Float,
				},
			},
		},
		&GeometryDescriptors{
			Camera: &descriptor.Uniform[CameraData]{
				Binding: 0,
				Stages:  vk.ShaderStageAll,
			},
			Objects: &descriptor.Storage[ObjectStorage]{
				Binding: 1,
				Stages:  vk.ShaderStageAll,
				Size:    10,
			},
		})

	quad := vertex.NewTriangles("screen_quad", []vertex.T{
		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
		{P: vec3.New(-1, 1, 0), T: vec2.New(0, 1)},
		{P: vec3.New(1, -1, 0), T: vec2.New(1, 0)},
	}, []uint16{
		0, 1, 2,
		0, 3, 1,
	})

	lightsh := vk_shader.New(
		backend,
		vk_shader.Args{
			Path:     "vk/light",
			Frames:   1,
			Pass:     pass,
			Subpass:  "lighting",
			Pointers: vertex.ParsePointers(vertex.T{}),
			Attributes: shader.AttributeMap{
				"position": {
					Loc:  0,
					Type: types.Float,
				},
			},
			Constants: []pipeline.PushConstant{
				{
					Stages: vk.ShaderStageFragmentBit,
					Offset: 0,
					Size:   4,
				},
			},
		},
		&LightDescriptors{
			Diffuse: &descriptor.InputAttachment{
				Binding: 0,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Normal: &descriptor.InputAttachment{
				Binding: 1,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Position: &descriptor.InputAttachment{
				Binding: 2,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Depth: &descriptor.InputAttachment{
				Binding: 3,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Camera: &descriptor.Uniform[CameraData]{
				Binding: 4,
				Stages:  vk.ShaderStageFragmentBit,
			},
			Light: &descriptor.UniformArray[light.Descriptor]{
				Binding: 5,
				Size:    10,
				Stages:  vk.ShaderStageFragmentBit,
			},
		})

	// lightsh.Descriptors(0).Depth.Set(gbuffer.Depth(0))

	return &GeometryPass{
		GeometryBuffer: gbuffer,

		backend:   backend,
		meshes:    meshes,
		quad:      quad,
		geom:      geomsh,
		light:     lightsh,
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

	geomDesc := p.geom.Descriptors(ctx.Index)

	camera := CameraData{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         args.Position,
	}

	geomDesc.Camera.Set(camera)

	geomDesc.Objects.Set(0, ObjectStorage{
		Model: mat4.Ident(),
	})

	geomDesc.Objects.Set(1, ObjectStorage{
		Model: mat4.Translate(vec3.New(-16, 0, 0)),
	})

	lightDesc := p.light.Descriptors(ctx.Index)

	lightDesc.Camera.Set(camera)

	lightDesc.Diffuse.Set(p.GeometryBuffer.Diffuse(ctx.Index))
	lightDesc.Normal.Set(p.GeometryBuffer.Normal(ctx.Index))
	lightDesc.Position.Set(p.GeometryBuffer.Position(ctx.Index))
	lightDesc.Depth.Set(p.GeometryBuffer.Depth(ctx.Index))

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
		p.geom.Bind(ctx.Index, cmd)
	})

	objects := query.New[mesh.T]().Where(isDrawDeferred).Collect(scene)
	for _, mesh := range objects {
		if err := p.DrawDeferred(cmds, args, mesh); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", mesh.Name(), err)
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdNextSubpass()
		p.light.Bind(ctx.Index, cmd)
	})

	ambient := light.Descriptor{
		Type:      light.Ambient,
		Color:     color.White,
		Intensity: 0.33,
	}
	p.DrawLight(cmds, args, 0, ambient)

	lights := query.New[light.T]().Collect(scene)
	for i, lit := range lights {
		if err := p.DrawLight(cmds, args, i+1, lit.LightDescriptor()); err != nil {
			fmt.Printf("light draw error in object %s: %s\n", lit.Name(), err)
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})

	worker := p.backend.Worker(ctx.Index)
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

func (p *GeometryPass) DrawLight(cmds command.Recorder, args render.Args, index int, lit light.Descriptor) error {
	p.light.Descriptors(0).Light.Set(index, lit)

	vkmesh := p.meshes.Fetch(p.quad, nil).(*cache.VkMesh)
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBindVertexBuffer(vkmesh.Vertices, 0)
		cmd.CmdBindIndexBuffers(vkmesh.Indices, 0, vk.IndexTypeUint16)

		push := LightConst{
			LightIndex: uint32(index),
		}
		cmd.CmdPushConstant(p.light.Layout(), vk.ShaderStageFlags(vk.ShaderStageFragmentBit), 0, &push)

		cmd.CmdDrawIndexed(vkmesh.Mesh.Elements(), 1, 0, 0, index)
	})

	return nil
}

func (p *GeometryPass) Destroy() {
	p.pass.Destroy()
	p.GeometryBuffer.Destroy()
	p.geom.Destroy()
	p.light.Destroy()
	p.completed.Destroy()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}
