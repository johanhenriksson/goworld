package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// MaterialTextureMap maps texture names to texture objects
type MaterialTextureMap map[string]*Texture

// BufferDescriptors is a list of vertex pointer descriptors
type BufferDescriptors []BufferDescriptor

// BufferDescriptor describes a vertex pointer into a buffer
type BufferDescriptor struct {
	Buffer    string
	Name      string
	Index     int
	Type      GLType
	Elements  int
	Stride    int
	Offset    int
	Normalize bool
	Integer   bool
}

// Material contains a shader reference and all resources required to draw a vertex buffer array
type Material struct {
	*Shader
	Textures    MaterialTextureMap
	Buffers     []string
	Descriptors []BufferDescriptor

	texslots []string // since map is unordered
	name     string
}

// CreateMaterial instantiates a new empty material
func CreateMaterial(name string, shader *Shader) *Material {
	return &Material{
		Shader:      shader,
		Textures:    make(MaterialTextureMap),
		Descriptors: make(BufferDescriptors, 0, 4),
		name:        name,
	}
}

func (mat *Material) String() string {
	return fmt.Sprintf("Material %s", mat.name)
}

// AddDescriptor adds a vertex pointer configuration
// Used to the describe the geometry format that will be drawn with this material
func (mat *Material) AddDescriptor(desc BufferDescriptor) {
	loc, exists := mat.getAttribute(desc.Name)
	if !exists {
		panic(fmt.Errorf("%s: No such attribute %s", mat, desc.Name))
	}
	desc.Index = int(loc.Index)
	mat.Descriptors = append(mat.Descriptors, desc)
	mat.addBuffer(desc.Buffer)
}

// AddDescriptors adds a list of vertex pointer configurations
// Used to the describe the geometry format that will be drawn with this material
func (mat *Material) AddDescriptors(descriptors []BufferDescriptor) {
	for _, desc := range descriptors {
		mat.AddDescriptor(desc)
	}
}

// addBuffer adds a buffer name to the buffer list if it does not already exist
func (mat *Material) addBuffer(buffer string) {
	for _, buffer := range mat.Buffers {
		if buffer == buffer {
			return
		}
	}
	mat.Buffers = append(mat.Buffers, buffer)
}

// AddTexture attaches a new texture to this material, and assings it to the next available texture slot.
func (mat *Material) AddTexture(name string, tex *Texture) {
	if _, exists := mat.Textures[name]; exists {
		mat.SetTexture(name, tex)
		return
	}
	mat.Textures[name] = tex
	mat.texslots = append(mat.texslots, name)
}

// SetTexture changes a bound texture
func (mat *Material) SetTexture(name string, tex *Texture) {
	mat.Textures[name] = tex
}

// Use sets the current shader and activates textures
func (mat *Material) Use() {
	mat.Shader.Use()
	i := uint32(0)
	for _, name := range mat.texslots {
		tex := mat.Textures[name]
		tex.Use(i)
		mat.Int32(name, int32(i))
		i++
	}
}

// EnablePointers enables vertex pointers used by this material
func (mat *Material) EnablePointers() {
	for _, desc := range mat.Descriptors {
		gl.EnableVertexAttribArray(uint32(desc.Index))
	}
}

// DisablePointers disables vertex pointers used by this material
func (mat *Material) DisablePointers() {
	for _, desc := range mat.Descriptors {
		gl.DisableVertexAttribArray(uint32(desc.Index))
	}
}

// SetupVertexPointers sets up vertex pointers used by this material.
// Use after binding the target vertex array object you want to configure!
func (mat *Material) SetupVertexPointers() {
	mat.EnablePointers()
	for _, desc := range mat.Descriptors {
		if desc.Integer {
			gl.VertexAttribIPointer(
				uint32(desc.Index),
				int32(desc.Elements),
				uint32(desc.Type),
				int32(desc.Stride),
				gl.PtrOffset(int(desc.Offset)))
		} else {
			gl.VertexAttribPointer(
				uint32(desc.Index),
				int32(desc.Elements),
				uint32(desc.Type),
				desc.Normalize,
				int32(desc.Stride),
				gl.PtrOffset(int(desc.Offset)))
		}
	}
}

// SetupBufferPointers sets up vertex pointers for a given buffer used by this material.
func (mat *Material) SetupBufferPointers(buffer string) {
	mat.EnablePointers()
	for _, desc := range mat.Descriptors {
		if desc.Buffer != buffer {
			continue
		}

		gl.EnableVertexAttribArray(uint32(desc.Index))
		if desc.Integer {
			gl.VertexAttribIPointer(
				uint32(desc.Index),
				int32(desc.Elements),
				uint32(desc.Type),
				int32(desc.Stride),
				gl.PtrOffset(int(desc.Offset)))
		} else {
			gl.VertexAttribPointer(
				uint32(desc.Index),
				int32(desc.Elements),
				uint32(desc.Type),
				desc.Normalize,
				int32(desc.Stride),
				gl.PtrOffset(int(desc.Offset)))
		}
	}
}
