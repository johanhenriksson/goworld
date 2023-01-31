package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type MaterialSorter struct {
	cache      map[uint64]material.Standard
	defaultMat *material.Def
	target     vulkan.Target
	pool       descriptor.Pool
	pass       renderpass.T
}

func NewMaterialSorter(target vulkan.Target, pool descriptor.Pool, pass renderpass.T, defaultMat *material.Def) *MaterialSorter {
	ms := &MaterialSorter{
		target:     target,
		pool:       pool,
		pass:       pass,
		defaultMat: defaultMat,
		cache:      map[uint64]material.Standard{},
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
	mat := material.FromDef(m.target.Device(), m.pool, m.pass, def)
	m.cache[id] = mat
}

func (m *MaterialSorter) Draw(cmds command.Recorder, args render.Args, forwardMeshes []mesh.T) {
	camera := uniform.Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         args.Position,
	}

	// sort meshes by material
	meshes := map[uint64][]mesh.T{}
	for _, msh := range forwardMeshes {
		mid := msh.MaterialID()
		if _, exists := m.cache[mid]; !exists {
			// initialize material
			m.Load(msh.Material())
		}
		meshes[mid] = append(meshes[mid], msh)
	}

	index := 0
	for mid, meshes := range meshes {
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
