package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
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

type ForwardDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Lights   *descriptor.Storage[uniform.Light]
	Textures *descriptor.SamplerArray
}

// Material Instance Wrapper - wraps material descriptors in helpers
type ForwardMatData struct {
	Instance *material.Instance[*ForwardDescriptors]
	Objects  *ObjectBuffer
	Lights   *LightBuffer
	Shadows  *ShadowCache
	Textures cache.SamplerCache
}

type ForwardMaterialMaker struct {
	app    vulkan.App
	pass   renderpass.T
	lookup ShadowmapLookupFn
}

var _ MaterialMaker[*ForwardMatData] = &ForwardMaterialMaker{}

func NewForwardMaterialMaker(app vulkan.App, pass renderpass.T, lookup ShadowmapLookupFn) *ForwardMaterialMaker {
	return &ForwardMaterialMaker{
		app:    app,
		pass:   pass,
		lookup: lookup,
	}
}

func (m *ForwardMaterialMaker) Instantiate(def *material.Def, count int) []*ForwardMatData {
	if def == nil {
		def = material.StandardForward()
	}

	desc := &ForwardDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   2000,
		},
		Lights: &descriptor.Storage[uniform.Light]{
			Stages: core1_0.StageAll,
			Size:   256,
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

	instances := make([]*ForwardMatData, count)
	for i := range instances {
		instance := mat.Instantiate(m.app.Pool())
		textures := cache.NewSamplerCache(m.app.Textures(), instance.Descriptors().Textures)
		instances[i] = &ForwardMatData{
			Instance: instance,
			Objects:  NewObjectBuffer(desc.Objects.Size),
			Lights:   NewLightBuffer(desc.Lights.Size),
			Shadows:  NewShadowCache(textures, m.lookup),
			Textures: textures,
		}
	}

	return instances
}

func (m *ForwardMaterialMaker) Destroy(mat *ForwardMatData) {
	mat.Instance.Material().Destroy()
}

func (m *ForwardMaterialMaker) BeginFrame(mat *ForwardMatData, camera uniform.Camera, lights []light.T) {
	mat.Instance.Descriptors().Camera.Set(camera)

	// multiple calls to this reset in a single frame will cause weird behaviour
	// we need to split this function somehow in order to be able to do depth sorting etc
	mat.Objects.Reset()

	if len(lights) > 0 {
		// how to get ambient light info?
		mat.Lights.Reset()
		for _, lit := range lights {
			mat.Lights.Store(lit.LightData(mat.Shadows))
		}
		mat.Lights.Flush(mat.Instance.Descriptors().Lights)
	}
}

func (m *ForwardMaterialMaker) PrepareMesh(mat *ForwardMatData, msh mesh.Mesh) int {
	textures := mat.Instance.Material().TextureSlots()
	textureIds := AssignMeshTextures(mat.Textures, msh, textures)

	return mat.Objects.Store(uniform.Object{
		Model:    msh.Transform().Matrix(),
		Textures: textureIds,
	})
}

func (m *ForwardMaterialMaker) EndFrame(mat *ForwardMatData) {
	mat.Objects.Flush(mat.Instance.Descriptors().Objects)
	mat.Textures.Flush()
}

func (m *ForwardMaterialMaker) Draw(cmds command.Recorder, camera uniform.Camera, group *MeshGroup[*ForwardMatData], lights []light.T) {
	mat := group.Material

	m.BeginFrame(mat, camera, lights)

	cmds.Record(func(cmd command.Buffer) {
		mat.Instance.Bind(cmd)
	})
	for _, msh := range group.Meshes {
		vkmesh, meshReady := m.app.Meshes().TryFetch(msh.Mesh().Get())
		if !meshReady {
			continue
		}

		index := m.PrepareMesh(mat, msh)

		cmds.Record(func(cmd command.Buffer) {
			vkmesh.Draw(cmd, index)
		})
	}

	m.EndFrame(mat)
}
