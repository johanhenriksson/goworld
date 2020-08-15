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

func NewObjModel(parent *engine.Object, path string, mat *render.Material) *ObjModel {
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
	file, err := gwob.NewObjFromFile(obj.Path, nil)
	if err != nil {
		return fmt.Errorf("parse error input=%s: %v", obj.Path, err)
	}

	// vertex data
	vertexSize := 8 // 8 floats per vertex (xyz uv NxNyNz)
	vertices := make(DefaultVertices, file.NumberOfElements())
	for i := range vertices {
		idx := i * vertexSize
		vertices[i] = DefaultVertex{
			X:  file.Coord[idx+0],
			Y:  file.Coord[idx+1],
			Z:  file.Coord[idx+2],
			U:  file.Coord[idx+3],
			V:  1.0 - file.Coord[idx+4],
			Nx: file.Coord[idx+5],
			Ny: file.Coord[idx+6],
			Nz: file.Coord[idx+7],
		}
	}
	fmt.Println("Vertices", len(vertices))

	// index data
	indices := make(render.UInt32Array, len(file.Indices))
	fmt.Println("Indices", len(indices))
	for i, index := range file.Indices {
		indices[i] = uint32(index)
		fmt.Println(vertices[index])
	}

	obj.Buffer("geometry", vertices)
	obj.BufferIndices(indices)
	return nil
}
