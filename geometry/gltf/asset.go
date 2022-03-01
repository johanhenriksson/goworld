package gltf

import (
	"fmt"
	"strings"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/qmuntal/gltf"
)

type gltfModel struct {
	mesh.T
	Path string
}

func Load(mat material.T, path string) *gltfModel {
	obj := &gltfModel{
		T:    mesh.New(mat, mesh.Deferred),
		Path: path,
	}
	if err := obj.load(); err != nil {
		fmt.Println("Error loading model", path, ":", err)
	}
	return obj
}

func (obj *gltfModel) load() error {
	assetPath := fmt.Sprintf("assets/%s", obj.Path)
	doc, _ := gltf.Open(assetPath)

	pointers, _, vertexData := createBuffer(doc, doc.Meshes[0].Primitives[0])
	indexElements, indexData := createIndexBuffer(doc, doc.Meshes[0].Primitives[0])

	for i, ptr := range pointers {
		pointers[i].Name = strings.ToLower(ptr.Name)
	}

	primitive := mapPrimitiveType(doc.Meshes[0].Primitives[0].Mode)

	gmesh := &gltfMesh{
		id:        assetPath,
		primitive: primitive,
		elements:  indexElements,
		pointers:  pointers,
		vertices:  vertexData,
		indices:   indexData,
		indexsize: len(indexData) / indexElements,
	}
	obj.SetMesh(gmesh)

	return nil
}

func extractPointers(doc *gltf.Document, primitive *gltf.Primitive) []vertex.Pointer {
	offset := 0
	pointers := make(vertex.Pointers, len(primitive.Attributes))
	for name, index := range primitive.Attributes {
		accessor := doc.Accessors[primitive.Attributes[name]]

		pointers[index] = vertex.Pointer{
			Name:      name,
			Source:    mapComponentType(accessor.ComponentType),
			Offset:    offset,
			Elements:  int(accessor.Type.Components()),
			Normalize: accessor.Normalized,
			Stride:    0, // filed in in next pass
		}

		size := int(accessor.ComponentType.ByteSize() * accessor.Type.Components())
		offset += size
	}

	// at this point, offset equals the final stride value. fill it in
	for _, index := range primitive.Attributes {
		pointers[index].Stride = offset
	}

	return pointers
}

func createBuffer(doc *gltf.Document, primitive *gltf.Primitive) (vertex.Pointers, int, []byte) {
	pointers := extractPointers(doc, primitive)

	count := int(doc.Accessors[primitive.Attributes[pointers[0].Name]].Count)
	size := count * pointers[0].Stride

	output := make([]byte, size)

	for _, ptr := range pointers {
		accessor := doc.Accessors[primitive.Attributes[ptr.Name]]
		view := doc.BufferViews[*accessor.BufferView]
		buffer := doc.Buffers[view.Buffer]
		size := int(accessor.ComponentType.ByteSize() * accessor.Type.Components())
		stride := size
		if view.ByteStride != 0 {
			stride = int(view.ByteStride)
		}

		for i := 0; i < count; i++ {
			srcStart := int(view.ByteOffset) + i*stride + int(accessor.ByteOffset)
			srcEnd := srcStart + size
			dstStart := i*ptr.Stride + ptr.Offset
			dstEnd := dstStart + size

			copy(output[dstStart:dstEnd], buffer.Data[srcStart:srcEnd])
		}
	}

	return pointers, count, output
}

func createIndexBuffer(doc *gltf.Document, primitive *gltf.Primitive) (int, []byte) {
	accessor := doc.Accessors[*primitive.Indices]
	view := doc.BufferViews[*accessor.BufferView]
	buffer := doc.Buffers[view.Buffer]

	count := int(accessor.Count)
	size := int(accessor.ComponentType.ByteSize() * accessor.Type.Components())
	stride := size
	if view.ByteStride != 0 {
		stride = int(view.ByteStride)
	}

	output := make([]byte, size*count)
	for i := 0; i < count; i++ {
		srcStart := int(view.ByteOffset) + i*stride + int(accessor.ByteOffset)
		srcEnd := srcStart + size
		dstStart := i * size
		dstEnd := dstStart + size

		copy(output[dstStart:dstEnd], buffer.Data[srcStart:srcEnd])
	}

	return count, output
}

func mapPrimitiveType(mode gltf.PrimitiveMode) vertex.Primitive {
	switch mode {
	case gltf.PrimitiveTriangles:
		return vertex.Triangles
	case gltf.PrimitiveLines:
		return vertex.Lines
	default:
		panic("unsupported render primitive")
	}
}

func mapComponentType(kind gltf.ComponentType) types.Type {
	switch kind {
	case gltf.ComponentFloat:
		return types.Float
	case gltf.ComponentByte:
		return types.Int8
	case gltf.ComponentUbyte:
		return types.UInt8
	case gltf.ComponentShort:
		return types.Int16
	case gltf.ComponentUshort:
		return types.UInt16
	case gltf.ComponentUint:
		return types.UInt32
	default:
		panic(fmt.Sprintf("unmapped type %s (%d)", kind, kind))
	}
}
