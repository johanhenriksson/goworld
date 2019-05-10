package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// MaterialTextureMap maps texture names to texture objects
type MaterialTextureMap map[string]*Texture

// BufferDescriptors is a list of vertex pointer descriptors
type BufferDescriptors []BufferDescriptor

// BufferDescriptor describes a vertex pointer into a buffer
type BufferDescriptor struct {
	Name      string
	Buffer    int
	Type      int
	Elements  int
	Stride    int
	Offset    int
	Normalize bool
	Integer   bool
}

// Material contains a shader reference and all resources required to draw a vertex buffer array
type Material struct {
	Shader   *ShaderProgram
	Textures MaterialTextureMap
	Buffers  []BufferDescriptor
	slots    []string // since map is unordered
}

// CreateMaterial instantiates a new empty material
func CreateMaterial(shader *ShaderProgram) *Material {
	return &Material{
		Shader:   shader,
		Textures: make(MaterialTextureMap),
		Buffers:  make(BufferDescriptors, 0, 0),
	}
}

// AddDescriptor adds a vertex pointer configuration
// Used to the describe the geometry format that will be drawn with this material
func (mat *Material) AddDescriptor(desc BufferDescriptor) {
	loc, exists := mat.Shader.GetAttrLoc(desc.Name)
	if !exists {
		panic("No such attribute " + desc.Name)
	}
	desc.Buffer = int(loc)
	mat.Buffers = append(mat.Buffers, desc)
}

// AddDescriptors adds a list of vertex pointer configurations
// Used to the describe the geometry format that will be drawn with this material
func (mat *Material) AddDescriptors(descriptors []BufferDescriptor) {
	for _, desc := range descriptors {
		mat.AddDescriptor(desc)
	}
}

// AddTexture attaches a new texture to this material, and assings it to the next available texture slot.
func (mat *Material) AddTexture(name string, tex *Texture) {
	mat.Textures[name] = tex
	mat.slots = append(mat.slots, name)
}

// SetTexture changes a bound texture
func (mat *Material) SetTexture(name string, tex *Texture) {
	mat.Textures[name] = tex
}

// Use sets the current shader and activates textures
func (mat *Material) Use() {
	mat.Shader.Use()
	i := uint32(0)
	for _, name := range mat.slots {
		tex := mat.Textures[name]
		tex.Use(i)
		mat.Shader.Int32(name, int32(i))
		i++
	}
}

// EnablePointers enables vertex pointers used by this material
func (mat *Material) EnablePointers() {
	for _, desc := range mat.Buffers {
		gl.EnableVertexAttribArray(uint32(desc.Buffer))
	}
}

// DisablePointers disables vertex pointers used by this material
func (mat *Material) DisablePointers() {
	for _, desc := range mat.Buffers {
		gl.DisableVertexAttribArray(uint32(desc.Buffer))
	}
}

// SetupVertexPointers sets up vertex pointers used by this material.
// Use after binding the target vertex array object you want to configure!
func (mat *Material) SetupVertexPointers() {
	mat.EnablePointers()
	for _, desc := range mat.Buffers {
		if desc.Integer {
			gl.VertexAttribIPointer(
				uint32(desc.Buffer),
				int32(desc.Elements),
				uint32(desc.Type),
				int32(desc.Stride),
				gl.PtrOffset(int(desc.Offset)))
		} else {
			gl.VertexAttribPointer(
				uint32(desc.Buffer),
				int32(desc.Elements),
				uint32(desc.Type),
				desc.Normalize,
				int32(desc.Stride),
				gl.PtrOffset(int(desc.Offset)))
		}
	}
}
