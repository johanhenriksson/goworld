package vertex

import (
	"iter"
	"reflect"

	"github.com/johanhenriksson/goworld/assets/fs"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Mesh interface {
	Key() string
	Version() int

	Primitive() Primitive
	Pointers() Pointers
	VertexCount() int
	VertexData() any
	VertexSize() int
	IndexCount() int
	IndexData() any
	IndexSize() int

	Min() vec3.T
	Max() vec3.T

	// Bounds returns a bounding sphere containing the mesh, centered at the given origin.
	Bounds(vec3.T) shape.Sphere

	Positions() iter.Seq[vec3.T]
	Triangles() iter.Seq[Triangle]

	// Self-referential
	LoadMesh(fs.Filesystem) Mesh
}

type Vertex interface {
	Position() vec3.T
}

type Index interface {
	uint8 | uint16 | uint32
}

type MutableMesh[V Vertex, I Index] interface {
	Mesh
	Vertices() []V
	Indices() []I
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
	center     vec3.T
	radius     float32
}

var _ Mesh = &mesh[P, uint8]{}

func (m *mesh[V, I]) Key() string          { return m.key }
func (m *mesh[V, I]) Version() int         { return m.version }
func (m *mesh[V, I]) Primitive() Primitive { return m.primitive }
func (m *mesh[V, I]) Pointers() Pointers   { return m.pointers }
func (m *mesh[V, I]) Vertices() []V        { return m.vertices }
func (m *mesh[V, I]) VertexData() any      { return m.vertices }
func (m *mesh[V, I]) VertexSize() int      { return m.vertexsize }
func (m *mesh[V, I]) VertexCount() int     { return len(m.vertices) }
func (m *mesh[V, I]) Indices() []I         { return m.indices }
func (m *mesh[V, I]) IndexData() any       { return m.indices }
func (m *mesh[V, I]) IndexSize() int       { return m.indexsize }
func (m *mesh[V, I]) IndexCount() int      { return len(m.indices) }
func (m *mesh[V, I]) String() string       { return m.key }
func (m *mesh[V, I]) Min() vec3.T          { return m.min }
func (m *mesh[V, I]) Max() vec3.T          { return m.max }

func (m *mesh[V, I]) Positions() iter.Seq[vec3.T] {
	return func(yield func(vec3.T) bool) {
		for _, index := range m.indices {
			vertex := m.vertices[index]
			if !yield(vertex.Position()) {
				break
			}
		}
	}
}

func (m *mesh[V, I]) Triangles() iter.Seq[Triangle] {
	return func(yield func(Triangle) bool) {
		for i := 0; i+3 < len(m.indices); i += 3 {
			if !yield(Triangle{
				A: m.vertices[m.indices[i+0]].Position(),
				B: m.vertices[m.indices[i+1]].Position(),
				C: m.vertices[m.indices[i+2]].Position(),
			}) {
				break
			}
		}
	}
}

func (m *mesh[V, I]) Bounds(origin vec3.T) shape.Sphere {
	return shape.Sphere{
		Center: origin.Add(m.center),
		Radius: m.radius,
	}
}

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
	m.center = m.max.Sub(m.min).Scaled(0.5)
	m.radius = m.center.Length()

	m.vertices = vertices
	m.indices = indices
	m.version++
}

// LoadMesh returns the mesh itself.
// Implementing LoadMesh ensures that the mesh satisfies the assets.Mesh interface.
func (m *mesh[V, I]) LoadMesh(fs.Filesystem) Mesh {
	return m
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
