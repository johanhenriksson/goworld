package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/render"
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

	cache      map[uint64][]material.Instance[*material.Descriptors]
	defaultMat *material.Def
	target     vulkan.Target
	pass       renderpass.T
}

func NewMaterialSorter(target vulkan.Target, pass renderpass.T, defaultMat *material.Def) *MaterialSorter {
	ms := &MaterialSorter{
		target:     target,
		pass:       pass,
		defaultMat: defaultMat,
		cache:      map[uint64][]material.Instance[*material.Descriptors]{},

		TransformFn: func(d *material.Def) *material.Def { return d },
	}
	ms.Load(defaultMat)
	return ms
}

func (m *MaterialSorter) Destroy() {
	for _, mat := range m.cache {
		mat[0].Material().Destroy()
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
	shader, shaderReady := m.target.Shaders().Fetch(shader.NewRef(def.Shader))
	if !shaderReady {
		// pending
		return false
	}

	// create material
	mat := material.New(
		m.target.Device(),
		material.Args{
			Shader:     shader,
			Pass:       m.pass,
			Subpass:    def.Subpass,
			Pointers:   pointers,
			DepthTest:  def.DepthTest,
			DepthWrite: def.DepthWrite,
			Primitive:  def.Primitive,
			CullMode:   def.CullMode,
		},
		desc)

	m.cache[id] = mat.InstantiateMany(m.target.Pool(), m.target.Frames())

	// indicate that the material is ready to be used
	return true
}

func (m *MaterialSorter) Draw(cmds command.Recorder, args render.Args, meshes []mesh.T) {
	camera := uniform.Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         args.Position,
	}
	m.DrawCamera(cmds, args, camera, meshes)
}

func (m *MaterialSorter) DrawCamera(cmds command.Recorder, args render.Args, camera uniform.Camera, meshes []mesh.T) {
	// sort meshGroups by material
	meshGroups := map[uint64][]mesh.T{}
	for _, msh := range meshes {
		mid := msh.MaterialID()
		if _, exists := m.cache[mid]; !exists {
			// initialize material
			if !m.Load(msh.Material()) {
				continue
			}
		}
		meshGroups[mid] = append(meshGroups[mid], msh)
	}

	descriptors := make([]uniform.Object, len(meshes))

	index := 0
	for mid, meshes := range meshGroups {
		mat := m.cache[mid][args.Context.Index]
		mat.Descriptors().Camera.Set(camera)

		cmds.Record(func(cmd command.Buffer) {
			mat.Bind(cmd)
		})

		begin := index
		for _, msh := range meshes {
			vkmesh, meshReady := m.target.Meshes().Fetch(msh.Mesh())
			if !meshReady {
				continue
			}

			descriptors[index] = uniform.Object{
				Model: msh.Transform().World(),
			}

			i := index - begin
			cmds.Record(func(cmd command.Buffer) {
				vkmesh.Draw(cmd, i)
			})

			index++
		}

		mat.Descriptors().Objects.SetRange(0, descriptors[begin:index-begin])
	}
}
