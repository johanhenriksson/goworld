package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
)

type DeferredDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type DeferredMaterial struct {
	Instance *material.Instance[*DeferredDescriptors]
	Objects  *ObjectBuffer
	Textures cache.SamplerCache
	Meshes   cache.MeshCache
	Commands *command.IndirectDrawBuffer

	id material.ID
}

func (m *DeferredMaterial) ID() material.ID {
	return m.id
}

func (m *DeferredMaterial) Begin(camera uniform.Camera, lights []light.T) {
	m.Instance.Descriptors().Camera.Set(camera)
	m.Objects.Reset()
	m.Commands.Reset()
}

func (m *DeferredMaterial) Bind(cmds command.Recorder) {
	cmds.Record(func(cmd command.Buffer) {
		m.Instance.Bind(cmd)
		m.Commands.BeginDrawIndirect()
	})
}

func (m *DeferredMaterial) Draw(cmds command.Recorder, msh mesh.Mesh) {
	gpuMesh, meshReady := m.Meshes.TryFetch(msh.Mesh().Get())
	if !meshReady {
		return
	}

	textures := m.Instance.Material().TextureSlots()
	textureIds := AssignMeshTextures(m.Textures, msh, textures)

	instanceId := m.Objects.Store(uniform.Object{
		Model:    msh.Transform().Matrix(),
		Textures: textureIds,
	})

	cmds.Record(func(cmd command.Buffer) {
		gpuMesh.Bind(cmd)
		gpuMesh.Draw(m.Commands, instanceId)
	})
}

func (m *DeferredMaterial) Unbind(cmds command.Recorder) {
	cmds.Record(func(cmd command.Buffer) {
		m.Commands.EndDrawIndirect(cmd)
	})
}

func (m *DeferredMaterial) End() {
	m.Objects.Flush(m.Instance.Descriptors().Objects)
	m.Textures.Flush()
	m.Commands.Flush()
}

func (m *DeferredMaterial) Destroy() {
	m.Instance.Material().Destroy()
	m.Commands.Destroy()
}
