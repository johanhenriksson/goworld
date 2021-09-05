package geometry

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/render"
	"github.com/udhos/gwob"
)

type ObjModel struct {
	mesh.T
	Path string
}

func NewObjModel(mat *render.Material, path string) *ObjModel {
	obj := &ObjModel{
		T:    mesh.New(mat),
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

	// flip texcoord Y
	elsize := len(file.Coord) / file.NumberOfElements()
	for i := 0; i < file.NumberOfElements(); i++ {
		vOffset := i*elsize + 4
		file.Coord[vOffset] = 1.0 - file.Coord[vOffset]
	}

	// vertex data
	meshdata := mesh.NewData(file.NumberOfElements(), file.Coord)

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
