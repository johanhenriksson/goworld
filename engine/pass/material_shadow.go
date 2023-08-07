package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ShadowMatData struct {
	Instance material.Instance[*DepthDescriptors]
	Objects  *ObjectBuffer
}

type ShadowMaterialMaker struct {
	app  vulkan.App
	pass renderpass.T
}

func NewShadowMaterialMaker(app vulkan.App, pass renderpass.T) MaterialMaker[*ShadowMatData] {
	return &ShadowMaterialMaker{
		app:  app,
		pass: pass,
	}
}

func (m *ShadowMaterialMaker) Instantiate(def *material.Def, count int) []*ShadowMatData {
	if def == nil {
		def = &material.Def{}
	}

	desc := &DepthDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   2000,
		},
	}

	// read vertex pointers from vertex format
	pointers := vertex.ParsePointers(def.VertexFormat)

	// fetch shader from cache
	shader, shaderReady := m.app.Shaders().TryFetch(shader.NewRef("shadow"))
	if !shaderReady {
		// pending
		return nil
	}

	// create material
	mat := material.New(
		m.app.Device(),
		material.Args{
			Shader:     shader,
			Pass:       m.pass,
			Subpass:    MainSubpass,
			Pointers:   pointers,
			CullMode:   vertex.CullFront,
			DepthTest:  true,
			DepthWrite: true,
			DepthFunc:  core1_0.CompareOpLess,
			DepthClamp: true,
			Primitive:  def.Primitive,
		},
		desc)

	instances := make([]*ShadowMatData, count)
	for i := range instances {
		instance := mat.Instantiate(m.app.Pool())
		instances[i] = &ShadowMatData{
			Instance: instance,
			Objects:  NewObjectBuffer(desc.Objects.Size),
		}
	}

	return instances
}

func (m *ShadowMaterialMaker) Destroy(mat *ShadowMatData) {
	mat.Instance.Material().Destroy()
}

func (m *ShadowMaterialMaker) Draw(cmds command.Recorder, camera uniform.Camera, group *MeshGroup[*ShadowMatData], lights []light.T) {
	mat := group.Material
	mat.Instance.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		mat.Instance.Bind(cmd)
	})

	mat.Objects.Reset()
	for i, msh := range group.Meshes {
		vkmesh, meshReady := m.app.Meshes().TryFetch(msh.Mesh().Get())
		if !meshReady {
			continue
		}

		mat.Objects.Store(uniform.Object{
			Model: msh.Transform().Matrix(),
		})

		index := i
		cmds.Record(func(cmd command.Buffer) {
			vkmesh.Draw(cmd, index)
		})
	}

	mat.Objects.Flush(mat.Instance.Descriptors().Objects)
}
