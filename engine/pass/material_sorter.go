package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type MaterialTransform func(*material.Def) *material.Def

type MaterialSorter struct {
	cache *MaterialCache[*MaterialData]
	maker MaterialMaker[*MaterialData]
	app   vulkan.App
}

type MaterialData struct {
	// shared between deferred & forward rendering:

	Instance material.Instance[*material.Descriptors]
	Objects  []uniform.Object
	Textures cache.SamplerCache

	// extra data needed for forward rendering:

	Lights  *LightBuffer
	Shadows *ShadowCache
}

func NewMaterialSorter(app vulkan.App, frames int, pass renderpass.T, lookup ShadowmapLookupFn, defaultMat *material.Def, transform MaterialTransform) *MaterialSorter {
	maker := &StdMaterialMaker{
		app:       app,
		pass:      pass,
		lookup:    lookup,
		transform: transform,
	}
	return &MaterialSorter{
		app:   app,
		cache: NewMaterialCache[*MaterialData](app, frames, pass, maker, defaultMat),
		maker: maker,
	}
}

func (m *MaterialSorter) Destroy() {
	m.cache.Destroy()
	m.cache = nil
}

func (m *MaterialSorter) Draw(cmds command.Recorder, args render.Args, meshes []mesh.Mesh, lights []light.T) {
	camera := uniform.Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         vec4.Extend(args.Position, 0),
		Forward:     vec4.Extend(args.Forward, 0),
		Viewport:    vec2.NewI(args.Viewport.Width, args.Viewport.Height),
	}
	m.DrawCamera(cmds, args.Context.Index, camera, meshes, lights)
}

type MeshGroup[T any] struct {
	Material T
	Meshes   []mesh.Mesh
}

func (m *MaterialSorter) DrawCamera(cmds command.Recorder, frame int, camera uniform.Camera, meshes []mesh.Mesh, lights []light.T) {
	// sort meshes by material
	meshGroups := map[uint64]*MeshGroup[*MaterialData]{}
	for _, msh := range meshes {
		mat, ready := m.cache.Get(msh, frame)
		if !ready {
			continue
		}
		group, exists := meshGroups[msh.MaterialID()]
		if !exists {
			group = &MeshGroup[*MaterialData]{
				Material: mat,
				Meshes:   make([]mesh.Mesh, 0, 32),
			}
			meshGroups[msh.MaterialID()] = group
		}
		group.Meshes = append(group.Meshes, msh)
	}

	// iterate the sorted material groups
	for _, group := range meshGroups {
		m.maker.Draw(cmds, camera, group, lights)
	}
}
