package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/texture"
)

type DeferredDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type DeferredMaterial struct {
	Pipeline    *pipeline.Pipeline
	Descriptors *DeferredDescriptors
	Objects     *ObjectBuffer
	textures    cache.SamplerCache
	Meshes      cache.MeshCache
	Commands    *command.IndirectDrawBuffer

	id    material.ID
	slots []texture.Slot
}

func (m *DeferredMaterial) ID() material.ID          { return m.id }
func (m *DeferredMaterial) Textures() []texture.Slot { return m.slots }

func (m *DeferredMaterial) Begin(camera uniform.Camera) {
	m.Descriptors.Camera.Set(camera)
	// todo: assign global descriptors

	m.Objects.Reset()
	m.Commands.Reset()
}

func (m *DeferredMaterial) Bind(cmds command.Recorder) {
	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBindGraphicsPipeline(m.Pipeline)
		cmd.CmdBindGraphicsDescriptor(m.Pipeline.Layout(), 0, m.Descriptors)
		m.Commands.BeginDrawIndirect()
	})
}

func (m *DeferredMaterial) Draw(cmds command.Recorder, msh mesh.Mesh) {
	gpuMesh, meshReady := m.Meshes.TryFetch(msh.Mesh())
	if !meshReady {
		return
	}

	textureIds := AssignMeshTextures(m.textures, msh, m.slots)

	instanceId := m.Objects.Store(uniform.Object{
		Model:    msh.Transform().Matrix(),
		Textures: textureIds,
	})

	// this is really the only thing that should be in the draw call
	cmds.Record(func(cmd *command.Buffer) {
		gpuMesh.Bind(cmd)
		gpuMesh.Draw(m.Commands, instanceId)
	})
}

func (m *DeferredMaterial) Unbind(cmds command.Recorder) {
	cmds.Record(func(cmd *command.Buffer) {
		m.Commands.EndDrawIndirect(cmd)
	})
}

func (m *DeferredMaterial) End() {
	m.Objects.Flush(m.Descriptors.Objects)
	m.textures.Flush(m.Descriptors.Textures)
	m.Commands.Flush()
}

func (m *DeferredMaterial) Destroy() {
	m.Descriptors.Destroy()
	m.Pipeline.Destroy()
	m.Commands.Destroy()
}
