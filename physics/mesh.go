package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"

import (
	"log"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	shapeBase
	object.T
	meshHandle C.goTriangleMeshHandle
}

var _ Shape = &Mesh{}
var _ mesh.UpdateHandler = &Mesh{}

func NewMesh() *Mesh {
	shape := object.New(&Mesh{
		shapeBase: shapeBase{
			kind: MeshShape,
		},
	})

	runtime.SetFinalizer(shape, func(m *Mesh) {
		m.destroy()
	})

	return shape
}

func (m *Mesh) SetMeshData(mesh vertex.Mesh) {
	// delete any existing physics mesh
	m.destroy()

	vertices := mesh.VertexData()
	vertexArray := reflect.ValueOf(vertices)
	if vertexArray.Kind() != reflect.Slice {
		panic("vertex data is not a slice")
	}
	if vertexArray.Len() < 1 {
		panic("vertex data is empty")
	}
	vertexPtr := vertexArray.Index(0).UnsafeAddr()
	vertexStride := mesh.VertexSize()
	vertexCount := vertexArray.Len()

	indices := mesh.IndexData()
	indexArray := reflect.ValueOf(indices)
	if indexArray.Kind() != reflect.Slice {
		panic("index data is not a slice")
	}
	if indexArray.Len() < 1 {
		panic("index data is empty")
	}
	indexPtr := indexArray.Index(0).UnsafeAddr()
	indexStride := mesh.IndexSize()
	indexCount := indexArray.Len()

	m.meshHandle = C.goNewTriangleMesh(
		unsafe.Pointer(vertexPtr), C.int(vertexCount), C.int(vertexStride),
		unsafe.Pointer(indexPtr), C.int(indexCount), C.int(indexStride))

	m.handle = C.goNewTriangleMeshShape((*C.char)(unsafe.Pointer(m)), m.meshHandle)
}

func (m *Mesh) destroy() {
	// todo: delete mesh handle

	// delete shape
	if m.handle != nil {
		C.goDeleteShape(m.handle)
		m.handle = nil
	}
}

func (m *Mesh) OnActivate() {
	mesh, ok := object.FindInSiblings[mesh.T](m)
	if ok {
		m.SetMeshData(mesh.Mesh())
		log.Println("added mesh data from", m.Parent().Name())
	} else {
		log.Println("no mesh found for collider :(", m.Parent().Name())
	}
}

func (m *Mesh) OnMeshUpdate(mesh vertex.Mesh) {
	log.Println("physics mesh: mesh update")
	m.SetMeshData(mesh)

	// todo: recreate/update the rigidbody
}
