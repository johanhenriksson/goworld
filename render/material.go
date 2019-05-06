package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type MaterialTextureMap map[string]*Texture
type BufferDescriptors []BufferDescriptor

type BufferDescriptor struct {
	Buffer    uint32
	DataType  uint32
	Count     int32
	Stride    int32
	Offset    int32
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

/* Add vertex attribute pointer */
func (mat *Material) AddDescriptor(attrName string, dataType uint32, count, stride, offset int, normalize, integer bool) {
	loc, exists := mat.Shader.GetAttrLoc(attrName)
	if !exists {
		panic("No such attribute " + attrName)
	}

	mat.Buffers = append(mat.Buffers, BufferDescriptor{
		Buffer:    uint32(loc),
		DataType:  dataType,
		Count:     int32(count),
		Stride:    int32(stride),
		Normalize: normalize,
		Offset:    int32(offset),
		Integer:   integer,
	})
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

func (mat *Material) SetupVertexPointers() {
	/* Enable vertex array attributes and set up vertex attribute pointers */
	for _, desc := range mat.Buffers {
		gl.EnableVertexAttribArray(desc.Buffer)
		if desc.Integer {
			gl.VertexAttribIPointer(
				desc.Buffer,
				desc.Count,
				desc.DataType,
				desc.Stride,
				gl.PtrOffset(int(desc.Offset)))
		} else {
			gl.VertexAttribPointer(
				desc.Buffer,
				desc.Count,
				desc.DataType,
				desc.Normalize,
				desc.Stride,
				gl.PtrOffset(int(desc.Offset)))
		}
	}
}
