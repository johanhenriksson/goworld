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

type DepthDescriptors struct {
	descriptor.Set
	Camera  *descriptor.Uniform[uniform.Camera]
	Objects *descriptor.Storage[uniform.Object]
}

type DepthMatData struct {
	Instance material.Instance[*DepthDescriptors]
	Objects  []uniform.Object
}

type DepthMaterialMaker struct {
	app  vulkan.App
	pass renderpass.T
}

func NewDepthMaterialMaker(app vulkan.App, pass renderpass.T) MaterialMaker[*DepthMatData] {
	return &DepthMaterialMaker{
		app:  app,
		pass: pass,
	}
}

func (m *DepthMaterialMaker) Instantiate(def *material.Def, count int) []*DepthMatData {
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
	shader, shaderReady := m.app.Shaders().TryFetch(shader.NewRef("depth"))
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
			CullMode:   vertex.CullBack,
			DepthTest:  true,
			DepthWrite: true,
			DepthFunc:  core1_0.CompareOpLess,
			DepthClamp: def.DepthClamp,
			Primitive:  def.Primitive,
		},
		desc)

	instances := make([]*DepthMatData, count)
	for i := range instances {
		instance := mat.Instantiate(m.app.Pool())
		instances[i] = &DepthMatData{
			Instance: instance,
			Objects:  make([]uniform.Object, 0, instance.Descriptors().Objects.Size),
		}
	}

	return instances
}

func (m *DepthMaterialMaker) Destroy(mat *DepthMatData) {
	mat.Instance.Material().Destroy()
}

func (m *DepthMaterialMaker) Draw(cmds command.Recorder, camera uniform.Camera, group *MeshGroup[*DepthMatData], lights []light.T) {
	mat := group.Material
	mat.Instance.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		mat.Instance.Bind(cmd)
	})

	mat.Objects = mat.Objects[:0]
	for i, msh := range group.Meshes {
		vkmesh, meshReady := m.app.Meshes().TryFetch(msh.Mesh().Get())
		if !meshReady {
			continue
		}

		mat.Objects = append(mat.Objects, uniform.Object{
			Model: msh.Transform().Matrix(),
		})

		index := i
		cmds.Record(func(cmd command.Buffer) {
			vkmesh.Draw(cmd, index)
		})
	}

	mat.Instance.Descriptors().Objects.SetRange(0, mat.Objects)
}
