package physics

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	shapeBase
	object.Component

	collision  vertex.Mesh
	meshHandle meshHandle

	Mesh *object.Property[vertex.Mesh]
}

var _ Shape = &Mesh{}

func NewMesh() *Mesh {
	shape := object.NewComponent(&Mesh{
		shapeBase: newShapeBase(MeshShape),

		Mesh: object.NewProperty[vertex.Mesh](nil),
	})

	// refresh physics mesh when the mesh property is changed
	// unsub to old mesh?
	// subscribe to new mesh?
	shape.Mesh.OnChange.Subscribe(shape, shape.refresh)

	runtime.SetFinalizer(shape, func(m *Mesh) {
		m.destroy()
	})

	return shape
}

func (m *Mesh) refresh(mesh vertex.Mesh) {
	// todo: if its the same mesh, dont do anything

	// delete any existing physics mesh
	m.destroy()

	// generate an optmized collision mesh from the given mesh
	m.collision = vertex.CollisionMesh(mesh)
	log.Println("computed collision mesh of", m.collision.VertexCount(), "vertices [", m.collision.IndexCount(), "], down from", mesh.VertexCount(), "[", mesh.IndexCount(), "]")

	m.meshHandle = mesh_new(m.collision)

	m.handle = shape_new_triangle_mesh(unsafe.Pointer(m), m.meshHandle)
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
	if m.Mesh.Get() != nil {
		// we already have a mesh handle
		return
	}
	if mesh := object.Get[mesh.Mesh](m); mesh != nil {
		log.Println("added mesh data from", m.Parent().Name())
		m.Mesh.Set(mesh.Mesh().Get())
		// subscribe?
	} else {
		log.Println("no mesh found for collider :(", m.Parent().Name())
	}
}
