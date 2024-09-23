package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
)

type ForwardDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Lights   *descriptor.Storage[uniform.Light]
	Textures *descriptor.SamplerArray
}

type ForwardMaterial struct {
	Pipeline    *pipeline.Pipeline
	Descriptors *ForwardDescriptors
	Objects     *ObjectBuffer
	Lights      *LightBuffer
	Shadows     *ShadowCache
	Textures    cache.SamplerCache
	Meshes      cache.MeshCache
	Commands    *command.IndirectDrawBuffer

	id material.ID
}

func (m *ForwardMaterial) ID() material.ID {
	return m.id
}

func (m *ForwardMaterial) Begin(camera uniform.Camera, lights []light.T) {
	m.Descriptors.Camera.Set(camera)

	m.Commands.Reset()

	if len(lights) > 0 {
		// how to get ambient light info?
		m.Lights.Reset()
		for _, lit := range lights {
			m.Lights.Store(lit.LightData(m.Shadows))
		}
		m.Lights.Flush(m.Descriptors.Lights)
	}
}

func (m *ForwardMaterial) Bind(cmds command.Recorder) {
	cmds.Record(func(cmd *command.Buffer) {

		// we actually want to do this once
		cmd.CmdBindGraphicsDescriptor(m.Pipeline.Layout(), 0, m.Descriptors)

		cmd.CmdBindGraphicsPipeline(m.Pipeline)
		m.Commands.BeginDrawIndirect()
	})
}

func (m *ForwardMaterial) Draw(cmds command.Recorder, msh mesh.Mesh) {
	gpuMesh, meshReady := m.Meshes.TryFetch(msh.Mesh())
	if !meshReady {
		return
	}

	textures := m.Pipeline.Shader().Textures()
	textureIds := AssignMeshTextures(m.Textures, msh, textures)

	instanceId := m.Objects.Store(uniform.Object{
		Model:    msh.Transform().Matrix(),
		Textures: textureIds,
	})

	cmds.Record(func(cmd *command.Buffer) {
		gpuMesh.Bind(cmd)
		gpuMesh.Draw(m.Commands, instanceId)
	})
}

func (m *ForwardMaterial) Unbind(cmds command.Recorder) {
	cmds.Record(func(cmd *command.Buffer) {
		m.Commands.EndDrawIndirect(cmd)
	})
}

func (m *ForwardMaterial) End() {
	m.Objects.Flush(m.Descriptors.Objects)
	m.Textures.Flush(m.Descriptors.Textures)
	m.Commands.Flush()
}

func (m *ForwardMaterial) Destroy() {
	m.Descriptors.Destroy()
	m.Pipeline.Destroy()
	m.Commands.Destroy()
}
