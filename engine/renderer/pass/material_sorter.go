package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
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
	pass       renderpass.T
}

type MaterialData struct {
	Instance material.Instance[*material.Descriptors]
	Objects  []uniform.Object
	Textures cache.SamplerCache
}

func NewMaterialSorter(app vulkan.App, pass renderpass.T, defaultMat *material.Def) *MaterialSorter {
	ms := &MaterialSorter{
		app:        app,
		pass:       pass,
		defaultMat: defaultMat,
		cache:      map[uint64][]*MaterialData{},

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
			Subpass:    def.Subpass,
			Pointers:   pointers,
			DepthTest:  def.DepthTest,
			DepthWrite: def.DepthWrite,
			DepthClamp: def.DepthClamp,
			Primitive:  def.Primitive,
			CullMode:   def.CullMode,
		},
		desc)

	m.cache[id] = make([]*MaterialData, m.app.Frames())
	for frame := 0; frame < m.app.Frames(); frame++ {
		instance := mat.Instantiate(m.app.Pool())
		m.cache[id][frame] = &MaterialData{
			Instance: instance,
			Objects:  make([]uniform.Object, instance.Descriptors().Objects.Size),
			Textures: cache.NewSamplerCache(m.app.Textures(), instance.Descriptors().Textures),
		}
	}

	// indicate that the material is ready to be used
	return true
}

func (m *MaterialSorter) Draw(cmds command.Recorder, args render.Args, meshes []mesh.Component) {
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
	m.DrawCamera(cmds, args, camera, meshes)
}

func (m *MaterialSorter) DrawCamera(cmds command.Recorder, args render.Args, camera uniform.Camera, meshes []mesh.Component) {
	// sort meshes by material
	meshGroups := map[uint64][]mesh.Component{}
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

	for matId, matMeshes := range meshGroups {
		mat := m.cache[matId][args.Context.Index]
		mat.Instance.Descriptors().Camera.Set(camera)

		cmds.Record(func(cmd command.Buffer) {
			mat.Instance.Bind(cmd)
		})

		for i, msh := range matMeshes {
			vkmesh, meshReady := m.app.Meshes().TryFetch(msh.Mesh())
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

		mat.Instance.Descriptors().Objects.SetRange(0, mat.Objects[:len(matMeshes)])
		mat.Textures.UpdateDescriptors()
	}
}
