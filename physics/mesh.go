package physics

import (
	"log"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	kind ShapeType
	*Collider

	collision  vertex.Mesh
	meshHandle meshHandle

	Mesh *object.Property[vertex.Mesh]
}

var _ = checkShape(NewMesh())

var emptyMesh = vertex.NewTriangles[vertex.P, uint32]("emptyMeshCollider", []vertex.P{{}}, []uint32{0, 0, 0})

func NewMesh() *Mesh {
	mesh := object.NewComponent(&Mesh{
		kind: MeshShape,

		Mesh: object.NewProperty[vertex.Mesh](nil),
	})

	mesh.Collider = newCollider(mesh)

	// refresh physics mesh when the mesh property is changed
	// unsub to old mesh?
	// subscribe to new mesh?
	mesh.Mesh.OnChange.Subscribe(func(m vertex.Mesh) {
		mesh.refresh()
	})

	return mesh
}

func (m *Mesh) colliderCreate() shapeHandle {
	// generate an optmized collision mesh from the given mesh
	mesh := m.Mesh.Get()
	if mesh == nil {
		mesh = emptyMesh
	}

	m.collision = vertex.CollisionMesh(mesh)
	m.meshHandle = mesh_new(m.collision)

	return shape_new_triangle_mesh(unsafe.Pointer(m), m.meshHandle)
}

func (m *Mesh) colliderDestroy() {
	if m.meshHandle != nil {
		mesh_delete(&m.meshHandle)
	}
}
func (m *Mesh) colliderIsCompound() bool {
	return defaultCompoundCheck(m)
}

func (m *Mesh) Name() string {
	return "MeshShape"
}

func (m *Mesh) OnEnable() {
	if m.Mesh.Get() == nil {
		if mesh := object.Get[mesh.Mesh](m); mesh != nil {
			log.Println("added mesh data from", m.Parent().Name())
			m.Mesh.Set(mesh.Mesh().Get())
			// subscribe?
		} else {
			log.Println("no mesh found for collider :(", m.Parent().Name())
		}
	}

	m.Collider.OnEnable()
}
