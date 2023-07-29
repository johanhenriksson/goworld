package lines

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

var Debug = &DebugLines{}

type DebugLines struct {
	enabled bool
	frame   int
	meshes  []*Mesh
}

func (li *DebugLines) Setup(frames int) {
	li.meshes = make([]*Mesh, frames)
	for i := range li.meshes {
		li.meshes[i] = New(Args{})
	}
	li.enabled = true
}

func (li *DebugLines) Add(from, to vec3.T, clr color.T) {
	if !li.enabled {
		return
	}
	mesh := li.meshes[li.frame]
	mesh.Lines = append(mesh.Lines, Line{
		Start: from,
		End:   to,
		Color: clr,
	})
}

func (li *DebugLines) Fetch() mesh.Mesh {
	// build mesh for current frame
	mesh := li.meshes[li.frame]
	mesh.Refresh()

	// set next frame
	li.frame = (li.frame + 1) % len(li.meshes)

	// prepare next mesh
	nextMesh := li.meshes[li.frame]
	nextMesh.Clear()

	return mesh
}
