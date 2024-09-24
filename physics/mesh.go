package physics

import (
	"log"
	"unsafe"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/vertex"
)

func init() {
	object.Register[*Mesh](object.Type{
		Name: "Mesh Collider",
		Path: []string{"Physics"},
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

	Mesh object.Property[assets.Mesh]

	key     string
	version int
}

var _ = checkShape(NewMesh(object.GlobalPool))

var emptyMesh = vertex.NewTriangles[Vertex, uint32]("emptyMeshCollider", []Vertex{{}}, []uint32{0, 0, 0})

func NewMesh(pool object.Pool) *Mesh {
	mesh := &Mesh{
		kind: MeshShape,
		Mesh: object.NewProperty[assets.Mesh](nil),
	}
	mesh.Collider = newCollider(pool, mesh, true)

	// refresh physics mesh when the mesh property is changed
	mesh.Mesh.OnChange.Subscribe(func(m assets.Mesh) {
		mesh.checkMesh()
	})

	return object.NewComponent(pool, mesh)
}

func (m *Mesh) colliderCreate() shapeHandle {
	// generate an optmized collision mesh from the given mesh
	ref := m.Mesh.Get()
	var mesh vertex.Mesh = emptyMesh
	if ref != nil {
		mesh = ref.LoadMesh(assets.FS)
	}
	log.Println("creating physics mesh", m.key, "tris:", mesh.IndexCount()/3)

	m.collision = CollisionMesh(mesh)
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
	log.Println("physics mesh enabled", m.key)
	if m.Mesh.Get() == nil {
		if mesh := object.Get[mesh.Mesh](m); mesh != nil {
			m.Mesh.Set(mesh.Mesh())
		} else {
			log.Println("mesh not found")
		}
	} else {
		log.Println("mesh is set")
	}
	m.Collider.OnEnable()
}

func (m *Mesh) checkMesh() {
	// watch for changes to the ref, and trigger a refresh if needed
	if ref := m.Mesh.Get(); ref != nil {
		if m.key != ref.Key() || m.version != ref.Version() {
			log.Println("physics mesh updated", m.key)
			m.key = ref.Key()
			m.version = ref.Version()
			m.refresh()
		}
	} else if m.key != "" {
		log.Println("physics mesh removed", m.key)
		m.key = ""
		m.version = 0
		m.Mesh.Set(nil)
	}
}

func (m *Mesh) Update(scene object.Component, dt float32) {
	m.checkMesh()
	m.Collider.Update(scene, dt)
}

func (m *Mesh) EditorUpdate(scene object.Component, dt float32) {
	m.checkMesh()
}

func (m *Mesh) OnDisable() {
	m.Collider.OnDisable()
}
