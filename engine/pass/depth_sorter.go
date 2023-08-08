package pass

import (
	"sort"

	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type DepthSorter struct {
	cache *MaterialCache[*ForwardMatData]
	maker *ForwardMaterialMaker
	app   vulkan.App
}

func NewDepthSorter(app vulkan.App, frames int, maker *ForwardMaterialMaker) *DepthSorter {
	return &DepthSorter{
		app:   app,
		cache: NewMaterialCache[*ForwardMatData](app, frames, maker),
		maker: maker,
	}
}

func (m *DepthSorter) Destroy() {
	m.cache.Destroy()
	m.cache = nil
}

func (m *DepthSorter) Draw(cmds command.Recorder, frame int, camera uniform.Camera, meshes []mesh.Mesh, lights []light.T) {
	// perform back-to-front depth sorting
	// we use the closest point on the meshes bounding sphere as a heuristic
	sort.SliceStable(meshes, func(i, j int) bool {
		// return true if meshes[i] is further away than meshes[j]
		first, second := meshes[i].BoundingSphere(), meshes[j].BoundingSphere()

		di := vec3.Distance(camera.Eye.XYZ(), first.Center) - first.Radius
		dj := vec3.Distance(camera.Eye.XYZ(), second.Center) - second.Radius
		return di > dj
	})

	// sort meshes by material
	meshGroups := []*MeshGroup[*ForwardMatData]{}
	var group *MeshGroup[*ForwardMatData]
	for _, msh := range meshes {
		mat, ready := m.cache.Get(msh, frame)
		if !ready {
			continue
		}

		id := msh.MaterialID()
		if group == nil || id != group.MatID {
			if group != nil && len(group.Meshes) > 0 {
				meshGroups = append(meshGroups, group)
			}
			group = &MeshGroup[*ForwardMatData]{
				MatID:    id,
				Material: mat,
				Meshes:   make([]mesh.Mesh, 0, 32),
			}
		}

		group.Meshes = append(group.Meshes, msh)
	}
	if len(group.Meshes) > 0 {
		meshGroups = append(meshGroups, group)
	}

	for _, group := range meshGroups {
		// can happen multiple times
		m.maker.BeginFrame(group.Material, camera, lights)
	}

	for _, group := range meshGroups {
		mat := group.Material
		cmds.Record(func(cmd command.Buffer) {
			mat.Instance.Bind(cmd)
		})
		for _, msh := range group.Meshes {
			vkmesh, meshReady := m.app.Meshes().TryFetch(msh.Mesh().Get())
			if !meshReady {
				continue
			}

			index := m.maker.PrepareMesh(mat, msh)

			cmds.Record(func(cmd command.Buffer) {
				vkmesh.Draw(cmd, index)
			})
		}
	}

	for _, group := range meshGroups {
		// can happen multiple times
		m.maker.EndFrame(group.Material)
	}
}
