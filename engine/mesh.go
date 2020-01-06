package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
)

type MeshBufferMap map[string]*render.VertexBuffer

type Mesh struct {
	*ComponentBase

	material *render.Material
	vao      *render.VertexArray
	vbos     MeshBufferMap
}

func NewMesh(material string) *Mesh {
	m := &Mesh{
		material: assets.GetMaterialCached(material),
		vao:      render.CreateVertexArray(),
		vbos:     MeshBufferMap{},
	}
	for _, buffer := range m.material.Buffers {
		m.addBuffer(buffer)
	}
	return m
}

func (m *Mesh) addBuffer(name string) *render.VertexBuffer {
	// create new vbo
	vbo := render.CreateVertexBuffer()

	// set up vertex array pointers for this buffer
	m.vao.Bind()
	vbo.Bind()
	m.material.SetupBufferPointers(name)

	// store reference & return vbo object
	m.vbos[name] = vbo
	return vbo
}

func (m *Mesh) Buffer(name string, data render.VertexData) error {
	m.vao.Bind()
	m.vao.Length = int32(data.Elements())
	vbo, exists := m.vbos[name]
	if !exists {
		return fmt.Errorf("Unknown VBO: %s", name)
	}
	return vbo.Buffer(data)
}

func (m *Mesh) Update(dt float32) {}

func (m *Mesh) Draw(args render.DrawArgs) {
	if args.Pass == render.GeometryPass {
		m.vao.Draw()
	}
}
