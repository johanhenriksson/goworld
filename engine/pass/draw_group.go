package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
)

type DrawGroup struct {
	ID       material.ID
	Material Material
	Meshes   []mesh.Mesh
}

type DrawGroups []DrawGroup

func (groups DrawGroups) Draw(cmds command.Recorder, camera uniform.Camera, lights []light.T) {
	for _, group := range groups {
		group.Material.Begin(camera, lights)
	}

	for _, group := range groups {
		group.Material.Bind(cmds)
		for _, msh := range group.Meshes {
			group.Material.Draw(cmds, msh)
		}
		group.Material.Unbind(cmds)
	}

	for _, group := range groups {
		group.Material.End()
	}
}

// Sort meshes by material according to depth.
// Consecutive meshes in the depth order are grouped if they have the same material
func OrderedGroups(cache MaterialCache, frame int, meshes []mesh.Mesh) DrawGroups {
	groups := make(DrawGroups, 0, 16)
	var group *DrawGroup
	for _, msh := range meshes {
		mats, ready := cache.TryFetch(msh.Material())
		if !ready {
			continue
		}

		id := msh.MaterialID()
		if group == nil || id != group.Material.ID() {
			groups = append(groups, DrawGroup{
				Material: mats[frame],
				Meshes:   make([]mesh.Mesh, 0, 32),
			})
			group = &groups[len(groups)-1]
		}
		group.Meshes = append(group.Meshes, msh)
	}
	return groups
}

// Sort meshes by material
func MaterialGroups(cache MaterialCache, frame int, meshes []mesh.Mesh) DrawGroups {
	groups := make(DrawGroups, 0, 16)
	matGroups := map[material.ID]*DrawGroup{}

	for _, msh := range meshes {
		mats, ready := cache.TryFetch(msh.Material())
		if !ready {
			continue
		}

		group, exists := matGroups[msh.MaterialID()]
		if !exists {
			groups = append(groups, DrawGroup{
				Material: mats[frame],
				Meshes:   make([]mesh.Mesh, 0, 32),
			})
			group = &groups[len(groups)-1]
			matGroups[msh.MaterialID()] = group
		}
		group.Meshes = append(group.Meshes, msh)
	}

	return groups
}
