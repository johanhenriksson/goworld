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

type BasicDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[uniform.Camera]
	Objects *descriptor.Storage[uniform.Object]
}

// Basic Materials only contain camera & object descriptors
// They can be used for various untextured objects, such
// as shadow/depth passes and lines.
type BasicMaterial struct {
	Material    *material.Material
	Descriptors *BasicDescriptors
	Objects     *ObjectBuffer
	Meshes      cache.MeshCache
	Commands    *command.IndirectDrawBuffer

	id material.ID
}

func (m *BasicMaterial) ID() material.ID {
	return m.id
}

func (m *BasicMaterial) Begin(camera uniform.Camera, lights []light.T) {
	m.Descriptors.Camera.Set(camera)
	m.Objects.Reset()
	m.Commands.Reset()
}

func (m *BasicMaterial) Bind(cmds command.Recorder) {
	cmds.Record(func(cmd *command.Buffer) {
		m.Material.Bind(cmd)
		cmd.CmdBindGraphicsDescriptor(m.Descriptors)
		m.Commands.BeginDrawIndirect()
	})
}

func (m *BasicMaterial) Draw(cmds command.Recorder, msh mesh.Mesh) {
	gpuMesh, meshReady := m.Meshes.TryFetch(msh.Mesh())
	if !meshReady {
		return
	}

	instanceId := m.Objects.Store(uniform.Object{
		Model: msh.Transform().Matrix(),
	})

	cmds.Record(func(cmd *command.Buffer) {
		gpuMesh.Bind(cmd)
		gpuMesh.Draw(m.Commands, instanceId)
	})
}

func (m *BasicMaterial) Unbind(cmds command.Recorder) {
	cmds.Record(func(cmd *command.Buffer) {
		m.Commands.EndDrawIndirect(cmd)
	})
}

func (m *BasicMaterial) End() {
	m.Objects.Flush(m.Descriptors.Objects)
	m.Commands.Flush()
}

func (m *BasicMaterial) Destroy() {
	m.Material.Destroy()
	m.Commands.Destroy()
}
