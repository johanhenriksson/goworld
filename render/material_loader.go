package render

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
    "github.com/go-gl/gl/v4.1-core/gl"
)

/** Material file JSON representation */
type f_material struct {
    Shader   string
    Pointers []*f_pointer
    Textures []*f_texture
}

/** Vertex pointer */
type f_pointer struct {
    Name      string
    Type      string
    GlType    uint32
    Size      int
    Offset    int
    Count     int
    Normalize bool
    Integer   bool
}

/** Texture definition */
type f_texture struct {
    Name string
    File string
}

/** Loads a material from a json definition file */
func LoadMaterial(shader *ShaderProgram, file string) *Material {
    json_bytes, err := ioutil.ReadFile(fmt.Sprintf("./%s.json", file))
    if err != nil {
        panic(err)
    }

    var matf f_material
    err = json.Unmarshal(json_bytes, &matf)
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
        mat.AddDescriptor(ptr.Name, ptr.GlType, ptr.Count, stride, ptr.Offset, ptr.Normalize, ptr.Integer)
    }

    /* Load textures */
    for _, txtf := range matf.Textures {
        texture, _ := LoadTexture(txtf.File)
        mat.AddTexture(txtf.Name, texture)
    }

    return mat
}

/** Returns the GL identifier & size of a data type name */
func getGlType(name string) (uint32,int) {
    switch(name) {
    case "byte":
        return gl.BYTE, 1
    case "unsigned byte":
        return gl.UNSIGNED_BYTE, 1
    case "float":
        return gl.FLOAT, 4
    }
    panic(fmt.Sprintf("Unknown GL type '%s'", name))
}
