package render

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// todo: move to assets package

// Material file json representation
type materialDef struct {
	Shader   string
	Buffers  map[string][]*pointerDef
	Textures map[string]*textureDef
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
	File   string
	Filter string
	Wrap   string
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
	for buffer, pointers := range matf.Buffers {
		stride := 0
		for _, ptr := range pointers {
			// padding
			if ptr.Name == "skip" {
				stride += ptr.Count
				continue
			}

			if ptr.Count <= 0 {
				panic(fmt.Errorf("Expected count >0 for pointer %s", ptr.Name))
			}

			ptr.GlType, ptr.Size = getGlType(ptr.Type)
			ptr.Size *= ptr.Count
			ptr.Offset = stride
			stride += ptr.Size
		}
		for _, ptr := range pointers {
			if ptr.Name == "skip" {
				continue
			}
			mat.AddDescriptor(BufferDescriptor{
				Buffer:    buffer,
				Name:      ptr.Name,
				Type:      int(ptr.GlType),
				Elements:  ptr.Count,
				Stride:    stride,
				Offset:    ptr.Offset,
				Normalize: ptr.Normalize,
				Integer:   ptr.Integer,
			})
		}
	}

	/* Load textures */
	for name, txtf := range matf.Textures {
		texture, err := TextureFromFile(txtf.File)
		if err != nil {
			panic(err)
		}
		if txtf.Filter == "nearest" {
			texture.SetFilter(NearestFilter)
		}
		mat.AddTexture(name, texture)
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
