package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	handle     C.goShapeHandle
	meshHandle C.goTriangleMeshHandle
}

var _ Shape = &Mesh{}

func NewMesh(mesh vertex.Mesh) *Mesh {
	vertices := mesh.VertexData()
	vertexValue := reflect.ValueOf(vertices)
	firstVertex := vertexValue.Index(0)
	vertexPtr := firstVertex.Addr().Pointer()
	vertexStride := mesh.VertexSize()
	vertexCount := vertexValue.Len()

	indices := mesh.IndexData()
	indexValue := reflect.ValueOf(indices)
	firstIndex := indexValue.Index(0)
	indexPtr := firstIndex.Addr().Pointer()
	indexStride := mesh.IndexSize()
	indexCount := indexValue.Len()

	meshHandle := C.goNewTriangleMesh(
		unsafe.Pointer(vertexPtr), C.int(vertexCount), C.int(vertexStride),
		unsafe.Pointer(indexPtr), C.int(indexCount), C.int(indexStride))

	handle := C.goNewStaticTriangleMeshShape(meshHandle)

	physMesh := &Mesh{
		handle:     handle,
		meshHandle: meshHandle,
	}
	runtime.SetFinalizer(physMesh, func(m *Mesh) {
		C.goDeleteShape(m.handle)
	})
	return physMesh
}

func (m *Mesh) shape() C.goShapeHandle {
	return m.handle
}
