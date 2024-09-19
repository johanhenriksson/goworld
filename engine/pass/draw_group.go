package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
)

type Drawable interface {
	Key() string
	Version() int

	Bind(cmd *command.Buffer)
	Draw(cmd command.DrawIndexedBuffer, instanceOffset int)
	DrawInstanced(cmd command.DrawIndexedBuffer, instanceOffset, instanceCount int)

	Model() mat4.T
	MaterialID() material.ID
	Material() *material.Def
	Textures() uniform.TextureIds

	Bounds() shape.Sphere
}

type DrawableMesh struct {
	*cache.GpuMesh
	model    mat4.T
	material *material.Def
	textures uniform.TextureIds
}

func (m DrawableMesh) Model() mat4.T                  { return m.model }
func (m DrawableMesh) Material() *material.Def        { return m.material }
func (m DrawableMesh) MaterialID() material.ID        { return m.material.Hash() }
func (m DrawableMesh) TextureIds() uniform.TextureIds { return m.textures }

func (m DrawableMesh) Bounds() shape.Sphere {
	b := m.GpuMesh.Bounds()
	return shape.Sphere{
		Center: b.Center.Add(m.model.Origin()),
		Radius: b.Radius,
	}
}

type DrawGroup struct {
	ID       material.ID
	Material Material
	Meshes   []mesh.Mesh
}

type DrawGroups []DrawGroup

func (groups DrawGroups) Draw(cmds command.Recorder, camera uniform.Camera, lights []light.T) {
	// todo: this can happen multiple times per frame if there are multiple draw groups for the same material.
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

	// todo: this can happen multiple times per frame if there are multiple draw groups for the same material.
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
