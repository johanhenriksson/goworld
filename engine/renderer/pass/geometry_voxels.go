package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

//
// each material requires its own subpass
// so - it would be better to look up all materials in the scene, and dynamically
// create a subpass for each material. the subpass would then query all meshes
// with its assigned material, and render them.
//

type voxelpass struct {
	backend vulkan.T
	mat     material.Standard
	shadows material.Standard
	meshes  cache.MeshCache
}

func NewVoxelSubpass(backend vulkan.T, meshes cache.MeshCache) DeferredSubpass {
	return &voxelpass{
		backend: backend,
		meshes:  meshes,
	}
}

func (p *voxelpass) Instantiate(pass renderpass.T) {
	p.mat = material.FromDef(
		p.backend,
		pass,
		&material.Def{
			Shader:       "vk/color_f",
			VertexFormat: game.VoxelVertex{},
		})
}

func (p *voxelpass) InstantiateShadows(pass renderpass.T) {
	p.shadows = material.FromDef(
		p.backend,
		pass,
		&material.Def{
			Shader:       "vk/shadow",
			VertexFormat: game.VoxelVertex{},
		})
}

func (p *voxelpass) Name() string {
	return "vk/color_f"
}

func (p *voxelpass) Record(cmds command.Recorder, camera uniform.Camera, scene object.T) {
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

func (p *voxelpass) RecordShadows(cmds command.Recorder, camera uniform.Camera, scene object.T) {
	p.mat.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		p.shadows.Bind(cmd)
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
	p.mat.Descriptors().Objects.Set(index, uniform.Object{
		Model: mesh.Transform().World(),
	})

	cmds.Record(func(cmd command.Buffer) {
		vkmesh.Draw(cmd, index)
	})

	return nil
}

func (p *voxelpass) Destroy() {
	p.mat.Material().Destroy()
	p.shadows.Material().Destroy()
}
