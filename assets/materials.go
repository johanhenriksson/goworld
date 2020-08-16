package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/johanhenriksson/goworld/render"
)

// MaterialDefinition file json representation
type MaterialDefinition struct {
	Shader   string
	Buffers  map[string][]*VertexPointerDefinition
	Textures map[string]*TextureDefinition
}

// VertexPointerDefinition json representation
type VertexPointerDefinition struct {
	Name      string
	Type      string
	GlType    render.GLType `json:"-"`
	Size      int
	Offset    int
	Count     int
	Normalize bool
	Integer   bool
}

// TextureDefinition json representation
type TextureDefinition struct {
	File   string
	Filter string
	Wrap   string
}

func LoadMaterialDefinition(file string) (*MaterialDefinition, error) {
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	matf := &MaterialDefinition{}
	err = json.Unmarshal(jsonBytes, matf)
	if err != nil {
		return nil, err
	}

	return matf, nil
}

func LoadMaterial(shader *render.ShaderProgram, matf *MaterialDefinition) (*render.Material, error) {
	mat := render.CreateMaterial(shader)

	/* Load vertex pointers */
	for buffer, pointers := range matf.Buffers {
		stride := 0
		for _, ptr := range pointers {
			if ptr.Count <= 0 {
				return nil, fmt.Errorf("Expected count >0 for pointer %s", ptr.Name)
			}

			// padding
			if ptr.Type == "skip" {
				stride += ptr.Count
				continue
			}

			// convert GL data type
			gltype, err := render.GLTypeFromString(ptr.Type)
			if err != nil {
				return nil, err
			}

			ptr.GlType = gltype
			ptr.Size = gltype.Size() * ptr.Count
			ptr.Offset = stride
			stride += ptr.Size
		}
		for _, ptr := range pointers {
			if ptr.Type == "skip" {
				continue
			}
			mat.AddDescriptor(render.BufferDescriptor{
				Buffer:    buffer,
				Name:      ptr.Name,
				Type:      ptr.GlType,
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
		texture, err := render.TextureFromFile(txtf.File)
		if err != nil {
			return nil, err
		}
		if txtf.Filter == "nearest" {
			texture.SetFilter(render.NearestFilter)
		}
		mat.AddTexture(name, texture)
	}

	return mat, nil
}

func GetMaterial(name string) *render.Material {
	path := fmt.Sprintf("assets/materials/%s.json", name)
	def, err := LoadMaterialDefinition(path)
	if err != nil {
		panic(err)
	}

	// use configured shader name if provided
	// otherwise fall back to material name
	shaderName := def.Shader
	if shaderName == "" {
		shaderName = name
	}

	shader := GetShader(shaderName)

	mat, err := LoadMaterial(shader, def)
	if err != nil {
		panic(err)
	}

	return mat
}
