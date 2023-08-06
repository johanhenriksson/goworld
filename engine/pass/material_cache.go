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

type MaterialCache[T any] struct {
	TransformFn MaterialTransform

	cache      map[uint64][]T
	defaultMat *material.Def
	app        vulkan.App
	frames     int
	pass       renderpass.T
	maker      MaterialMaker[T]
}

func NewMaterialCache[T any](app vulkan.App, frames int, pass renderpass.T, maker MaterialMaker[T], defaultMat *material.Def, transform MaterialTransform) *MaterialCache[T] {
	ms := &MaterialCache[T]{
		app:        app,
		frames:     frames,
		pass:       pass,
		defaultMat: defaultMat,
		cache:      map[uint64][]T{},
		maker:      maker,

		TransformFn: transform,
	}
	ms.Load(defaultMat)
	return ms
}

func (m *MaterialCache[T]) Get(matId uint64, frame int) T {
	return m.cache[matId][frame]
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

func (m *MaterialCache[T]) Load(def *material.Def) bool {
	id := material.Hash(def)
	if def == nil {
		def = m.defaultMat
	}

	// apply material transform
	if m.TransformFn != nil {
		def = m.TransformFn(def)
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
		return false
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

	m.cache[id] = make([]T, m.frames)
	for frame := 0; frame < m.frames; frame++ {
		m.cache[id][frame] = m.maker.Instantiate(mat)
	}

	// indicate that the material is ready to be used
	return true
}

type MaterialMaker[T any] interface {
	Instantiate(mat material.T[*material.Descriptors]) T
	Destroy(T)
	Draw(cmds command.Recorder, mat T, camera uniform.Camera, objects []mesh.Mesh, lights []light.T)
}

type StdMaterialMaker struct {
	app    vulkan.App
	lookup ShadowmapLookupFn
}

func (m *StdMaterialMaker) Instantiate(mat material.T[*material.Descriptors]) *MaterialData {
	instance := mat.Instantiate(m.app.Pool())
	textures := cache.NewSamplerCache(m.app.Textures(), instance.Descriptors().Textures)
	return &MaterialData{
		Instance: instance,
		Objects:  make([]uniform.Object, instance.Descriptors().Objects.Size),
		Textures: textures,
		Lights:   NewLightBuffer(),
		Shadows:  NewShadowCache(textures, m.lookup),
	}
}

func (m *StdMaterialMaker) Destroy(mat *MaterialData) {
	mat.Instance.Material().Destroy()
}

func (m *StdMaterialMaker) Draw(cmds command.Recorder, mat *MaterialData, camera uniform.Camera, objects []mesh.Mesh, lights []light.T) {
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

	for i, msh := range objects {
		vkmesh, meshReady := m.app.Meshes().TryFetch(msh.Mesh().Get())
		if !meshReady {
			continue
		}

		textureIds := [4]uint32{}
		textures := mat.Instance.Material().TextureSlots()
		for id, textureSlot := range textures {
			ref := msh.Texture(textureSlot)
			if ref != nil {
				handle, exists := mat.Textures.TryFetch(ref)
				if exists {
					textureIds[id] = uint32(handle.ID)
				}
			}
		}

		mat.Objects[i] = uniform.Object{
			Model:    msh.Transform().Matrix(),
			Textures: textureIds,
		}

		index := i
		cmds.Record(func(cmd command.Buffer) {
			vkmesh.Draw(cmd, index)
		})
	}

	mat.Instance.Descriptors().Objects.SetRange(0, mat.Objects[:len(objects)])
	mat.Textures.Flush()
}
