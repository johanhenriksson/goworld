package vertex

import (
	"fmt"
	"reflect"
)

var nextID int = 1

type Mesh interface {
	Id() string
	Version() int
	Elements() int
	Primitive() Primitive
	Pointers() Pointers
	VertexData() any
	IndexData() any
	IndexSize() int
}

type mesh[T any, K Indices] struct {
	id        string
	version   int
	indexsize int
	primitive Primitive
	pointers  Pointers
	vertices  []T
	indices   []K
}

func (m *mesh[T, K]) Id() string           { return m.id }
func (m *mesh[T, K]) Version() int         { return m.version }
func (m *mesh[T, K]) IndexData() any       { return m.indices }
func (m *mesh[T, K]) VertexData() any      { return m.vertices }
func (m *mesh[T, K]) IndexSize() int       { return m.indexsize }
func (m *mesh[T, K]) Primitive() Primitive { return m.primitive }
func (m *mesh[T, K]) Pointers() Pointers   { return m.pointers }

func (m *mesh[T, K]) Elements() int {
	if len(m.indices) > 0 {
		return len(m.indices)
	}
	return len(m.vertices)
}

func (m *mesh[T, K]) Update(vertices []T, indices []K) {
	m.vertices = vertices
	m.indices = indices
	m.version++
}

func NewMesh[T any, K Indices](id string, primitive Primitive, vertices []T, indices []K) Mesh {
	id = fmt.Sprintf("%s-%d", id, nextID)
	nextID++

	if len(indices) == 0 {
		indices = make([]K, len(vertices))
		for i := 0; i < len(indices); i++ {
			indices[i] = K(i)
		}
	}
	var vertex T
	var index K
	return &mesh[T, K]{
		id:        id,
		version:   1,
		pointers:  ParsePointers(vertex),
		vertices:  vertices,
		indices:   indices,
		indexsize: int(reflect.TypeOf(index).Size()),
		primitive: primitive,
	}
}

func NewTriangles[T any, K Indices](id string, vertices []T, indices []K) Mesh {
	return NewMesh(id, Triangles, vertices, indices)
}

func NewLines[T any, K Indices](id string, vertices []T, indices []K) Mesh {
	return NewMesh(id, Lines, vertices, indices)
}
