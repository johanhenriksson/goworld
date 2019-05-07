package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type MaterialTextureMap map[string]*Texture
type BufferDescriptors []BufferDescriptor

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

type Material struct {
	Shader    *ShaderProgram
	Textures  MaterialTextureMap
	Buffers   []BufferDescriptor
	tex_slots []string // since map is unordered
}

/* Instantiate a new material */
func CreateMaterial(shader *ShaderProgram) *Material {
	return &Material{
		Shader:   shader,
		Textures: make(MaterialTextureMap),
		Buffers:  make(BufferDescriptors, 0, 0),
	}
}

func (mat *Material) AddDescriptor(desc BufferDescriptor) {
	loc, exists := mat.Shader.GetAttrLoc(desc.Name)
	if !exists {
		panic("No such attribute " + desc.Name)
	}
	desc.Buffer = int(loc)
	mat.Buffers = append(mat.Buffers, desc)
}

func (mat *Material) AddDescriptors(descriptors []BufferDescriptor) {
	for _, desc := range descriptors {
		mat.AddDescriptor(desc)
	}
}

/* Attach a texture to this material */
func (mat *Material) AddTexture(name string, tex *Texture) {
	mat.Textures[name] = tex
	mat.tex_slots = append(mat.tex_slots, name)
}

func (mat *Material) SetTexture(name string, tex *Texture) {
	mat.Textures[name] = tex
}

/* Set current shader and activate textures */
func (mat *Material) Use() {
	mat.Shader.Use()
	var i uint32 = 0
	//for name, tex := range mat.Textures {
	for _, name := range mat.tex_slots {
		tex := mat.Textures[name]
		tex.Use(i)
		mat.Shader.Int32(name, int32(i))
		i++
	}
	gl.ActiveTexture(gl.TEXTURE0)
}

func (mat *Material) EnablePointers() {
	for _, desc := range mat.Buffers {
		gl.EnableVertexAttribArray(uint32(desc.Buffer))
	}
}

func (mat *Material) DisablePointers() {
	for _, desc := range mat.Buffers {
		gl.DisableVertexAttribArray(uint32(desc.Buffer))
	}
}

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
