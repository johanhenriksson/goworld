package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
)

type DepthDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[uniform.Camera]
	Objects *descriptor.Storage[uniform.Object]
}

type DepthMaterial struct {
	Instance *material.Instance[*DepthDescriptors]
	Objects  *ObjectBuffer
	Meshes   cache.MeshCache

	id material.ID
}

func (m *DepthMaterial) ID() material.ID {
	return m.id
}

func (m *DepthMaterial) Begin(camera uniform.Camera, lights []light.T) {
	m.Instance.Descriptors().Camera.Set(camera)
	m.Objects.Reset()
}

func (m *DepthMaterial) Bind(cmds command.Recorder) {
	cmds.Record(func(cmd command.Buffer) {
		m.Instance.Bind(cmd)
	})
}

func (m *DepthMaterial) End() {
	m.Objects.Flush(m.Instance.Descriptors().Objects)
}

func (m *DepthMaterial) Draw(cmds command.Recorder, msh mesh.Mesh) {
	vkmesh, meshReady := m.Meshes.TryFetch(msh.Mesh().Get())
	if !meshReady {
		return
	}

	index := m.Objects.Store(uniform.Object{
		Model: msh.Transform().Matrix(),
	})

	cmds.Record(func(cmd command.Buffer) {
		vkmesh.Draw(cmd, index)
	})
}

func (m *DepthMaterial) Destroy() {
	m.Instance.Material().Destroy()
}
