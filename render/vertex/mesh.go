package vertex

import (
	"reflect"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type Mesh interface {
	Key() string
	Version() int
	Indices() int
	Vertices() int
	Primitive() Primitive
	Pointers() Pointers
	VertexData() any
	IndexData() any
	IndexSize() int
	VertexSize() int
	Min() vec3.T
	Max() vec3.T
}

type Vertex interface {
	Position() vec3.T
}

type Index interface {
	uint8 | uint16 | uint32
}

type MutableMesh[V Vertex, I Index] interface {
	Mesh
	Update(vertices []V, indices []I)
}

type mesh[V Vertex, I Index] struct {
	key        string
	version    int
	indexsize  int
	vertexsize int
	primitive  Primitive
	pointers   Pointers
	vertices   []V
	indices    []I
	min        vec3.T
	max        vec3.T
}

func (m *mesh[V, I]) Key() string          { return m.key }
func (m *mesh[V, I]) Version() int         { return m.version }
func (m *mesh[V, I]) IndexData() any       { return m.indices }
func (m *mesh[V, I]) VertexData() any      { return m.vertices }
func (m *mesh[V, I]) IndexSize() int       { return m.indexsize }
func (m *mesh[V, I]) Primitive() Primitive { return m.primitive }
func (m *mesh[V, I]) Pointers() Pointers   { return m.pointers }
func (m *mesh[V, I]) Indices() int         { return len(m.indices) }
func (m *mesh[V, I]) Vertices() int        { return len(m.vertices) }
func (m *mesh[V, I]) VertexSize() int      { return m.vertexsize }
func (m *mesh[V, I]) String() string       { return m.key }
func (m *mesh[V, I]) Min() vec3.T          { return m.min }
func (m *mesh[V, I]) Max() vec3.T          { return m.max }

func (m *mesh[V, I]) Update(vertices []V, indices []I) {
	if len(indices) == 0 {
		indices = make([]I, len(vertices))
		for i := 0; i < len(indices); i++ {
			indices[i] = I(i)
		}
	}

	// update mesh bounds
	m.min = Min(vertices)
	m.max = Max(vertices)

	m.vertices = vertices
	m.indices = indices
	m.version++
}

func NewMesh[V Vertex, I Index](key string, primitive Primitive, vertices []V, indices []I) MutableMesh[V, I] {
	var vertex V
	var index I
	ptrs := ParsePointers(vertex)

	// calculate mesh bounds
	min := Min(vertices)
	max := Max(vertices)

	mesh := &mesh[V, I]{
		key:        key,
		pointers:   ptrs,
		vertexsize: ptrs.Stride(),
		indexsize:  int(reflect.TypeOf(index).Size()),
		primitive:  primitive,
		min:        min,
		max:        max,
	}
	mesh.Update(vertices, indices)
	return mesh
}

func NewTriangles[V Vertex, I Index](key string, vertices []V, indices []I) MutableMesh[V, I] {
	return NewMesh(key, Triangles, vertices, indices)
}

func NewLines[T Vertex, K Index](key string, vertices []T, indices []K) MutableMesh[T, K] {
	return NewMesh(key, Lines, vertices, indices)
}
