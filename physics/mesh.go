package physics

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	shapeBase
	object.Component

	compound   bool
	collision  vertex.Mesh
	meshHandle meshHandle
	unsubTf    func()

	Mesh *object.Property[vertex.Mesh]
}

var _ Shape = &Mesh{}

var emptyMesh = vertex.NewTriangles[vertex.P, uint32]("emptyMeshCollider", []vertex.P{{}}, []uint32{0, 0, 0})

func NewMesh() *Mesh {
	mesh := object.NewComponent(&Mesh{
		shapeBase: newShapeBase(MeshShape),

		Mesh: object.NewProperty[vertex.Mesh](nil),
	})

	// refresh physics mesh when the mesh property is changed
	// unsub to old mesh?
	// subscribe to new mesh?
	mesh.Mesh.OnChange.Subscribe(mesh.refresh)

	// initialize with the empty mesh
	mesh.meshHandle = mesh_new(emptyMesh)
	mesh.handle = shape_new_triangle_mesh(unsafe.Pointer(mesh), mesh.meshHandle)

	runtime.SetFinalizer(mesh, func(m *Mesh) {
		m.destroy()
	})

	return mesh
}

func (m *Mesh) Name() string {
	return "MeshShape"
}

func (m *Mesh) scale() vec3.T {
	if m.compound {
		return m.Transform().Scale()
	}
	return m.Transform().WorldScale()
}

func (m *Mesh) refresh(mesh vertex.Mesh) {
	// todo: if its the same mesh, dont do anything

	// delete any existing physics mesh
	m.destroy()

	// generate an optmized collision mesh from the given mesh
	m.collision = vertex.CollisionMesh(mesh)
	log.Println("computed collision mesh of", m.collision.VertexCount(), "vertices down from", mesh.VertexCount(), "[", mesh.IndexCount(), "]")

	m.meshHandle = mesh_new(m.collision)

	m.handle = shape_new_triangle_mesh(unsafe.Pointer(m), m.meshHandle)
	shape_scaling_set(m.handle, m.scale())
	m.OnChange().Emit(m)
}

func (m *Mesh) destroy() {
	// todo: delete mesh handle
	if m.meshHandle != nil {
		mesh_delete(&m.meshHandle)
	}

	// delete shape
	if m.handle != nil {
		shape_delete(&m.handle)
	}
}

func (m *Mesh) OnEnable() {
	log.Println("enable mesh", m.Parent().Name())
	if m.Mesh.Get() == nil {
		if mesh := object.Get[mesh.Mesh](m); mesh != nil {
			log.Println("added mesh data from", m.Parent().Name())
			m.Mesh.Set(mesh.Mesh().Get())
			// subscribe?
		} else {
			log.Println("no mesh found for collider :(", m.Parent().Name())
		}
	}

	// react to scale changes
	lastScale := m.scale()
	shape_scaling_set(m.handle, lastScale)
	m.OnChange().Emit(m)

	m.unsubTf = m.Transform().OnChange().Subscribe(func(t transform.T) {
		if t.Scale() != lastScale {
			lastScale = t.Scale()
			log.Println("update mesh", m.Parent().Name(), "scale to", lastScale)
			shape_scaling_set(m.handle, lastScale)
			m.OnChange().Emit(m)
		}
	})
}

func (m *Mesh) OnDisable() {
	m.unsubTf()
}
