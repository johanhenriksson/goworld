package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type DeferredDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type DeferredMatData struct {
	Instance material.Instance[*DeferredDescriptors]
	Objects  []uniform.Object
	Textures cache.SamplerCache
}

type DeferredMaterialMaker struct {
	app  vulkan.App
	pass renderpass.T
}

func NewDeferredMaterialMaker(app vulkan.App, pass renderpass.T) MaterialMaker[*DeferredMatData] {
	return &DeferredMaterialMaker{
		app:  app,
		pass: pass,
	}
}

func (m *DeferredMaterialMaker) Instantiate(def *material.Def, count int) []*DeferredMatData {
	if def == nil {
		def = material.StandardDeferred()
	}

	desc := &DeferredDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   2000,
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  100,
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

	instances := make([]*DeferredMatData, count)
	for i := range instances {
		instance := mat.Instantiate(m.app.Pool())
		textures := cache.NewSamplerCache(m.app.Textures(), instance.Descriptors().Textures)
		instances[i] = &DeferredMatData{
			Instance: instance,
			Objects:  make([]uniform.Object, 0, instance.Descriptors().Objects.Size),
			Textures: textures,
		}
	}

	return instances
}

func (m *DeferredMaterialMaker) Destroy(mat *DeferredMatData) {
	mat.Instance.Material().Destroy()
}

func (m *DeferredMaterialMaker) Draw(cmds command.Recorder, camera uniform.Camera, group *MeshGroup[*DeferredMatData], lights []light.T) {
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

		textures := mat.Instance.Material().TextureSlots()
		textureIds := FetchMaterialTextures(mat.Textures, msh, textures)

		mat.Objects = append(mat.Objects, uniform.Object{
			Model:    msh.Transform().Matrix(),
			Textures: textureIds,
		})

		index := i
		cmds.Record(func(cmd command.Buffer) {
			vkmesh.Draw(cmd, index)
		})
	}

	mat.Instance.Descriptors().Objects.SetRange(0, mat.Objects)
	mat.Textures.Flush()
}
