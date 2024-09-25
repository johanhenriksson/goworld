package pass

import (
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
	}
}

// RenderGroup is a batch of objects that share the same material.
// All the objects in a group will be rendered using a single indirect draw call.
type RenderGroup struct {
	Pipeline *Pipeline
	Objects  []RenderObject
}

func (m *RenderGroup) Clear() {
	m.Objects = m.Objects[:0]
}

func (m *RenderGroup) Add(mat *Pipeline, object RenderObject) {
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
	for i, g := range r.groups {
		g.Clear()
		r.groups[i] = g
	}
}

// Add an object to the end of the render plan
// If the material is already in the last item in the plan, the object will be added to the existing group
func (r *RenderPlan) AddOrdered(pipe *Pipeline, object RenderObject) {
	index := len(r.groups)
	if index > 0 && r.groups[index-1].Pipeline.id == pipe.id {
		r.groups[index-1].Add(pipe, object)
		return
	}
	r.groups = append(r.groups, RenderGroup{
		Pipeline: nil,
		Objects:  make([]RenderObject, 0, 32),
	})
	r.mapping[pipe.id] = index
	r.groups[index].Add(pipe, object)
}

// Add an object to the render plan
// If the material is already in the plan, the object will be added to the existing group
// Otherwise a new group will be created
func (r *RenderPlan) Add(pipe *Pipeline, object RenderObject) {
	if index, exists := r.mapping[pipe.id]; exists {
		r.groups[index].Add(pipe, object)
		return
	}
	r.AddOrdered(pipe, object)
}

func (r *RenderPlan) Groups() []RenderGroup {
	return r.groups
}
