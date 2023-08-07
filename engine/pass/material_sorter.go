package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type MeshSorter[T any] struct {
	cache *MaterialCache[T]
	maker MaterialMaker[T]
	app   vulkan.App
}

type MeshGroup[T any] struct {
	Material T
	Meshes   []mesh.Mesh
}

func NewMeshSorter[T any](app vulkan.App, frames int, maker MaterialMaker[T]) *MeshSorter[T] {
	return &MeshSorter[T]{
		app:   app,
		cache: NewMaterialCache[T](app, frames, maker),
		maker: maker,
	}
}

func (m *MeshSorter[T]) Destroy() {
	m.cache.Destroy()
	m.cache = nil
}

func (m *MeshSorter[T]) Draw(cmds command.Recorder, args render.Args, meshes []mesh.Mesh, lights []light.T) {
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

func (m *MeshSorter[T]) DrawCamera(cmds command.Recorder, frame int, camera uniform.Camera, meshes []mesh.Mesh, lights []light.T) {
	// sort meshes by material
	meshGroups := map[uint64]*MeshGroup[T]{}
	for _, msh := range meshes {
		mat, ready := m.cache.Get(msh, frame)
		if !ready {
			continue
		}
		group, exists := meshGroups[msh.MaterialID()]
		if !exists {
			group = &MeshGroup[T]{
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
