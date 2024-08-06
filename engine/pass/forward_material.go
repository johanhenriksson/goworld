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

type ForwardDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Lights   *descriptor.Storage[uniform.Light]
	Textures *descriptor.SamplerArray
}

type ForwardMaterial struct {
	Instance *material.Instance[*ForwardDescriptors]
	Objects  *ObjectBuffer
	Lights   *LightBuffer
	Shadows  *ShadowCache
	Textures cache.SamplerCache
	Meshes   cache.MeshCache

	id material.ID
}

func (m *ForwardMaterial) ID() material.ID {
	return m.id
}

func (m *ForwardMaterial) Begin(camera uniform.Camera, lights []light.T) {
	m.Instance.Descriptors().Camera.Set(camera)

	// multiple calls to this reset in a single frame will cause weird behaviour
	// we need to split this function somehow in order to be able to do depth sorting etc
	m.Objects.Reset()

	if len(lights) > 0 {
		// how to get ambient light info?
		m.Lights.Reset()
		for _, lit := range lights {
			m.Lights.Store(lit.LightData(m.Shadows))
		}
		m.Lights.Flush(m.Instance.Descriptors().Lights)
	}
}

func (m *ForwardMaterial) Bind(cmds command.Recorder) {
	cmds.Record(func(cmd command.Buffer) {
		m.Instance.Bind(cmd)
	})
}

func (m *ForwardMaterial) Draw(cmds command.Recorder, msh mesh.Mesh) {
	vkmesh, meshReady := m.Meshes.TryFetch(msh.Mesh().Get())
	if !meshReady {
		return
	}

	textures := m.Instance.Material().TextureSlots()
	textureIds := AssignMeshTextures(m.Textures, msh, textures)

	index := m.Objects.Store(uniform.Object{
		Model:    msh.Transform().Matrix(),
		Textures: textureIds,
	})

	cmds.Record(func(cmd command.Buffer) {
		vkmesh.Bind(cmd)
		vkmesh.Draw(cmd, index)
	})
}

func (m *ForwardMaterial) Unbind(cmds command.Recorder) {
}

func (m *ForwardMaterial) End() {
	m.Objects.Flush(m.Instance.Descriptors().Objects)
	m.Textures.Flush()
}

func (m *ForwardMaterial) Destroy() {
	m.Instance.Material().Destroy()
}
