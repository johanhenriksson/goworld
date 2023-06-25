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
	shapeBase
	meshHandle C.goTriangleMeshHandle
}

var _ Shape = &Mesh{}

func NewMesh(mesh vertex.Mesh) *Mesh {
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

	meshHandle := C.goNewTriangleMesh(
		unsafe.Pointer(vertexPtr), C.int(vertexCount), C.int(vertexStride),
		unsafe.Pointer(indexPtr), C.int(indexCount), C.int(indexStride))

	shape := &Mesh{
		shapeBase: shapeBase{
			kind: MeshShape,
		},
		meshHandle: meshHandle,
	}

	shape.handle = C.goNewTriangleMeshShape((*C.char)(unsafe.Pointer(shape)), meshHandle)

	runtime.SetFinalizer(shape, func(m *Mesh) {
		C.goDeleteShape(m.handle)
	})
	return shape
}
