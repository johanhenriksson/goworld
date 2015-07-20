package render

import (
    "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type MaterialTextureMap map[uint32]*Texture
type BufferDescriptors []BufferDescriptor

type BufferDescriptor struct {
    Buffer      uint32
    DataType    uint32
    Count       int32
    Stride      int32
    Offset      int32
    Normalize   bool
}

type Material struct {
    Shader      *ShaderProgram
    Textures    MaterialTextureMap
    Buffers     []BufferDescriptor
}

func CreateMaterial(shader *ShaderProgram) *Material {
    return &Material {
        Shader: shader,
        Textures:  make(MaterialTextureMap),
        Buffers: make(BufferDescriptors,0,3),
    }
}

func (mat *Material) AddDescriptor(attrName string, dataType uint32, count, stride, offset int, normalize bool) {
    loc := uint32(mat.Shader.GetAttrLoc(attrName))
    if loc < 0 {
        return
    }

    fmt.Printf("Adding descriptor %s at offset %d, stride: %d\n", attrName, offset, stride)

    mat.Buffers = append(mat.Buffers, BufferDescriptor {
        Buffer: loc,
        DataType: dataType,
        Count: int32(count),
        Stride: int32(stride),
        Normalize: normalize,
        Offset: int32(offset),
    })
}

func (mat *Material) AddTexture(slot uint32, tex *Texture) {
    mat.Textures[slot] = tex
}

func (mat *Material) Use() {
    i := 0
    for slot, tex := range mat.Textures {
        tex.Bind(slot)
        mat.Shader.UInt32(fmt.Sprintf("tex%d", i), slot)
        i++
    }
}

func (mat *Material) Setup() {
    for _, desc := range mat.Buffers {
        gl.EnableVertexAttribArray(desc.Buffer)
        gl.VertexAttribPointer(desc.Buffer, desc.Count, desc.DataType, desc.Normalize, desc.Stride, gl.PtrOffset(int(desc.Offset)))
    }
    /*
	gl.EnableVertexAttribArray(uint32(mat.vertices))
	gl.VertexAttribPointer(uint32(mat.vertices), int32(mat.VertexCount), mat.VertexType, false, int32(size), gl.PtrOffset(offset))
    */
}
