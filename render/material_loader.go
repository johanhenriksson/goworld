package render

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Material file json representation
type materialDef struct {
	Shader   string
	Pointers []*pointerDef
	Textures []*textureDef
}

// Vertex pointer json representation
type pointerDef struct {
	Name      string
	Type      string
	GlType    uint32
	Size      int
	Offset    int
	Count     int
	Normalize bool
	Integer   bool
}

// Texture definition json representation
type textureDef struct {
	Name string
	File string
}

// LoadMaterial loads a material from a json definition file
func LoadMaterial(shader *ShaderProgram, file string) *Material {
	jsonBytes, err := ioutil.ReadFile(fmt.Sprintf("./%s.json", file))
	if err != nil {
		panic(err)
	}

	matf := materialDef{}
	err = json.Unmarshal(jsonBytes, &matf)
	if err != nil {
		panic(err)
	}

	shader.Use()
	mat := CreateMaterial(shader)

	/* Load vertex pointers */
	stride := 0
	for _, ptr := range matf.Pointers {
		if ptr.Name == "skip" {
			stride += ptr.Count
			continue
		}
		ptr.GlType, ptr.Size = getGlType(ptr.Type)
		ptr.Size *= ptr.Count
		ptr.Offset = stride
		stride += ptr.Size
	}
	for _, ptr := range matf.Pointers {
		if ptr.Name == "skip" {
			continue
		}
		mat.AddDescriptor(BufferDescriptor{
			Name:      ptr.Name,
			Type:      int(ptr.GlType),
			Elements:  ptr.Count,
			Stride:    stride,
			Offset:    ptr.Offset,
			Normalize: ptr.Normalize,
			Integer:   ptr.Integer,
		})
	}

	/* Load textures */
	for _, txtf := range matf.Textures {
		texture, _ := TextureFromFile(txtf.File)
		mat.AddTexture(txtf.Name, texture)
	}

	return mat
}

// getGlType returns the GL identifier & size of a data type name
func getGlType(name string) (uint32, int) {
	switch name {
	case "byte":
		return gl.BYTE, 1
	case "unsigned byte":
		return gl.UNSIGNED_BYTE, 1
	case "float":
		return gl.FLOAT, 4
	}
	panic(fmt.Sprintf("Unknown GL type '%s'", name))
}
