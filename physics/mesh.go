package physics

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	object.Register[*Mesh](object.TypeInfo{
		Name:        "Mesh Collider",
		Path:        []string{"Physics"},
		Deserialize: DeserializeMesh,
		Create: func(ctx object.Pool) (object.Component, error) {
			return NewMesh(ctx), nil
		},
	})
}

type Mesh struct {
	kind ShapeType
	*Collider

	collision  vertex.Mesh
	meshHandle meshHandle
	unsub      func()

	Mesh object.Property[vertex.Mesh]
}

var _ = checkShape(NewMesh(object.GlobalPool))

var emptyMesh = vertex.NewTriangles[vertex.P, uint32]("emptyMeshCollider", []vertex.P{{}}, []uint32{0, 0, 0})

func NewMesh(pool object.Pool) *Mesh {
	mesh := &Mesh{
		kind: MeshShape,
		Mesh: object.NewProperty[vertex.Mesh](nil),
	}
	mesh.Collider = newCollider(pool, mesh, true)

	// refresh physics mesh when the mesh property is changed
	mesh.Mesh.OnChange.Subscribe(func(m vertex.Mesh) {
		mesh.refresh()
	})

	return object.NewComponent(pool, mesh)
}

func (m *Mesh) colliderCreate() shapeHandle {
	// generate an optmized collision mesh from the given mesh
	mesh := m.Mesh.Get()
	if mesh == nil || mesh.IndexCount() == 0 {
		mesh = emptyMesh
	}

	m.collision = vertex.CollisionMesh(mesh)
	m.meshHandle = mesh_new(m.collision)

	return shape_new_triangle_mesh(unsafe.Pointer(m), m.meshHandle)
}

func (m *Mesh) colliderRefresh() {}

func (m *Mesh) colliderIsCompound() bool { return false }

func (m *Mesh) colliderDestroy() {
	if m.meshHandle != nil {
		mesh_delete(&m.meshHandle)
	}
}

func (m *Mesh) Name() string {
	return "MeshShape"
}

func (m *Mesh) OnEnable() {
	if m.Mesh.Get() == nil {
		if mesh := object.Get[mesh.Mesh](m); mesh != nil {
			m.Mesh.Set(mesh.Mesh().Get())
			if m.unsub != nil {
				m.unsub()
				m.unsub = nil
			}
			m.unsub = mesh.Mesh().OnChange.Subscribe(func(update vertex.Mesh) {
				m.Mesh.Set(update)
			})
		}
	}
	m.Collider.OnEnable()
}

func (m *Mesh) OnDisable() {
	if m.unsub != nil {
		m.unsub()
		m.unsub = nil
	}
	m.Collider.OnDisable()
}

func (m *Mesh) Serialize(enc object.Encoder) error {
	return nil
}

func DeserializeMesh(ctx object.Pool, dec object.Decoder) (object.Component, error) {
	return NewMesh(ctx), nil
}
