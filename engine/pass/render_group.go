package pass

import (
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
)

type RenderObject struct {
	// Handle is the index of the object in the object storage buffer
	Handle int

	// Indices is the number of indices to render
	Indices int
}

// DrawIndirect returns a command.Draw object that can be used to render the object
func (r RenderObject) DrawIndirect() command.Draw {
	return command.Draw{
		InstanceCount: 1,

		// InstanceOffset is the index of the object properties in the object buffer
		InstanceOffset: uint32(r.Handle),

		// Vertex count is actually the number of indices, since indexing is implemented in the shader
		VertexCount: uint32(r.Indices),

		VertexOffset: 0,
	}
}

// RenderGroup is a batch of objects that share the same material.
// All the objects in a group will be rendered using a single indirect draw call.
type RenderGroup struct {
	Pipeline *cache.Pipeline
	Objects  []RenderObject
}

func (m *RenderGroup) Clear() {
	clear(m.Objects)
	m.Objects = m.Objects[:0]
}

func (m *RenderGroup) Add(mat *cache.Pipeline, object RenderObject) {
	m.Pipeline = mat
	m.Objects = append(m.Objects, object)
}

// RenderPlan is a collection of render groups, each group containing objects that share the same material.
// The render plan maintains the ordering of object batches
type RenderPlan struct {
	groups  []RenderGroup
	mapping map[material.ID]int
}

func NewRenderPlan() *RenderPlan {
	return &RenderPlan{
		groups:  make([]RenderGroup, 0, 32),
		mapping: make(map[material.ID]int, 32),
	}
}

// Clear the rendre plan, preserving the allocated memory
func (r *RenderPlan) Clear() {
	clear(r.mapping)
	clear(r.groups)
	r.groups = r.groups[:0]
}

// Add an object to the end of the render plan
// If the material is already in the last item in the plan, the object will be added to the existing group
func (r *RenderPlan) AddOrdered(pipe *cache.Pipeline, object RenderObject) {
	index := len(r.groups)
	if index > 0 && r.groups[index-1].Pipeline.ID == pipe.ID {
		r.groups[index-1].Add(pipe, object)
		return
	}
	r.groups = append(r.groups, RenderGroup{
		Pipeline: nil,
		Objects:  make([]RenderObject, 0, 32),
	})
	r.mapping[pipe.ID] = index
	r.groups[index].Add(pipe, object)
}

// Add an object to the render plan
// If the material is already in the plan, the object will be added to the existing group
// Otherwise a new group will be created
func (r *RenderPlan) Add(pipe *cache.Pipeline, object RenderObject) {
	if index, exists := r.mapping[pipe.ID]; exists {
		r.groups[index].Add(pipe, object)
		return
	}
	r.AddOrdered(pipe, object)
}

func (r *RenderPlan) Draw(cmd *command.Buffer, indirect *command.IndirectDrawBuffer) {
	indirect.Reset()
	for _, group := range r.groups {
		if len(group.Objects) == 0 {
			continue
		}
		group.Pipeline.Bind(cmd)
		indirect.BeginDrawIndirect()
		for _, obj := range group.Objects {
			indirect.CmdDraw(obj.DrawIndirect())
		}
		indirect.EndDrawIndirect(cmd)
	}
}
