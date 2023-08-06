package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render"
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

type MaterialTransform func(*material.Def) *material.Def

type MaterialSorter struct {
	TransformFn MaterialTransform

	cache      map[uint64][]*MaterialData
	defaultMat *material.Def
	app        vulkan.App
	frames     int
	pass       renderpass.T
	lookup     ShadowmapLookupFn
}

type MaterialData struct {
	Instance material.Instance[*material.Descriptors]
	Objects  []uniform.Object
	Lights   *LightBuffer
	Textures cache.SamplerCache
}

func NewMaterialSorter(app vulkan.App, frames int, pass renderpass.T, lookup ShadowmapLookupFn, defaultMat *material.Def) *MaterialSorter {
	ms := &MaterialSorter{
		app:        app,
		frames:     frames,
		pass:       pass,
		defaultMat: defaultMat,
		cache:      map[uint64][]*MaterialData{},
		lookup:     lookup,

		TransformFn: func(d *material.Def) *material.Def { return d },
	}
	ms.Load(defaultMat)
	return ms
}

func (m *MaterialSorter) Destroy() {
	for _, mat := range m.cache {
		mat[0].Instance.Material().Destroy()
		// todo: destroy more
	}
}

func (m *MaterialSorter) Load(def *material.Def) bool {
	id := material.Hash(def)
	if def == nil {
		def = m.defaultMat
	}

	// apply material transform
	def = m.TransformFn(def)

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

	m.cache[id] = make([]*MaterialData, m.frames)
	for frame := 0; frame < m.frames; frame++ {
		instance := mat.Instantiate(m.app.Pool())
		textures := cache.NewSamplerCache(m.app.Textures(), instance.Descriptors().Textures)
		m.cache[id][frame] = &MaterialData{
			Instance: instance,
			Objects:  make([]uniform.Object, instance.Descriptors().Objects.Size),
			Textures: textures,
			Lights:   NewLightBuffer(instance.Descriptors().Lights, textures, m.lookup),
		}
	}

	// indicate that the material is ready to be used
	return true
}

func (m *MaterialSorter) Draw(cmds command.Recorder, args render.Args, meshes []mesh.Mesh, lights []light.T) {
	camera := uniform.Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         args.Position,
		Forward:     args.Forward,
	}
	m.DrawCamera(cmds, args, camera, meshes, lights)
}

func (m *MaterialSorter) DrawCamera(cmds command.Recorder, args render.Args, camera uniform.Camera, meshes []mesh.Mesh, lights []light.T) {
	// sort meshes by material
	meshGroups := map[uint64][]mesh.Mesh{}
	for _, msh := range meshes {
		matId := msh.MaterialID()
		if _, exists := m.cache[matId]; !exists {
			// initialize material
			if !m.Load(msh.Material()) {
				continue
			}
		}
		meshGroups[matId] = append(meshGroups[matId], msh)
	}

	// iterate the sorted material groups
	for matId, objects := range meshGroups {
		mat := m.cache[matId][args.Context.Index]
		mat.Instance.Descriptors().Camera.Set(camera)

		cmds.Record(func(cmd command.Buffer) {
			mat.Instance.Bind(cmd)
		})

		if len(lights) > 0 {
			// how to get ambient light info?
			mat.Lights.Reset()
			for _, lit := range lights {
				mat.Lights.Store(args, lit)
			}
			mat.Lights.Flush()
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
		mat.Textures.UpdateDescriptors()
	}
}
