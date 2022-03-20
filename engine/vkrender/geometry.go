package vkrender

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/texture"
	"github.com/johanhenriksson/goworld/render/color"
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

	meshes    MeshCache
	quad      vertex.Mesh
	backend   vulkan.T
	pass      renderpass.T
	geom      material.Instance[*GeometryDescriptors]
	light     material.Instance[*LightDescriptors]
	completed sync.Semaphore

	shadows ShadowPass
}

func NewGeometryPass(backend vulkan.T, meshes MeshCache, shadows ShadowPass) *GeometryPass {
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
				StoreOp:     vk.AttachmentStoreOpStore,
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
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit,
			},
			{
				Name:        "position",
				Format:      positionFmt,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit,
			},
		},
		DepthAttachment: &renderpass.DepthAttachment{
			LoadOp:      vk.AttachmentLoadOpClear,
			StoreOp:     vk.AttachmentStoreOpStore,
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

	geomsh := material.New(
		backend.Device(),
		material.Args{
			Shader: shader.New(
				backend.Device(),
				"vk/color_f",
				shader.Inputs{
					"position": {
						Index: 0,
						Type:  types.Float,
					},
					"normal_id": {
						Index: 1,
						Type:  types.UInt8,
					},
					"color_0": {
						Index: 2,
						Type:  types.Float,
					},
					"occlusion": {
						Index: 3,
						Type:  types.Float,
					},
				},
				shader.Descriptors{
					"Camera":  0,
					"Objects": 1,
				},
			),
			Pass:     pass,
			Subpass:  "geometry",
			Pointers: vertex.ParsePointers(game.VoxelVertex{}),
		},
		&GeometryDescriptors{
			Camera: &descriptor.Uniform[CameraData]{
				Stages: vk.ShaderStageAll,
			},
			Objects: &descriptor.Storage[ObjectStorage]{
				Stages: vk.ShaderStageAll,
				Size:   10,
			},
		}).Instantiate()

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

	shadowtex := texture.FromView(backend.Device(), shadows.Shadowmap(), texture.Args{
		Filter: vk.FilterNearest,
		Wrap:   vk.SamplerAddressModeClampToEdge,
	})
	lightsh.Descriptors().Shadow.Set(shadowtex)

	return &GeometryPass{
		GeometryBuffer: gbuffer,

		backend:   backend,
		meshes:    meshes,
		quad:      quad,
		geom:      geomsh,
		light:     lightsh,
		pass:      pass,
		completed: sync.NewSemaphore(backend.Device()),

		shadows: shadows,
	}
}

func (p *GeometryPass) Completed() sync.Semaphore {
	return p.completed
}

func (p *GeometryPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	cmds := command.NewRecorder()

	geomDesc := p.geom.Descriptors()

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

	lightDesc := p.light.Descriptors()

	lightDesc.Camera.Set(camera)

	lightDesc.Diffuse.Set(p.GeometryBuffer.Diffuse(ctx.Index))
	lightDesc.Normal.Set(p.GeometryBuffer.Normal(ctx.Index))
	lightDesc.Position.Set(p.GeometryBuffer.Position(ctx.Index))
	lightDesc.Depth.Set(p.GeometryBuffer.Depth(ctx.Index))

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
		p.geom.Bind(cmd)
	})

	objects := query.New[mesh.T]().Where(isDrawDeferred).Collect(scene)
	for _, mesh := range objects {
		if err := p.DrawDeferred(cmds, args, mesh); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", mesh.Name(), err)
		}
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdNextSubpass()
		p.light.Bind(cmd)
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
		Signal: []sync.Semaphore{p.completed},
		Wait: []command.Wait{
			{
				Semaphore: p.shadows.Completed(),
				Mask:      vk.PipelineStageFragmentShaderBit,
			},
		},
	})
	// worker.Wait()
}

func (p *GeometryPass) DrawDeferred(cmds command.Recorder, args render.Args, mesh mesh.T) error {
	args = args.Apply(mesh.Transform().World())

	vkmesh := p.meshes.Fetch(mesh.Mesh())
	if vkmesh == nil {
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
	p.light.Descriptors().Light.Set(index, lit)

	vkmesh := p.meshes.Fetch(p.quad)
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBindVertexBuffer(vkmesh.Vertices, 0)
		cmd.CmdBindIndexBuffers(vkmesh.Indices, 0, vk.IndexTypeUint16)

		push := LightConst{
			LightIndex: uint32(index),
		}
		cmd.CmdPushConstant(p.light.Material().Layout(), vk.ShaderStageFlags(vk.ShaderStageFragmentBit), 0, &push)

		cmd.CmdDrawIndexed(vkmesh.Mesh.Elements(), 1, 0, 0, index)
	})

	return nil
}

func (p *GeometryPass) Destroy() {
	p.pass.Destroy()
	p.GeometryBuffer.Destroy()
	p.geom.Material().Destroy()
	p.light.Material().Destroy()
	p.completed.Destroy()
}

func isDrawDeferred(m mesh.T) bool {
	return m.Mode() == mesh.Deferred
}
