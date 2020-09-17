package geometry

import (
	"fmt"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render"
	"github.com/udhos/gwob"
)

type ObjModel struct {
	*engine.IndexedMesh
	Path string
}

func NewObjModel(parent *engine.Object, mat *render.Material, path string) *ObjModel {
	obj := &ObjModel{
		IndexedMesh: engine.NewIndexedMesh(mat),
		Path:        path,
	}
	if err := obj.load(); err != nil {
		fmt.Println("Error loading model", path, ":", err)
	}

	obj.ComponentBase = engine.NewComponent(parent, obj)
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
	indices := make(render.UInt32Buffer, len(file.Indices))
	for i, index := range file.Indices {
		indices[i] = uint32(index)
	}

	obj.Buffer("geometry", meshdata)
	obj.BufferIndices(indices)
	return nil
}