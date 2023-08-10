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

type BasicDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[uniform.Camera]
	Objects *descriptor.Storage[uniform.Object]
}

// Basic Materials only contain camera & object descriptors
// They can be used for various untextured objects, such
// as shadow/depth passes and lines.
type BasicMaterial struct {
	Instance *material.Instance[*BasicDescriptors]
	Objects  *ObjectBuffer
	Meshes   cache.MeshCache

	id material.ID
}

func (m *BasicMaterial) ID() material.ID {
	return m.id
}

func (m *BasicMaterial) Begin(camera uniform.Camera, lights []light.T) {
	m.Instance.Descriptors().Camera.Set(camera)
	m.Objects.Reset()
}

func (m *BasicMaterial) Bind(cmds command.Recorder) {
	cmds.Record(func(cmd command.Buffer) {
		m.Instance.Bind(cmd)
	})
}

func (m *BasicMaterial) End() {
	m.Objects.Flush(m.Instance.Descriptors().Objects)
}

func (m *BasicMaterial) Draw(cmds command.Recorder, msh mesh.Mesh) {
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

func (m *BasicMaterial) Destroy() {
	m.Instance.Material().Destroy()
}
