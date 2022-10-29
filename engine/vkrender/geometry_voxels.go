package vkrender

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/cache"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/types"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

//
// each material requires its own subpass
// so - it would be better to look up all materials in the scene, and dynamically
// create a subpass for each material. the subpass would then query all meshes
// with its assigned material, and render them.
//

type voxelpass struct {
	backend vulkan.T
	mat     material.Instance[*GeometryDescriptors]
	meshes  cache.MeshCache
}

func NewVoxelSubpass(backend vulkan.T, meshes cache.MeshCache) DeferredSubpass {
	return &voxelpass{
		backend: backend,
		meshes:  meshes,
	}
}

func (p *voxelpass) Instantiate(pass renderpass.T) {
	p.mat = material.New(
		p.backend.Device(),
		material.Args{
			Shader: shader.New(
				p.backend.Device(),
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
					"Camera":   0,
					"Objects":  1,
					"Textures": 2,
				},
			),
			Pass:       pass,
			Subpass:    p.Name(),
			Pointers:   vertex.ParsePointers(game.VoxelVertex{}),
			DepthTest:  true,
			DepthWrite: true,
		},
		&GeometryDescriptors{
			Camera: &descriptor.Uniform[Camera]{
				Stages: vk.ShaderStageAll,
			},
			Objects: &descriptor.Storage[ObjectStorage]{
				Stages: vk.ShaderStageAll,
				Size:   10,
			},
			Textures: &descriptor.SamplerArray{
				Stages: vk.ShaderStageFragmentBit,
				Count:  1,
			},
		}).Instantiate()
}

func (p *voxelpass) Name() string {
	return "voxels"
}

func (p *voxelpass) Record(cmds command.Recorder, camera Camera, scene object.T) {
	p.mat.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		p.mat.Bind(cmd)
	})

	objects := query.New[mesh.T]().Where(isDrawDeferred).Collect(scene)
	for index, mesh := range objects {
		if err := p.DrawDeferred(cmds, index, mesh); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", mesh.Name(), err)
		}
	}
}

func (p *voxelpass) DrawDeferred(cmds command.Recorder, index int, mesh mesh.T) error {
	vkmesh := p.meshes.Fetch(mesh.Mesh())
	if vkmesh == nil {
		fmt.Println("mesh is nil")
		return nil
	}

	// write object properties to ssbo
	// todo: this should be reused between frames - maybe
	//       how to detect changes?
	p.mat.Descriptors().Objects.Set(index, ObjectStorage{
		Model: mesh.Transform().World(),
	})

	cmds.Record(func(cmd command.Buffer) {
		vkmesh.Draw(cmd, index)
	})

	return nil
}

func (p *voxelpass) Destroy() {
	p.mat.Material().Destroy()
}
