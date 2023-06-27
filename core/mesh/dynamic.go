package mesh

import (
	"log"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Generator[V vertex.Vertex, I vertex.Index] func() Data[V, I]

type Data[V vertex.Vertex, I vertex.Index] struct {
	Vertices []V
	Indices  []I
}

type Dynamic[V vertex.Vertex, I vertex.Index] interface {
	T
	Refresh()
	RefreshSync()
}

type dynamic[V vertex.Vertex, I vertex.Index] struct {
	T
	name     string
	refresh  Generator[V, I]
	updated  chan Data[V, I]
	meshdata vertex.MutableMesh[V, I]
}

func NewDynamic[V vertex.Vertex, I vertex.Index](name string, mode DrawMode, mat *material.Def, fn Generator[V, I]) Dynamic[V, I] {
	m := &dynamic[V, I]{
		T:       New(mode, mat),
		name:    name,
		refresh: fn,
		updated: make(chan Data[V, I], 2),
	}
	m.meshdata = vertex.NewTriangles(object.Key(name, m), []V{}, []I{})
	m.SetMesh(m.meshdata)
	m.RefreshSync()

	return m
}

func (m *dynamic[V, I]) Name() string {
	return m.name
}

func (m *dynamic[V, I]) Refresh() {
	log.Println("mesh", m, ": async refresh")
	go func() {
		data := m.refresh()
		m.updated <- data
	}()
}

func (m *dynamic[V, I]) RefreshSync() {
	log.Println("mesh", m, ": blocking refresh")
	data := m.refresh()
	m.meshdata.Update(data.Vertices, data.Indices)
	m.SetMesh(m.meshdata)
}

func (m *dynamic[V, I]) Update(scene object.T, dt float32) {
	m.T.Update(scene, dt)
	select {
	case data := <-m.updated:
		m.meshdata.Update(data.Vertices, data.Indices)
		m.SetMesh(m.meshdata)
	default:
	}
}
