package engine

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// MeshBufferMap maps buffer names to vertex buffer objects
type MeshBufferMap map[string]*render.VertexBuffer

// Mesh base
type Mesh struct {
	*Transform
	Pass render.Pass

	name     string
	material *render.Material
	vao      *render.VertexArray
	pointers []render.VertexPointer
}

// NewMesh creates a new mesh object
func NewMesh(name string, material *render.Material) *Mesh {
	return NewPrimitiveMesh(name, render.Triangles, render.Geometry, material)
}

// NewLineMesh creates a new mesh for drawing lines
func NewLineMesh(name string) *Mesh {
	material := assets.GetMaterialShared("lines")
	return NewPrimitiveMesh(name, render.Lines, render.LinePass, material)
}

// NewPrimitiveMesh creates a new mesh composed of a given GL primitive
func NewPrimitiveMesh(name string, primitive render.GLPrimitive, pass render.Pass, material *render.Material) *Mesh {
	m := &Mesh{
		Transform: Identity(),
		Pass:      pass,
		name:      name,
		material:  material,
		vao:       render.CreateVertexArray(primitive),
	}
	return m
}

// Returns the name of the mesh
func (m *Mesh) Name() string {
	return m.name
}

// Buffer mesh data to GPU memory
// func (m *Mesh) Buffer(name string, data render.VertexData) error {
// 	m.vao.Buffer(name, data)
// 	for _, buffer := range m.material.Buffers {
// 		m.material.SetupBufferPointers(buffer)
// 	}
// 	return nil
// }

func (m *Mesh) SetIndexType(t render.GLType) {
	// get rid of this later
	m.vao.SetIndexType(t)
}

func (m *Mesh) Collect(pass DrawPass, args DrawArgs) {
	if m.Pass == pass.Type() && pass.Visible(m, args) {
		pass.Queue(m, args.Apply(m.Transform))
	}
}

func (m *Mesh) DrawDeferred(args DrawArgs) {
	m.material.Use()
	shader := m.material.Shader // UsePass(render.Geometry)

	// set up uniforms
	shader.Mat4("model", &args.Transform)
	shader.Mat4("view", &args.View)
	shader.Mat4("projection", &args.Projection)
	shader.Mat4("mvp", &args.MVP)
	shader.Vec3("eye", &args.Position)

	m.vao.Draw()
}

func (m *Mesh) DrawForward(args DrawArgs) {
	m.material.Use()
	shader := m.material.Shader // UsePass(render.Geometry)

	// set up uniforms
	shader.Mat4("model", &args.Transform)
	shader.Mat4("view", &args.View)
	shader.Mat4("projection", &args.Projection)
	shader.Mat4("mvp", &args.MVP)

	m.vao.Draw()
}

func (m *Mesh) DrawLines(args DrawArgs) {
	m.material.Use()
	m.material.Mat4("mvp", &args.MVP)

	m.vao.Draw()
}

func (m Mesh) Buffer(buffer string, data interface{}) {
	// infer pointers from struct tags
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Slice {
		return
	}

	el := t.Elem()
	if el.Kind() != reflect.Struct {
		fmt.Println("not a struct")
		return
	}

	size := int(el.Size())

	offset := 0
	pointers := make([]render.VertexPointer, 0, el.NumField())
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		tag, err := vertex.ParseTag(f.Tag.Get("vtx"))
		if err != nil {
			fmt.Printf("tag error on %s.%s: %s\n", el.Name(), f.Name, err)
			continue
		}

		gltype, err := render.GLTypeFromString(tag.Type)
		if err != nil {
			panic(fmt.Errorf("tag error on %s.%s: invalid GL type", el.Name(), f.Name))
		}

		attr, exists := m.material.Attribute(tag.Name)
		if !exists {
			// skip?
			fmt.Printf("attribute %s does not exist on %s\n", tag.Name, el.Name())
			offset += gltype.Size() * tag.Count
			continue
		}

		ptr := render.VertexPointer{
			Index:     int(attr.Index),
			Name:      tag.Name,
			Type:      gltype,
			Elements:  tag.Count,
			Normalize: tag.Normalize,
			Integer:   attr.Type.Integer(),
			Offset:    offset,
			Stride:    size,
		}

		pointers = append(pointers, ptr)

		fmt.Printf("%+v\n", ptr)

		offset += gltype.Size() * tag.Count
	}

	names := make([]string, len(pointers))
	for i := range pointers {
		names[i] = pointers[i].Name
	}
	bufferName := strings.Join(names, ",")

	m.vao.Buffer(bufferName, data)

	for _, ptr := range pointers {
		ptr.Enable()
	}

	// compatibility hack
	if len(pointers) == 0 {
		m.material.SetupVertexPointers()
	}
}
