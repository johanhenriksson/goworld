package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type MaterialTransform func(*material.Def) *material.Def

type MaterialSorter struct {
	TransformFn MaterialTransform

	cache      map[uint64]material.Standard
	defaultMat *material.Def
	target     vulkan.Target
	pass       renderpass.T
}

func NewMaterialSorter(target vulkan.Target, pass renderpass.T, defaultMat *material.Def) *MaterialSorter {
	ms := &MaterialSorter{
		target:     target,
		pass:       pass,
		defaultMat: defaultMat,
		cache:      map[uint64]material.Standard{},

		TransformFn: func(d *material.Def) *material.Def { return d },
	}
	ms.Load(defaultMat)
	return ms
}

func (m *MaterialSorter) Destroy() {
	for _, mat := range m.cache {
		mat.Material().Destroy()
	}
}

func (m *MaterialSorter) Load(def *material.Def) {
	id := material.Hash(def)
	if def == nil {
		def = m.defaultMat
	}
	log.Println("instantiate new forward material id", id, def)
	mat := material.FromDef(m.target.Device(), m.target.Pool(), m.pass, def)
	m.cache[id] = mat
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
	m.DrawCamera(cmds, camera, meshes)
}

func (m *MaterialSorter) DrawCamera(cmds command.Recorder, camera uniform.Camera, meshes []mesh.T) {
	// sort meshGroups by material
	meshGroups := map[uint64][]mesh.T{}
	for _, msh := range meshes {
		mid := msh.MaterialID()
		if _, exists := m.cache[mid]; !exists {
			// initialize material
			m.Load(msh.Material())
		}
		meshGroups[mid] = append(meshGroups[mid], msh)
	}

	index := 0
	for mid, meshes := range meshGroups {
		mat := m.cache[mid]
		mat.Descriptors().Camera.Set(camera)

		cmds.Record(func(cmd command.Buffer) {
			mat.Bind(cmd)
		})

		for _, msh := range meshes {
			vkmesh := m.target.Meshes().Fetch(msh.Mesh())
			if vkmesh == nil {
				return
			}

			mat.Descriptors().Objects.Set(index, uniform.Object{
				Model: msh.Transform().World(),
			})

			i := index
			cmds.Record(func(cmd command.Buffer) {
				vkmesh.Draw(cmd, i)
			})

			index++
		}
	}
}
