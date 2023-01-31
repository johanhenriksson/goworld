package game

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer/pass"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
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
	target vulkan.Target
	mat    material.Standard
	shader string
}

func NewVoxelSubpass(target vulkan.Target) pass.DeferredSubpass {
	return &voxelpass{
		target: target,
		shader: "vk/color_f",
	}
}

func NewVoxelShadowpass(target vulkan.Target) pass.DeferredSubpass {
	return &voxelpass{
		target: target,
		shader: "vk/shadow",
	}
}

func (p *voxelpass) Instantiate(pool descriptor.Pool, rpass renderpass.T) {
	p.mat = material.FromDef(
		p.target.Device(),
		pool,
		rpass,
		&material.Def{
			Shader:       p.shader,
			Subpass:      p.Name(),
			VertexFormat: voxel.Vertex{},
			DepthTest:    true,
			DepthWrite:   true,
		})
}

func (p *voxelpass) Name() renderpass.Name {
	return "voxels"
}

func (p *voxelpass) Record(cmds command.Recorder, camera uniform.Camera, scene object.T) {
	p.mat.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		p.mat.Bind(cmd)
	})

	objects := object.Query[mesh.T]().Where(isVoxelMesh).Collect(scene)
	for index, mesh := range objects {
		if err := p.DrawDeferred(cmds, index, mesh, p.mat); err != nil {
			fmt.Printf("deferred draw error in object %s: %s\n", mesh.Name(), err)
		}
	}
}

func (p *voxelpass) DrawDeferred(cmds command.Recorder, index int, mesh mesh.T, mat material.Standard) error {
	vkmesh := p.target.Meshes().Fetch(mesh.Mesh())
	if vkmesh == nil {
		fmt.Println("mesh is nil")
		return nil
	}

	// write object properties to ssbo
	// todo: this should be reused between frames - maybe
	//       how to detect changes?
	mat.Descriptors().Objects.Set(index, uniform.Object{
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

func isVoxelMesh(m mesh.T) bool {
	// todo: improve this
	return m.Mode() == mesh.Deferred
}
