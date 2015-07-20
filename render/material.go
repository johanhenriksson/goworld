package render

import (
    "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type MaterialTextureMap map[uint32]*Texture

type Material struct {
    Shader      *ShaderProgram
    Textures    MaterialTextureMap

    vertices    AttributeLocation
    normals     AttributeLocation
    texCoords   AttributeLocation
}

func CreateMaterial(shader *ShaderProgram, vertexAttr, normalsAttr, texAttr string) *Material {
    return &Material {
        Shader: shader,
        Textures:  make(MaterialTextureMap),

        vertices:  shader.GetAttrLoc(vertexAttr),
        normals:   shader.GetAttrLoc(normalsAttr),
        texCoords: shader.GetAttrLoc(texAttr),
    }
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
    size := 3 * 4
    offset := 0

    if mat.normals >= 0 { size += 3 * 4 }
    if mat.texCoords >= 0 { size += 2 * 4 }

	gl.EnableVertexAttribArray(uint32(mat.vertices))
	gl.VertexAttribPointer(uint32(mat.vertices), 3, gl.FLOAT, false, int32(size), gl.PtrOffset(offset))
    offset += 3 *  4

    if mat.normals >= 0 {
        gl.EnableVertexAttribArray(uint32(mat.normals))
        gl.VertexAttribPointer(uint32(mat.normals), 3, gl.FLOAT, false, int32(size), gl.PtrOffset(offset))
        offset += 3 *  4
    }

    if mat.texCoords >= 0 {
        gl.EnableVertexAttribArray(uint32(mat.texCoords))
        gl.VertexAttribPointer(uint32(mat.texCoords), 2, gl.FLOAT, false, int32(size), gl.PtrOffset(offset))
        offset += 2 *  4
    }
}
