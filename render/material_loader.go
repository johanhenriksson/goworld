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

			ptr.GlType, ptr.Size = GLTypeFromString(ptr.Type)
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

// GLTypeFromString returns the GL identifier & size of a data type name
func GLTypeFromString(name string) (uint32, int) {
	switch name {
	case "byte":
		fallthrough
	case "int8":
		return gl.BYTE, 1

	case "ubyte":
		fallthrough
	case "uint8":
		fallthrough
	case "unsigned byte":
		return gl.UNSIGNED_BYTE, 1

	case "short":
		fallthrough
	case "int16":
		return gl.SHORT, 2

	case "ushort":
		fallthrough
	case "uint16":
		fallthrough
	case "unsigned short":
		return gl.UNSIGNED_SHORT, 2

	case "int":
		fallthrough
	case "int32":
		fallthrough
	case "integer":
		return gl.INT, 4

	case "uint":
		fallthrough
	case "uint32":
		fallthrough
	case "unsigned integer":
		return gl.UNSIGNED_INT, 4

	case "float":
		fallthrough
	case "float32":
		return gl.FLOAT, 4

	case "float64":
		fallthrough
	case "double":
		return gl.DOUBLE, 8
	}
	panic(fmt.Sprintf("Unknown GL type '%s'", name))
}
