package geometry

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render"
	"github.com/udhos/gwob"
)

type ObjModel struct {
	*engine.Mesh
	Path string
}

func NewObjModel(mat *render.Material, path string) *ObjModel {
	obj := &ObjModel{
		Mesh: engine.NewMesh(fmt.Sprintf("Mesh:%s", path), mat),
		Path: path,
	}
	obj.SetIndexType(render.UInt32)
	if err := obj.load(); err != nil {
		fmt.Println("Error loading model", path, ":", err)
	}
	return obj
}

func (obj *ObjModel) load() error {
	// load obj
	assetPath := fmt.Sprintf("assets/%s", obj.Path)
	file, err := gwob.NewObjFromFile(assetPath, nil)
	if err != nil {
		return fmt.Errorf("parse error input=%s: %v", obj.Path, err)
	}

	// vertex data
	meshdata := &engine.MeshData{
		Items:  file.NumberOfElements(),
		Buffer: file.Coord,
	}

	// flip texcoord Y
	for i := 0; i < meshdata.Items; i++ {
		vOffset := i*meshdata.Size()/4 + 4
		meshdata.Buffer[vOffset] = 1.0 - meshdata.Buffer[vOffset]
	}

	// index data
	indices := make(UInt32Buffer, len(file.Indices))
	for i, index := range file.Indices {
		indices[i] = uint32(index)
	}

	obj.Buffer(meshdata)
	// obj.Buffer("index", indices)
	return nil
}

type UInt32Buffer []uint32

func (a UInt32Buffer) Elements() int {
	return len(a)
}

func (a UInt32Buffer) Size() int {
	return 4
}

func (a UInt32Buffer) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&a[0])
}
