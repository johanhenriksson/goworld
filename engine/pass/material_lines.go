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

type LineMatData struct {
	Instance material.Instance[*DepthDescriptors]
	Objects  []uniform.Object
}

type LineMaterialMaker struct {
	app  vulkan.App
	pass renderpass.T
}

func NewLineMaterialMaker(app vulkan.App, pass renderpass.T) MaterialMaker[*LineMatData] {
	return &LineMaterialMaker{
		app:  app,
		pass: pass,
	}
}

func (m *LineMaterialMaker) Instantiate(def *material.Def, count int) []*LineMatData {
	if def == nil {
		def = material.Lines()
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
	shader, shaderReady := m.app.Shaders().TryFetch(shader.NewRef(def.Shader))
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
			DepthTest:  def.DepthTest,
			DepthWrite: def.DepthWrite,
			DepthClamp: def.DepthClamp,
			DepthFunc:  def.DepthFunc,
			Primitive:  def.Primitive,
			CullMode:   def.CullMode,
		},
		desc)

	instances := make([]*LineMatData, count)
	for i := range instances {
		instance := mat.Instantiate(m.app.Pool())
		instances[i] = &LineMatData{
			Instance: instance,
			Objects:  make([]uniform.Object, 0, instance.Descriptors().Objects.Size),
		}
	}

	return instances
}

func (m *LineMaterialMaker) Destroy(mat *LineMatData) {
	mat.Instance.Material().Destroy()
}

func (m *LineMaterialMaker) Draw(cmds command.Recorder, camera uniform.Camera, group *MeshGroup[*LineMatData], lights []light.T) {
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
