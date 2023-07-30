package gltf

import (
	"fmt"
	"strings"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/types"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/qmuntal/gltf"
)

func Load(path string) object.Component {
	assetPath := fmt.Sprintf("assets/%s", path)
	doc, _ := gltf.Open(assetPath)

	// load default scene
	scene := doc.Scenes[*doc.Scene]

	return loadScene(doc, scene)
}

func loadScene(doc *gltf.Document, scene *gltf.Scene) object.Component {
	root := object.Empty(scene.Name)

	for _, nodeId := range scene.Nodes {
		node := loadNode(doc, doc.Nodes[nodeId])
		object.Attach(root, node)
	}

	// rotate to get Y+ up
	root.Transform().SetRotation(quat.Euler(90, 0, 0))

	return root
}

func loadNode(doc *gltf.Document, node *gltf.Node) object.Component {
	obj := object.Empty(node.Name)

	// mesh components
	if node.Mesh != nil {
		msh := doc.Meshes[*node.Mesh]
		for _, primitive := range msh.Primitives {
			renderer := loadPrimitive(doc, msh.Name, primitive)
			object.Attach(obj, renderer)
		}
	}

	// object transform
	obj.Transform().SetPosition(vec3.FromSlice(node.Translation[:3]))
	obj.Transform().SetRotation(quat.T{W: node.Rotation[3], V: vec3.FromSlice(node.Rotation[:3])})
	obj.Transform().SetScale(vec3.FromSlice(node.Scale[:3]))

	// child objects
	for _, child := range node.Children {
		object.Attach(obj, loadNode(doc, doc.Nodes[child]))
	}

	return obj
}

func loadPrimitive(doc *gltf.Document, name string, primitive *gltf.Primitive) mesh.Mesh {
	kind := mapPrimitiveType(primitive.Mode)

	// create interleaved buffers
	pointers, vertexData := createBuffer(doc, primitive)
	indexElements, indexData := createIndexBuffer(doc, primitive)

	// ensure vertex attribute names are in lowercase
	for i, ptr := range pointers {
		pointers[i].Name = strings.ToLower(ptr.Name)
	}

	// mesh data
	gmesh := &gltfMesh{
		key:       name,
		primitive: kind,
		elements:  indexElements,
		pointers:  pointers,
		vertices:  vertexData,
		indices:   indexData,
		indexsize: len(indexData) / indexElements,
	}

	// create mesh component
	mesh := mesh.NewPrimitiveMesh(kind, nil)
	mesh.VertexData.Set(gmesh)
	return mesh
}

func extractPointers(doc *gltf.Document, primitive *gltf.Primitive) []vertex.Pointer {
	offset := 0
	pointers := make(vertex.Pointers, 0, len(primitive.Attributes))
	for name, index := range primitive.Attributes {
		accessor := doc.Accessors[index]

		pointers = append(pointers, vertex.Pointer{
			Name:      name,
			Source:    mapComponentType(accessor.ComponentType),
			Offset:    offset,
			Elements:  int(accessor.Type.Components()),
			Normalize: accessor.Normalized,
			Stride:    0, // filed in in next pass
		})

		size := int(accessor.ComponentType.ByteSize() * accessor.Type.Components())
		offset += size
	}

	// at this point, offset equals the final stride value. fill it in
	for index := range pointers {
		pointers[index].Stride = offset
	}

	return pointers
}

func createBuffer(doc *gltf.Document, primitive *gltf.Primitive) (vertex.Pointers, []byte) {
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

	return pointers, output
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
