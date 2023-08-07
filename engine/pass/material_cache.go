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
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type MaterialCache[T any] struct {
	cache      map[uint64][]T
	defaultMat *material.Def
	app        vulkan.App
	frames     int
	maker      MaterialMaker[T]
}

func NewMaterialCache[T any](app vulkan.App, frames int, pass renderpass.T, maker MaterialMaker[T], defaultMat *material.Def) *MaterialCache[T] {
	ms := &MaterialCache[T]{
		app:        app,
		frames:     frames,
		defaultMat: defaultMat,
		cache:      map[uint64][]T{},
		maker:      maker,
	}
	ms.Load(defaultMat)
	return ms
}

func (m *MaterialCache[T]) Get(msh mesh.Mesh, frame int) (T, bool) {
	matId := msh.MaterialID()
	mat, exists := m.cache[matId]
	if !exists {
		// initialize material
		var ready bool
		mat, ready = m.Load(msh.Material())
		if !ready {
			// not ready yet
			var empty T
			return empty, false
		}
	}
	return mat[frame], true
}

func (m *MaterialCache[T]) Exists(matId uint64) bool {
	_, exists := m.cache[matId]
	return exists
}

func (m *MaterialCache[T]) Destroy() {
	for _, mat := range m.cache {
		m.maker.Destroy(mat[0])
	}
	m.cache = nil
}

func (m *MaterialCache[T]) Load(def *material.Def) ([]T, bool) {
	id := material.Hash(def)
	if def == nil {
		def = m.defaultMat
	}

	mat := m.maker.Instantiate(def, m.frames)
	if len(mat) == 0 {
		// not ready yet
		return nil, false
	}

	m.cache[id] = mat
	return mat, true
}

type MaterialMaker[T any] interface {
	Instantiate(mat *material.Def, count int) []T
	Destroy(T)
	Draw(cmds command.Recorder, camera uniform.Camera, group *MeshGroup[T], lights []light.T)
}

type StdMaterialMaker struct {
	app       vulkan.App
	pass      renderpass.T
	lookup    ShadowmapLookupFn
	transform MaterialTransform
}

func (m *StdMaterialMaker) Instantiate(def *material.Def, count int) []*MaterialData {
	// apply material transform
	if m.transform != nil {
		def = m.transform(def)
	}

	desc := &material.Descriptors{
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

	instances := make([]*MaterialData, count)
	for i := range instances {
		instance := mat.Instantiate(m.app.Pool())
		textures := cache.NewSamplerCache(m.app.Textures(), instance.Descriptors().Textures)
		instances[i] = &MaterialData{
			Instance: instance,
			Objects:  make([]uniform.Object, 0, instance.Descriptors().Objects.Size),
			Textures: textures,
			Lights:   NewLightBuffer(),
			Shadows:  NewShadowCache(textures, m.lookup),
		}
	}

	return instances
}

func (m *StdMaterialMaker) Destroy(mat *MaterialData) {
	mat.Instance.Material().Destroy()
}

func (m *StdMaterialMaker) Draw(cmds command.Recorder, camera uniform.Camera, group *MeshGroup[*MaterialData], lights []light.T) {
	mat := group.Material
	mat.Instance.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		mat.Instance.Bind(cmd)
	})

	if len(lights) > 0 {
		// how to get ambient light info?
		mat.Lights.Reset()
		for _, lit := range lights {
			mat.Lights.Store(lit.LightData(mat.Shadows))
		}
		mat.Lights.Flush(mat.Instance.Descriptors().Lights)
	}

	mat.Objects = mat.Objects[:0]
	for i, msh := range group.Meshes {
		vkmesh, meshReady := m.app.Meshes().TryFetch(msh.Mesh().Get())
		if !meshReady {
			continue
		}

		textures := mat.Instance.Material().TextureSlots()
		textureIds := AssignMaterialTextures(mat.Textures, msh, textures)

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

func AssignMaterialTextures(samplers cache.SamplerCache, msh mesh.Mesh, slots []texture.Slot) [4]uint32 {
	textureIds := [4]uint32{}
	for id, slot := range slots {
		ref := msh.Texture(slot)
		if ref != nil {
			handle, exists := samplers.TryFetch(ref)
			if exists {
				textureIds[id] = uint32(handle.ID)
			}
		}
	}
	return textureIds
}
