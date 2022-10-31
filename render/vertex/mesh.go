package vertex

import (
	"fmt"
	"reflect"
)

var nextID int = 1

type Mesh interface {
	Id() string
	Version() int
	Indices() int
	Vertices() int
	Primitive() Primitive
	Pointers() Pointers
	VertexData() any
	IndexData() any
	IndexSize() int
	VertexSize() int
}

type MutableMesh[V any, I Indices] interface {
	Mesh
	Update(vertices []V, indices []I)
}

type mesh[V any, I Indices] struct {
	id         string
	version    int
	indexsize  int
	vertexsize int
	primitive  Primitive
	pointers   Pointers
	vertices   []V
	indices    []I
}

func (m *mesh[V, I]) Id() string           { return m.id }
func (m *mesh[V, I]) Version() int         { return m.version }
func (m *mesh[V, I]) IndexData() any       { return m.indices }
func (m *mesh[V, I]) VertexData() any      { return m.vertices }
func (m *mesh[V, I]) IndexSize() int       { return m.indexsize }
func (m *mesh[V, I]) Primitive() Primitive { return m.primitive }
func (m *mesh[V, I]) Pointers() Pointers   { return m.pointers }
func (m *mesh[V, I]) Indices() int         { return len(m.indices) }
func (m *mesh[V, I]) Vertices() int        { return len(m.vertices) }
func (m *mesh[V, I]) VertexSize() int      { return m.vertexsize }

func (m *mesh[V, I]) Update(vertices []V, indices []I) {
	if len(indices) == 0 {
		indices = make([]I, len(vertices))
		for i := 0; i < len(indices); i++ {
			indices[i] = I(i)
		}
	}

	m.vertices = vertices
	m.indices = indices
	m.version++
}

func NewMesh[V any, I Indices](id string, primitive Primitive, vertices []V, indices []I) MutableMesh[V, I] {
	id = fmt.Sprintf("%s-%d", id, nextID)
	nextID++

	var vertex V
	var index I
	ptrs := ParsePointers(vertex)
	mesh := &mesh[V, I]{
		id:         id,
		pointers:   ptrs,
		vertexsize: ptrs.Stride(),
		indexsize:  int(reflect.TypeOf(index).Size()),
		primitive:  primitive,
	}
	mesh.Update(vertices, indices)
	return mesh
}

func NewTriangles[V any, I Indices](id string, vertices []V, indices []I) MutableMesh[V, I] {
	return NewMesh(id, Triangles, vertices, indices)
}

func NewLines[T any, K Indices](id string, vertices []T, indices []K) MutableMesh[T, K] {
	return NewMesh(id, Lines, vertices, indices)
}
